package main

import (
	"fmt"
	"html-link-parser/parser"
)

const htmlParserType = "html"

func main() {
	fmt.Println("html-link-parser")
	files := []string {"example1.html", "example2.html", "example3.html", "example4.html"}
	for _, v := range files {
		parser.Parse(htmlParserType, v)
	}
}
