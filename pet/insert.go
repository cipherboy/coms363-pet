package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
	"github.com/shavac/readline"
)

func TableInsert(filename string) {
    fmt.Println("Call to insert with:", filename)

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        fmt.Println("Error: file `", filename, "` does not exist...")
        return
    }

    var file []string

    f, err := os.Open(filename)
    s := bufio.NewScanner(f)
    for s.Scan() {
        file = append(file, s.Text())
    }

    if file[0][0] != '[' {
        fmt.Println("Error: malformed file. Unknown character `", file[0][0], "` at line 0 position 0.")
        return
    }

    var header []string = strings.Split(file[0], "][")
    var columns int
    var records int

    header[0] = header[0][1:]
    header[len(header)-1] = header[len(header)-1][0:len(header[len(header)-1])-1]

    columns, err = strconv.Atoi(header[0])
    if err != nil {
        fmt.Println("Fatal Error: malformed file. Cannot parse column count as integer:", err)
        return
    }

    if len(header)-2 != columns {
        fmt.Println("Fatal Error: malformed file. Number of column does not match header column count:", len(header)-2, "!=", columns)
        return
    }

    records, err = strconv.Atoi(header[len(header)-1])
    if err != nil {
        fmt.Println("Fatal Error: malformed file. Cannot parse record count as integer:", err)
        return
    }

    if len(file)-1 != records {
        fmt.Println("Recoverable Error: Number of records do not match header record count. Using number in file:", len(file)-1, "vs", records)

        records = len(file)-1
    }

    var attribute_names []string
    var attribute_types []int

    for i := range(header) {
        if i == 0 || i == len(header)-1 {
            continue
        }

        var item []string = strings.Split(header[i], ":")
        if len(item) != 2 {
            fmt.Println("Fatal Error: malformed header. Expected two attributes in column", i, ": got", len(item))
            return
        }

        var attribute_name string = item[0]
        var attribute_type int = 0

        attribute_type, err = strconv.Atoi(item[1])
        if err != nil || attribute_type < 1 || attribute_type > 4 {
            fmt.Println("Fatal Error: malformed header. In column", i, ": cannot parse `", item[1], "` as integer. Error: ", err)
            return
        }

        attribute_names = append(attribute_names, attribute_name)
        attribute_types = append(attribute_types, attribute_type)
    }

    var record_data []string

    for i := range(attribute_names) {
        var attribute_data string

        for {
            prompt := attribute_names[i] + " (" + columTypeToName[attribute_types[i]] + ")> "
            rst := readline.ReadLine(&prompt)

            if rst == nil {
                fmt.Println("Unknown input.")
                continue
            }

            attribute_data = strings.Trim(*rst, " \n\t")

            if attribute_types[i] == 1 {
                _, err = strconv.Atoi(attribute_data)

                if err != nil {
                    fmt.Println("Unable to convert input to integer; please try again:", err)
                    continue
                } else {
                    break
                }
            } else if attribute_types[i] == 2 {
                _, err = strconv.ParseFloat(attribute_data, 64)

                if err != nil {
                    fmt.Println("Unable to convert input to double; please try again:", err)
                    continue
                } else {
                    break
                }
            } else if attribute_types[i] == 3 {
                attribute_data = strings.ToUpper(attribute_data)

                if attribute_data != "T" && attribute_data != "F" {
                    fmt.Println("Unknown boolean value: expected either T or F.")
                    continue
                } else {
                    break
                }
            } else if attribute_types[i] == 4 {
                if strings.Contains(attribute_data, "|") || strings.Contains(attribute_data, "{") || strings.Contains(attribute_data, "}") {
                    fmt.Println("Invalid character in string value. Invalid characters are '|', '{'. and '}'.")
                    continue
                } else {
                    break
                }
            }
        }

        record_data = append(record_data, attribute_data)
    }

    // Remove header, append new record
    file = file[1:]
    file = append(file, "{" + strings.Join(record_data, "|") + "}")

    var header_string string
    // Build a new header with updated info
    header_string = "[" + strconv.Itoa(len(attribute_names)) + "]"

    for i := range(attribute_names) {
        header_string += "[" + attribute_names[i] + ":" + strconv.Itoa(attribute_types[i]) + "]"
    }
    header_string += "[" + strconv.Itoa(len(file)) + "]\n"

    err = os.Remove(filename)
    if err != nil {
        fmt.Println("Fatal Error: cannot remove file:", err)
        return
    }

    fw, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer fw.Close()

    wl, err := fw.Write([]byte(header_string))
    if err != nil {
        fmt.Println("Fatal Error writing file:", err)
        return
    }

    if wl != len(header_string) {
        fmt.Println("Fatal Error writing file: wrote", wl, "bytes but expected to write", len(header_string))
        return
    }

    for i := range(file) {
        wl, err := fw.Write([]byte(file[i] + "\n"))
        if err != nil {
            fmt.Println("Fatal Error writing file:", err)
            return
        }

        if wl != len(file[i] + "\n") {
            fmt.Println("Fatal Error writing file: wrote", wl, "bytes but expected to write", len(file[i]))
            return
        }
    }

    fmt.Println("Successfully inserted into table `", filename, "`!")
}
