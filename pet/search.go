package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func bytes_contains(needle byte, haystack []byte) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

func strings_contains(needle string, haystack []string) int {
	for i := range haystack {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

/**
 * Value: literal token (string)
 * Type:
 *      Undefined:  -1
 *      Operator:   0
 *      Bareword:   1
 *      Join:       2
 *      String:     3
 *      Number:     4
**/
var unknown_token_type int = -1
var operator_token_type int = 0
var bareword_token_type int = 1
var join_token_type int = 2
var string_token_type int = 3
var number_token_type int = 4
var token_types_to_names map[int]string = map[int]string{-1: "unknown", 0: "operator", 1: "bareword", 2: "join", 3: "string", 4: "number"}

type token struct {
	Value string
	Type  int
}

func tokenizeQuery(query string) ([]token, error) {
	var result []token

	var whitespace_parts []byte = []byte(" \t\n")
	var join_parts []byte = []byte("&|")
	var operator_parts []byte = []byte("><=!")
	var bareword_parts []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")
	var number_parts []byte = []byte("0123456789.")
	var string_start byte = '\''
	var string_end byte = '\''

	for i := 0; i < len(query); i++ {
		var current token
		current.Type = unknown_token_type

		// Ignore whitespace
		if bytes_contains(query[i], whitespace_parts) != -1 {
			continue
		} else if bytes_contains(query[i], operator_parts) != -1 {
			current.Value += string(query[i])
			current.Type = operator_token_type

			// Look ahead and catch next operator part, if it exists
			for i+1 < len(query) && bytes_contains(query[i+1], operator_parts) != -1 {
				current.Value += string(query[i+1])
				i += 1
			}
		} else if bytes_contains(query[i], number_parts) != -1 {
			current.Value += string(query[i])
			current.Type = number_token_type

			// Look ahead and catch next number part, if it exists
			for i+1 < len(query) && bytes_contains(query[i+1], number_parts) != -1 {
				current.Value += string(query[i+1])
				i += 1
			}
		} else if bytes_contains(query[i], bareword_parts) != -1 {
			current.Value += string(query[i])
			current.Type = bareword_token_type

			// Look ahead and catch next bareword part, if it exists
			for i+1 < len(query) && bytes_contains(query[i+1], bareword_parts) != -1 {
				current.Value += string(query[i+1])
				i += 1
			}
		} else if bytes_contains(query[i], join_parts) != -1 {
			current.Value += string(query[i])
			current.Type = join_token_type

			// Look ahead and catch next join part, if it exists
			for i+1 < len(query) && bytes_contains(query[i+1], join_parts) != -1 {
				current.Value += string(query[i+1])
				i += 1
			}
		} else if query[i] == string_start {
			current.Value += string(query[i])
			current.Type = string_token_type

			// Add to string until end of string or end of query
			var found_end bool = false
			for i+1 < len(query) {
				current.Value += string(query[i+1])
				i += 1
				if query[i] == string_end {
					found_end = true
					break
				}
			}

			if !found_end {
				return []token(nil), errors.New("Unterminated string!")
			}

			current.Value = current.Value[1 : len(current.Value)-1]
		} else {
			return []token(nil), errors.New("Unknown character: `" + string(query[i]) + "`")
		}

		result = append(result, current)
	}

	return result, nil
}

/**
 * Value: set of tokens
 * Type:
 *      Undefined:  -1
 *      Relation:   0
 *      Join:       1
**/
var undefined_rtoken_type int = -1
var relation_rtoken_type int = 0
var join_rtoken_type int = 1
var rtoken_types_to_names map[int]string = map[int]string{-1: "unknown", 0: "relation", 1: "join"}

type rtoken struct {
	Value []token
	Type  int
}

func relationizeTokens(set []token) ([]rtoken, error) {
	var result []rtoken

	for i := 0; i < len(set); i++ {
		var current rtoken

		if set[i].Type == unknown_token_type {
			return []rtoken(nil), errors.New("Invalid token: Unknown token type: -1")
		} else if set[i].Type == bareword_token_type {
			current.Type = relation_rtoken_type
			current.Value = append(current.Value, set[i])

			i += 1

			if set[i].Type == operator_token_type {
				current.Value = append(current.Value, set[i])

				i += 1

				if set[i].Type == string_token_type || set[i].Type == number_token_type || set[i].Type == bareword_token_type {
					current.Value = append(current.Value, set[i])
				} else {
					return []rtoken(nil), errors.New("Invalid relation: cannot have type " + token_types_to_names[set[i].Type] + " (" + strconv.Itoa(set[i].Type) + ") after type bareword (" + strconv.Itoa(bareword_token_type) + ")")
				}
			} else {
				return []rtoken(nil), errors.New("Invalid relation: cannot have type " + token_types_to_names[set[i].Type] + " (" + strconv.Itoa(set[i].Type) + ") after type bareword (" + strconv.Itoa(bareword_token_type) + ")")
			}
		} else if set[i].Type == join_token_type {
			current.Type = join_rtoken_type
			current.Value = append(current.Value, set[i])
		} else {
			return []rtoken(nil), errors.New("Invalid relation: cannot have type " + token_types_to_names[set[i].Type] + " (" + strconv.Itoa(set[i].Type) + ") at this location.")
		}

		result = append(result, current)
	}

	return result, nil
}

func validateRelations(set []rtoken, column_names []string, column_types []int) error {
	if set[0].Type == join_rtoken_type || set[len(set)-1].Type == join_rtoken_type {
		return errors.New("Invalid relation: cannot have relation set begin or end with type join.")
	}

	for i := 0; i < len(set)-1; i++ {
		if set[i].Type == set[i+1].Type {
			return errors.New("Invalid relation: cannot have adjacent tokens of type " + rtoken_types_to_names[set[i].Type] + "(" + strconv.Itoa(set[i].Type) + ").")
		}
	}

	var valid_number_operators []string = []string{"==", "=", "!=", ">", "<", "<=", ">="}
	var valid_string_operators []string = []string{"==", "=", "!="}
	var valid_join_operators []string = []string{"&", "&&", "||", "|"}
	var valid_boolean_types []string = []string{"t", "T", "f", "F"}

	for i := range set {
		var tokens []token = set[i].Value
		if set[i].Type == relation_rtoken_type {
			if len(tokens) != 3 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting three tokens in relation")
			}

			if tokens[0].Type != bareword_token_type {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting left most token to be bareword")
			}

			var found_column_id int = strings_contains(tokens[0].Value, column_names)

			if found_column_id == -1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown bareword column name: " + tokens[0].Value)
			}

			if tokens[1].Type != operator_token_type {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting middle token to be operator")
			}

			if tokens[2].Type != bareword_token_type && tokens[2].Type != string_token_type && tokens[2].Type != number_token_type {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting right most token to be one of bareword, string, or number type.")
			}

			if tokens[2].Type == number_token_type && strings_contains(tokens[1].Value, valid_number_operators) == -1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown operator for numbers: " + tokens[1].Value)
			}

			if tokens[2].Type == number_token_type && column_types[found_column_id] > 2 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): column is not of numerical type: " + columnTypeToName[column_types[found_column_id]] + " vs " + tokens[2].Value)
			}

			if tokens[2].Type != number_token_type && strings_contains(tokens[1].Value, valid_string_operators) == -1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown operator for strings: " + tokens[1].Value)
			}

			if tokens[2].Type != number_token_type && column_types[found_column_id] < 3 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): column is of numerical type: " + columnTypeToName[column_types[found_column_id]] + " vs " + tokens[2].Value)
			}

			if tokens[2].Type != number_token_type && column_types[found_column_id] == 3 && strings_contains(tokens[2].Value, valid_boolean_types) == -1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): search value is not of boolean type: " + tokens[2].Value)
			}
		} else if set[i].Type == join_rtoken_type {
			if len(tokens) != 1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting only one tokens in join")
			}

			if tokens[0].Type != join_token_type {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting tokens in relation join to be of type join.")
			}

			if strings_contains(tokens[0].Value, valid_join_operators) == -1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown join operator: " + tokens[0].Value)
			}
		} else {
			if len(tokens) != 1 {
				return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown type: " + strconv.Itoa(i))
			}
		}
	}

	return nil
}

