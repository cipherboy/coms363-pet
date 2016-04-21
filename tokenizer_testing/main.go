package main

import (
    "fmt"
    "strconv"
    "errors"
)

func bytes_contains(needle byte, haystack []byte) int {
    for i := range(haystack) {
        if haystack[i] == needle {
            return i
        }
    }
    return -1
}

func strings_contains(needle string, haystack []string) int {
    for i := range(haystack) {
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
var unknown_token_type int = -1;
var operator_token_type int = 0;
var bareword_token_type int = 1;
var join_token_type int = 2;
var string_token_type int = 3;
var number_token_type int = 4;
var token_types_to_names map[int]string = map[int]string{-1: "unknown", 0: "operator", 1: "bareword", 2: "join", 3: "string", 4: "number"}

type token struct {
	Value   string
    Type    int
}

func tokenizeQuery(query string) ([]token, error) {
	var result []token

    var whitespace_parts []byte = []byte(" \t\n")
    var join_parts []byte = []byte("&|")
    var operator_parts []byte = []byte("><=!")
    var bareword_parts []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")
    var number_parts []byte = []byte("0123456789")
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

            current.Value = current.Value[1:len(current.Value)-1]
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
	Value   []token
    Type    int
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

func validateRelations(set []rtoken) error {
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

    for i := range(set) {
        var tokens []token = set[i].Value
        if set[i].Type == relation_rtoken_type {
            if len(tokens) != 3 {
                return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting three tokens in relation")
            }

            if tokens[0].Type != bareword_token_type {
                return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Expecting left most token to be bareword")
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

            if tokens[2].Type != number_token_type && strings_contains(tokens[1].Value, valid_string_operators) == -1 {
                return errors.New("Invalid relation (" + strconv.Itoa(i) + "): Unknown operator for strings: " + tokens[1].Value)
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

type evalTree struct {
    Join        int
    Left        *evalTree
    Right       *evalTree
    Relation    rtoken
    Value       bool
}

func evalTreeizeRelation(set []rtoken) (evalTree, error) {
    var result evalTree
    result.Left = nil
    result.Right = nil
    result.Join = -1

    for i := 0; i < len(set); i++ {

    }

    return result, nil
}

func main() {
    query := "wonderful == 2341 && other == value || something = 'Testing' && magical != 2"
    fmt.Println("Parsing query: `" + query + "`")
    var tokens []token
    tokens, err := tokenizeQuery(query)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("\nParsed tokens:")
    for i := range(tokens) {
        fmt.Println(":::: Token", i, "::::")
        fmt.Println("\tType:", tokens[i].Type)
        fmt.Println("\tValue: `" + tokens[i].Value + "`")
    }

    var relations []rtoken
    relations, err = relationizeTokens(tokens)

    if err != nil {
        fmt.Println(err)
        return
    }

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

    err = validateRelations(relations)
    if err != nil {
        fmt.Println(err)
        return
    }
}
