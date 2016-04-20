package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"strconv"
	"strings"
)

var columTypeToName map[int]string = map[int]string{1: "integer", 2: "double", 3: "boolean", 4: "string"}

func main() {
	prompt := []string{"pet> ", "Attribute name> ", "Valid attribute types:\n 1) Integer ;; 2) Double ;; 3) Boolean ;; 4) String\n\nType> ", "Additional attribute (y/n)> ", "rid> "}
    help_text := "PET: PET Editing of Tables\n--------------------------\nBy Alexander Scheel\n\nCommands\n========\ncreate <filename>\t\t\t--\tcreates a database; prompts for attributes\nheader <filename>\t\t\t--\tdisplays attributes of a database\ninsert <filename>\t\t\t--\tinserts into a database; prompts for values\ndisplay <rid> <filename>\t\t--\tdisplays the <rid>th entry of the database\ndelete <rid> <filename>\t\t\t--\tdeletes the <rid>th entry of the database\nsearch \"<condition>\" <filename>\t\t--\tsearches for the given condition in the database.\nhelp\t\t\t\t\t--\tprints this help message\n\n\n"


	var completer = readline.NewPrefixCompleter(
	    readline.PcItem("create"),
	    readline.PcItem("delete"),
	    readline.PcItem("display"),
	    readline.PcItem("header"),
	    readline.PcItem("insert"),
	    readline.PcItem("search"),
	    readline.PcItem("quit"),
	    readline.PcItem("exit"),
	    readline.PcItem("help"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt: prompt[0],
		InterruptPrompt: "^C",
		AutoComplete: completer,
	})

	if err != nil {
	    fmt.Println("Readline error:", err)
		return
	}
	defer rl.Close()


	for {
		line, err := rl.Readline()
		if err == nil {
			result := strings.Split(strings.ToLower(strings.Trim(line, " \t\n")), " ")
			switch result[0] {
			case "quit":
				return
			case "exit":
				return
			case "create":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to create: have", len(result), "but expected 2.")
					break
				}

				var attribute_names []string
				var attribute_types []int

				var stop bool = false

				for !stop {
					var attribute_name string
					var attribute_type int

					for {
						rl2, err := readline.New(prompt[1])
						if err != nil {
						    fmt.Println("Readline error:", err)
							return
						}
						defer rl2.Close()
						line2, err := rl2.Readline()
						if err == nil {
							attribute_name = line2

							if strings.Contains(attribute_name, ":") || strings.Contains(attribute_name, "[") || strings.Contains(attribute_name, "]") {
								fmt.Println("Invalid character in attribute name. Invalid characters are ':', '['. and ']'.")
								continue
							} else {
								var found bool = false
								for i := range attribute_names {
									if attribute_names[i] == attribute_name {
										found = true
										break
									}
								}

								if found {
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
						rl2, err := readline.New(prompt[2])
						if err != nil {
						    fmt.Println("Readline error:", err)
							return
						}
						defer rl2.Close()
						line2, err := rl2.Readline()
						if err == nil {
							tmp_string := line2
							attribute_type, err = strconv.Atoi(tmp_string)

							if err != nil || attribute_type < 1 || attribute_type > 4 {
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
						rl2, err := readline.New(prompt[3])
						if err != nil {
						    fmt.Println("Readline error:", err)
							return
						}
						defer rl2.Close()
						line2, err := rl2.Readline()
						if err == nil {
							tmp_string := strings.ToLower(strings.Trim(line2, " \n"))

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
					fmt.Println("Error; invalid number of arguments to header: have", len(result), "but expected 2.")
					break
				}

				TableHeader(result[1])
			case "insert":
				if len(result) != 2 {
					fmt.Println("Error; invalid number of arguments to insert: have", len(result), "but expected 2.")
					break
				}

				TableInsert(result[1])
			case "display":
				if len(result) != 3 {
					fmt.Println("Error; invalid number of arguments to display: have", len(result), "but expected 3.")
					break
				}

				var row_id int = -1

				row_id, err := strconv.Atoi(result[1])

				if err != nil || row_id < 0 {
					for {
						rl2, err := readline.New(prompt[4])
						if err != nil {
						    fmt.Println("Readline error:", err)
							return
						}
						defer rl2.Close()
						line2, err := rl2.Readline()
						if err == nil {
							tmp_string := line2
							row_id, err = strconv.Atoi(tmp_string)

							if err != nil || row_id < 0 {
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
					fmt.Println("Error; invalid number of arguments to delete: have", len(result), "but expected 3.")
					break
				}

				var row_id int = -1

				row_id, err := strconv.Atoi(result[1])

				if err != nil || row_id < 0 {
					for {
						rl2, err := readline.New(prompt[4])
						if err != nil {
						    fmt.Println("Readline error:", err)
							return
						}
						defer rl2.Close()
						line2, err := rl2.Readline()
						if err == nil {
							tmp_string := line2
							row_id, err = strconv.Atoi(tmp_string)

							if err != nil || row_id < 0 {
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
    			query := strings.Split(strings.Trim(line, " \t\n"), "\"")
                if len(query) != 3 {
					fmt.Println("Error; invalid number of arguments to search: have", len(query), "but expected at least 3.")
					break
                }

                query[1] = strings.Trim(query[1], " \t\n")
                query[2] = strings.Trim(query[2], " \t\n")

                if len(query[2]) == 0 || query[2] != result[len(result)-1] {
					fmt.Println("Error; invalid filename after search query.")
					break
                }

                TableSearch(query[1], query[2])
			case "help":
				fmt.Print(help_text)
			default:
				fmt.Println("Unknown command:", result[0])
				fmt.Print(help_text)
			}
		} else {
			fmt.Print("\n")
			return
		}
	}
}