/**
 * Join:
 *      single: -1
 *      and:    0
 *      or:     1
 * Left: left relation
 * Right: right relation
 * Relation: relation value
 * Value: evaluated relation
 * Evaluated: status of node
**/
var join_evalTree_types map[string]int = map[string]int{"&&": 0, "&": 0, "||": 1, "|": 1}
var single_evalTree_type = -1
var and_evalTree_type = 0
var or_evalTree_type = 1

type evalTree struct {
	Join      int
	Left      *evalTree
	Right     *evalTree
	Relation  []token
	Value     bool
	Evaluated bool
}

func evalTreeizeRelation(set []rtoken) (evalTree, error) {
	var result evalTree
	result.Left = nil
	result.Right = nil
	result.Join = -1
	result.Value = false
	result.Evaluated = false

	for i := 0; i < len(set); i++ {
		if set[i].Type == relation_rtoken_type {
			var lone_relation evalTree
			lone_relation.Left = nil
			lone_relation.Right = nil
			lone_relation.Join = -1
			lone_relation.Value = false
			lone_relation.Relation = set[i].Value
			lone_relation.Evaluated = false

			var assigned bool
			top := &result

			for top != nil {
				if top.Left == nil {
					top.Left = &lone_relation
					assigned = true
					top = nil
				} else if top.Right == nil {
					top.Right = &lone_relation
					top = nil
					assigned = true
				} else if top.Right.Join != -1 {
					top = top.Right
				} else {
					top = nil
					assigned = false
				}
			}

			if !assigned {
				return result, errors.New("Invalid Evaluation Tree: Unable to add new relation (" + strconv.Itoa(i) + ") to root: all full")
			}

		} else if set[i].Type == join_rtoken_type {
			// And takes precedence, i.e., goes lower, than or, left to right
			if result.Join == single_evalTree_type {
				var ok bool
				result.Join, ok = join_evalTree_types[set[i].Value[0].Value]
				if !ok {
					return result, errors.New("Invalid Evaluation Tree: Unknown join operator: " + set[i].Value[0].Value)
				}
			} else if result.Join == and_evalTree_type {
				var new_root evalTree
				new_root.Left = nil
				new_root.Right = nil
				new_root.Join = -1
				new_root.Value = false
				new_root.Evaluated = false

				var ok bool
				new_root.Join, ok = join_evalTree_types[set[i].Value[0].Value]
				if !ok {
					return result, errors.New("Invalid Evaluation Tree: Unknown join operator: " + set[i].Value[0].Value)
				}

				var current_root evalTree = result
				new_root.Left = &current_root
				result = new_root
			} else if result.Join == or_evalTree_type {
				var new_right evalTree
				new_right.Left = nil
				new_right.Right = nil
				new_right.Join = -1
				new_right.Value = false
				new_right.Evaluated = false

				var ok bool
				new_right.Join, ok = join_evalTree_types[set[i].Value[0].Value]
				if !ok {
					return result, errors.New("Invalid Evaluation Tree: Unknown join operator: " + set[i].Value[0].Value)
				}

				new_right.Left = result.Right
				result.Right = &new_right
			}
		}
	}

	return result, nil
}

