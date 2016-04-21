package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type query_part struct {
	Name      string
	ID        int
	Value     string
	Operation int
	Required  bool
}

/**
 * Basic tokenizer; splits on space
 *
 * Operators:
 * 		=  -- 0
 * 		>  -- 1
 * 		<  -- 2
 *  	>= -- 3
 * 		<= -- 4
 * 		!= -- 5
 *
 * name = 1
 * name     =             1
 * ^ name   ^ operation   ^ value
**/
func parseQuery(query string) []query_part {
	var parts []string = strings.Split(query, " ")
	var result []query_part = make([]query_part, 0)

	var operators map[string]int = map[string]int{"=": 0, ">": 1, "<": 2, ">=": 3, "<=": 4, "!=": 5}

	if len(parts) % 3 != 0 {
		fmt.Println("Fatal Error: invalid query.")
		return nil
	}

	for i := 0; i < len(parts); i+=3 {
		var current query_part
		current.Name = parts[i]
		current.ID = -1
		current.Value = parts[i+2]
		current.Required = true

		value, ok := operators[parts[i+1]]
		if !ok {
			fmt.Println("Fatal Error: Invalid operator.")
			return nil
		}

		current.Operation = value

		result = append(result, current)
	}

	return result
}

func TableSearch(query string, filename string) {
	fmt.Println("Call to search with:", filename, "and query", query)

	var full_parsed_query []query_part = parseQuery(query)

	if full_parsed_query == nil {
		return
	}

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


    for i := range(attribute_names) {
		for j := range(full_parsed_query) {
			var parsed_query query_part = full_parsed_query[j]

	        if attribute_names[i] == parsed_query.Name {
	            parsed_query.ID = i
				if attribute_types[i] == 1 {
					_, err = strconv.Atoi(parsed_query.Value)

					if err != nil {
						fmt.Println("Fatal Error: Unable to convert search query to integer", err)
						return
					}
				} else if attribute_types[i] == 2 {
					_, err = strconv.ParseFloat(parsed_query.Value, 64)

					if err != nil {
						fmt.Println("Fatal Error: Unable to convert search query to double", err)
						return
					}
				} else if attribute_types[i] == 3 {
					parsed_query.Value = strings.ToUpper(parsed_query.Value)

					if parsed_query.Value != "T" && parsed_query.Value != "F" {
						fmt.Println("Fatal Error: search query unknown boolean value: expected either T or F.")
						return
					}
				} else if attribute_types[i] == 4 {
					if strings.Contains(parsed_query.Value, "|") || strings.Contains(parsed_query.Value, "{") || strings.Contains(parsed_query.Value, "}") {
						fmt.Println("Invalid character in search query string value. Invalid characters are '|', '{'. and '}'.")
						return
					}
				}
	        }
		}
    }

    if parsed_query.ID == -1 {
        fmt.Println("Fatal error: Unknown column name `", parsed_query.Name, "`!")
        return
    }

    file = file[1:]
    for i := range(file) {
        var line string = file[i][1:len(file[i])-1]
        var values []string = strings.Split(line, "|")

		for j := range(full_parsed_query) {
			var parsed_query query_part = full_parsed_query[j]
		}
    }

	fmt.Println("Successfully searched in table `", filename, "`!")
}
