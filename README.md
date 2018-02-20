# Diglet: JSON Compiler and Schema Validator
[![GoDoc](https://godoc.org/github.com/gwicks/diglet?status.svg)](https://godoc.org/github.com/gwicks/diglet)

A JSON compiler intended for use in configuration systems, it supports references to files, references to specific objects in external file, and local document query references. Also supports parenting and inheritence between JSON files, allowing a given file to have any number of parents. Also incorporates schema validation, where a schema is referenced, conforming to the JSON Schema Draft 06 spec.

# Installation

To set up diglet, you'll need to install it's dependencies using `glide`, found [here](https://github.com/Masterminds/glide)

Then, simply run

```bash
glide update
go build
```

# Usage as a command line application

To install the built-in command line implementation, assuming your system `$PATH` includes your `$GOPATH/bin` simply run

```bash
make install
diglet
```

There are two commands available in the sample command line implementation

## diglet compile

Compiles a single file, resolving references, parenting, and validating any present schemas. If not output file is specified, prints to stdout.

Usage 
```bash
diglet compile infile <output file>
```

## diglet batchfile

Takes a text file listing of input files and their respective output locations, essentially a wrapper around the compile command.

Batchfile example 

```
test/a.json out/a_done.json
test/b.json out/b_done.json
```

Usage
```bash
diglet batchfile batch.txt
```

# Usage as a library

To use this as a library in another application, the relevent package to import is `compiler`

Usage example

```go
    import (
        "fmt"
        "github.com/gwicks/diglet/compiler"
    )

    func example() {
        resultString, _ := compiler.CompileFile("test/foo.json")
        fmt.Println(resultString)
    }
```