func prettyEvalTree(root *evalTree) string {
	if root == nil {
		return ""
	}

	var result string
	if root.Join == -1 {
		if root.Left != nil {
			for i := range root.Left.Relation {
				result += " " + root.Left.Relation[i].Value
			}
		} else if root.Relation != nil {
			for i := range root.Relation {
				result += " " + root.Relation[i].Value
			}
			result = result[1:]
		}
	} else if root.Join == 0 {
		result = "(" + prettyEvalTree(root.Left) + " && " + prettyEvalTree(root.Right) + ")"
	} else if root.Join == 1 {
		result = "(" + prettyEvalTree(root.Left) + " || " + prettyEvalTree(root.Right) + ")"
	}

	return result
}

func evaluateTreeForRow(root evalTree, column_names []string, column_types []int, row []string) bool {
	var copy evalTree = root
    recursiveEvaluateTreeForRow(&copy, column_names, column_types, row)

    if copy.Evaluated == false {
        fmt.Println("Error evaluating tree...")
    }

    return copy.Value
}

func recursiveEvaluateTreeForRow(root *evalTree, column_names []string, column_types []int, row []string) {
	if root == nil {
		return
	}

	if root.Join == -1 {
		if root.Left != nil {
            root.Left.Evaluated = true
            root.Left.Value = evaluateRelationForRow(root.Left.Relation, column_names, column_types, row)
            root.Evaluated = true
            root.Value = root.Left.Value
		} else if root.Relation != nil {
            root.Evaluated = true
            root.Value = evaluateRelationForRow(root.Relation, column_names, column_types, row)
		}
	} else if root.Join == 0 {
        recursiveEvaluateTreeForRow(root.Left, column_names, column_types, row)
        recursiveEvaluateTreeForRow(root.Right, column_names, column_types, row)

        if root.Left != nil && root.Left.Evaluated == true {
            root.Value = root.Left.Value
            root.Evaluated = true

            if root.Right != nil && root.Right.Evaluated == true {
                root.Value = root.Value && root.Right.Value
            }
        } else {
            if root.Right != nil && root.Right.Evaluated == true {
                root.Value = root.Right.Value
                root.Evaluated = true
            } else {
                root.Evaluated = false
            }
        }
	} else if root.Join == 1 {
        recursiveEvaluateTreeForRow(root.Left, column_names, column_types, row)
        recursiveEvaluateTreeForRow(root.Right, column_names, column_types, row)

        if root.Left != nil && root.Left.Evaluated == true {
            root.Value = root.Left.Value
            root.Evaluated = true

            if root.Right != nil && root.Right.Evaluated == true {
                root.Value = root.Value || root.Right.Value
            }
        } else {
            if root.Right != nil && root.Right.Evaluated == true {
                root.Value = root.Right.Value
                root.Evaluated = true
            } else {
                root.Evaluated = false
            }
        }
	}

	return
}

