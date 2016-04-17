package main

import (
	"fmt"
	"github.com/shavac/readline"
	"strings"
    "strconv"
)

func main() {
	prompt := []string{"pet> ", "Attribute name> ", "Valid attribute types:\n 1) Integer ;; 2) Double ;; 3) Boolean ;; 4) String\n\nType> ", "Additional attribute (y/n)> ", "rid> "}
	for {
		rst := readline.ReadLine(&(prompt[0]))
		if rst != nil {
			result := strings.Split(strings.ToLower(strings.Trim(*rst, " \n\t")), " ")
			switch result[0] {
			case "quit":
				return
			case "create":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to create: have", len(result), " but expected 2.")
                    break
				}

                var attribute_names []string
                var attribute_types []int

                var stop bool = false

                for !stop {
                    var attribute_name string
                    var attribute_type int
                    var err error

                    for {
                        rst_2 := readline.ReadLine(&(prompt[1]))
                        if rst_2 != nil {
                            attribute_name = *rst_2;

                            if strings.Contains(attribute_name, ":") || strings.Contains(attribute_name, "[") || strings.Contains(attribute_name, "]") || strings.Contains(attribute_name, ",") {
                                fmt.Println("Invalid character in attribute name. Invalid characters are ':', ',', '['. and ']'.")
                                continue
                            } else {
                                var found bool = false
                                for i := range(attribute_names) {
                                    if attribute_names[i] == attribute_name {
                                        found = true
                                        break
                                    }
                                }

                                if (found) {
                                    fmt.Println("Name already in use; please specify another.")
                                    continue
                                }
                                break
                            }
                        } else {
                            fmt.Println("Invalid input. Please try again.")
                        }
                    }

                    for {
                        rst_2 := readline.ReadLine(&(prompt[2]))
                        if rst_2 != nil {
                            tmp_string := *rst_2;
                            attribute_type, err = strconv.Atoi(tmp_string)

                            if  err != nil || attribute_type < 1 || attribute_type > 4 {
                                fmt.Println("Invalid character in attribute type. Must be an integer, [1...4].")
                                continue
                            } else {
                                break
                            }
                        } else {
                            fmt.Println("Invalid input. Please try again.")
                        }
                    }

                    attribute_names = append(attribute_names, attribute_name)
                    attribute_types = append(attribute_types, attribute_type)

                    for {
                        rst_2 := readline.ReadLine(&(prompt[3]))
                        if rst_2 != nil {
                            tmp_string := strings.ToLower(strings.Trim(*rst_2, " \n\t"));

                            if tmp_string != "y" && tmp_string != "n" {
                                fmt.Println("Invalid character in attribute type. Must be either y or n.")
                                continue
                            } else {
                                if tmp_string == "n" {
                                    stop = true
                                }
                                break
                            }
                        } else {
                            fmt.Println("Invalid input. Please try again.")
                        }
                    }
                }

                TableCreate(attribute_names, attribute_types, result[1])
			case "header":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to header: have", len(result), " but expected 2.")
                    break
				}

                TableHeader(result[1])
			case "insert":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to insert: have", len(result), " but expected 2.")
                    break
				}

                TableInsert(result[1])
			case "display":
				if len(result) != 3 {
					fmt.Println("Error; invalid number of arguments to display: have", len(result), " but expected 3.")
                    break
				}

                var row_id int = -1

                row_id, err := strconv.Atoi(result[1])

                if err != nil || row_id < 0 {
                    for {
                        rst_2 := readline.ReadLine(&(prompt[4]))
                        if rst_2 != nil {
                            tmp_string := *rst_2;
                            row_id, err = strconv.Atoi(tmp_string)

                            if  err != nil || row_id < 0 {
                                fmt.Println("Invalid character in row id. Must be an integer greater than zero.")
                                continue
                            } else {
                                break
                            }
                        } else {
                            fmt.Println("Invalid input. Please try again.")
                        }
                    }
                }

                TableDisplay(row_id, result[2])
			case "delete":
				if len(result) != 3 {
					fmt.Println("Error; invalid number of arguments to delete: have", len(result), " but expected 3.")
                    break
				}

                var row_id int = -1

                row_id, err := strconv.Atoi(result[1])

                if err != nil || row_id < 0 {
                    for {
                        rst_2 := readline.ReadLine(&(prompt[4]))
                        if rst_2 != nil {
                            tmp_string := *rst_2;
                            row_id, err = strconv.Atoi(tmp_string)

                            if  err != nil || row_id < 0 {
                                fmt.Println("Invalid character in row id. Must be an integer greater than zero.")
                                continue
                            } else {
                                break
                            }
                        } else {
                            fmt.Println("Invalid input. Please try again.")
                        }
                    }
                }

                TableDelete(row_id, result[2])
			case "search":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to search: have", len(result), " but expected 3.")
                    break
				}

                // TODO : Search
            case "help":
                // TODO : help text
			default:
				fmt.Println("Unknown command:", result[0])
			}
		} else {
			fmt.Print("\n")
			return
		}
	}
}
