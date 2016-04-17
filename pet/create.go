package main

import (
    "fmt"
    "strconv"
)

func TableCreate(attributes []string, types []int, filename string) {
    fmt.Println("Call to create with:", filename)
    var header string
    header = "[" + strconv.Itoa(len(attributes)) + "]"

    for i := range(attributes) {
        header += "[" + attributes[i] + ":" + strconv.Itoa(types[i]) + "]"
    }
    header += "[0]"

    fmt.Println(header)
}