func evaluateRelationForRow(tokens []token, column_names []string, column_types []int, row []string) bool {
    if len(tokens) != 3 {
        return false
    }

    var found_column_id int = strings_contains(tokens[0].Value, column_names)

    if found_column_id == -1 {
        fmt.Println("Unknown bareword column name: " + tokens[0].Value)
        return false
    }

    var row_value string = row[found_column_id]
    var comparison_value string = tokens[2].Value

	if column_types[found_column_id] == 1 {
		real_row_value, err := strconv.Atoi(row_value )

		if err != nil {
			fmt.Println("Unable to convert row value to integer:", err)
            return false
		} else {
    		real_comparison_value, err := strconv.Atoi(comparison_value)

    		if err != nil {
    			fmt.Println("Unable to convert comparison value to integer:", err)
                return false
    		} else {
                if tokens[1].Value == "=" || tokens[1].Value == "==" {
                    return real_row_value == real_comparison_value
                } else if tokens[1].Value == "!=" {
                    return real_row_value != real_comparison_value
                } else if tokens[1].Value == ">" {
                    return real_row_value > real_comparison_value
                } else if tokens[1].Value == "<" {
                    return real_row_value < real_comparison_value
                } else if tokens[1].Value == "<=" {
                    return real_row_value <= real_comparison_value
                } else if tokens[1].Value == ">=" {
                    return real_row_value >= real_comparison_value
                } else {
                    fmt.Println("Unknown comparison operator: " + tokens[1].Value)
                    return false
                }
    		}
		}
	} else if column_types[found_column_id] == 2 {
		real_row_value, err := strconv.ParseFloat(row_value, 64)

		if err != nil {
			fmt.Println("Unable to convert row value to double:", err)
			return false
		} else {
    		real_comparison_value, err := strconv.ParseFloat(comparison_value, 64)

    		if err != nil {
    			fmt.Println("Unable to convert comparison value to double:", err)
    			return false
    		} else {
                if tokens[1].Value == "=" || tokens[1].Value == "==" {
                    return real_row_value == real_comparison_value
                } else if tokens[1].Value == "!=" {
                    return real_row_value != real_comparison_value
                } else if tokens[1].Value == ">" {
                    return real_row_value > real_comparison_value
                } else if tokens[1].Value == "<" {
                    return real_row_value < real_comparison_value
                } else if tokens[1].Value == "<=" {
                    return real_row_value <= real_comparison_value
                } else if tokens[1].Value == ">=" {
                    return real_row_value >= real_comparison_value
                } else {
                    fmt.Println("Unknown comparison operator: " + tokens[1].Value)
                    return false
                }
    		}
		}
	} else if column_types[found_column_id] == 3 {
		real_row_value := strings.ToUpper(row_value)

		if real_row_value != "T" && real_row_value != "F" {
			fmt.Println("Unable to convert row value to boolean; must either be T or F:", real_row_value)
			return false
		} else {
    		real_comparison_value := strings.ToUpper(comparison_value)

    		if real_row_value != "T" && real_row_value != "F" {
    			fmt.Println("Unable to convert row value to boolean; must either be T or F:", real_comparison_value)
    			return false
    		} else {
                if tokens[1].Value == "=" || tokens[1].Value == "==" {
                    return real_row_value == real_comparison_value
                } else if tokens[1].Value == "!=" {
                    return real_row_value != real_comparison_value
                } else {
                    fmt.Println("Unknown comparison operator: " + tokens[1].Value)
                    return false
                }
    		}
		}
	} else if column_types[found_column_id] == 4 {
        if tokens[1].Value == "=" || tokens[1].Value == "==" {
            return row_value == comparison_value
        } else if tokens[1].Value == "!=" {
            return row_value != comparison_value
        } else {
            fmt.Println("Unknown comparison operator: " + tokens[1].Value)
            return false
        }
	}

    return false
}

