package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var defPath string
	var templateFile string
	var name string

	flag.StringVar(&defPath, "p", "", "Path to load service definition from")
	flag.StringVar(&templateFile, "t", "", "Path to the template file")
	flag.StringVar(&name, "n", "", "Name of the service definition")

	flag.Parse()

	result, err := ParseDefinition(defPath)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(templateFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	io.Copy(&buf, f)
	output, err := render(string(buf.Bytes()), result, name)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
