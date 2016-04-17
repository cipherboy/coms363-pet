package main

import (
    "fmt"
    "os"
    "strconv"
)

func TableCreate(attributes []string, types []int, filename string) {
    fmt.Println("Call to create with:", filename)
    var header string
    header = "[" + strconv.Itoa(len(attributes)) + "]"

    for i := range(attributes) {
        header += "[" + attributes[i] + ":" + strconv.Itoa(types[i]) + "]"
    }
    header += "[0]\n"

    if _, err := os.Stat(filename); err == nil {
        fmt.Println("Error: file `", filename, "` already exists... Refusing to overwrite.")
        return
    }

    f, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer f.Close()

    wl, err := f.Write([]byte(header))
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    if wl != len(header) {
        fmt.Println("Error writing file: wrote", wl, "bytes but expected to write", len(header))
        return
    }

    fmt.Println("Successfully created table `", filename, "`!")
}