func TableSearch(query string, filename string) {
	fmt.Println("Call to search with:", filename, "and query", query)

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

	fmt.Println("Parsing query: `" + query + "`")
	var tokens []token
	tokens, err = tokenizeQuery(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	   fmt.Println("\nParsed tokens:")
	   for i := range(tokens) {
	       fmt.Println(":::: Token", i, "::::")
	       fmt.Println("\tType:", tokens[i].Type)
	       fmt.Println("\tValue: `" + tokens[i].Value + "`")
	   }
	*/

	var relations []rtoken
	relations, err = relationizeTokens(tokens)

	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	   fmt.Println("\nParsed relations:")
	   for i := range(relations) {
	       fmt.Println(":::: Relation", i, "::::")
	       fmt.Println("\tType:", relations[i].Type)
	       for j := range(relations[i].Value) {
	           fmt.Println("\t:::: Token", j, "::::")
	           fmt.Println("\t\tType:", relations[i].Value[j].Type)
	           fmt.Println("\t\tValue: `" + relations[i].Value[j].Value + "`")
	       }
	   }
	*/

	err = validateRelations(relations, attribute_names, attribute_types)
	if err != nil {
		fmt.Println(err)
		return
	}

	var tree evalTree
	tree, err = evalTreeizeRelation(relations)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(prettyEvalTree(&tree))

	file = file[1:]
	for i := range file {
		var line string = file[i][1 : len(file[i])-1]
		var values []string = strings.Split(line, "|")

		if evaluateTreeForRow(tree, attribute_names, attribute_types, values) {
			fmt.Println("Found match:", i)
		}
	}

	fmt.Println("Successfully searched in table `", filename, "`!")
}
