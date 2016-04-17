package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

func TableDelete(row_id int, filename string) {
    fmt.Println("Call to delete with:", filename, "and row id", row_id)

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

    if row_id >= records {
        fmt.Println("Fatal Error: Display index out of bounds; only have", records, "records.")
        return
    }

    // Remove header
    file = file[1:]

    // Remove given row
    if row_id == 0 {
        file = file[1:]
    } else if (row_id == records-1) {
        file = file[0:len(file)-1]
    } else {
        var saved []string = file[row_id+1:]
        file = file[0:row_id]
        file = append(file, saved...)
    }

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

    fmt.Println("Successfully deleted record id", row_id, "in table `", filename, "`!")
}
