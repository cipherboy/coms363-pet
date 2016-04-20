# Project 3

    u: Alexander Scheel
    e: <scheel@iastate.edu>
    l: BSD: 2-clause

## Overview
This is my final project for COM S 363. It is a pet database written in Golang.

It supports the following operations:

- creating a table
- displaying a table's header
- inserting into a table
- displaying a table row
- deleting a table
- searching a table on a particular column

This is provided via an interactive prompt with readline support. Tested on
Mac OS X and Linux, using Go 1.6.1. Interactive commands can be listed via
the built-in help text. Type 'help' to get started.  

## Building
To build and run, make sure go = 1.6.1 is installed. Then execute the following:

    cd ./pet
    go get
    go build
    ./pet
