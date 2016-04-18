package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func TableSearch(query string, filename string) {
	fmt.Println("Call to search with:", filename, "and query", query)

    var extended_query []string = strings.Split(query, "=")
    var target_column string
    var target_value string
    if len(extended_query) < 2 {
		fmt.Println("Error: malformed search query; ")
		return
    }

    target_column = strings.Trim(extended_query[0], " \t\n")
    target_value = strings.Trim(strings.Join(extended_query[1:], "="), " \t\n")

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
	header[len(header)-1] = header[len(header)-1][0 : len(header[len(header)-1])-1]

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

		records = len(file) - 1
	}

	var attribute_names []string
	var attribute_types []int

	for i := range header {
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


    var found_name int = -1;
    for i := range(attribute_names) {
        if attribute_names[i] == target_column {
            found_name = i
        }
    }

    if found_name == -1 {
        fmt.Println("Fatal error: Unknown column name `", target_column, "`!")
        return
    }

    // TODO : Validate that type of search value matches given value

    file = file[1:]
    for i := range(file) {
        var line string = file[i][1:len(file[i])-1]
        var values []string = strings.Split(line, "|")
        if (values[found_name] == target_value) {
            fmt.Println("rid:", i, "matches!")
        }
    }

	fmt.Println("Successfully searched in table `", filename, "`!")
}
