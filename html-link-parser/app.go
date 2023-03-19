package main

import (
	"fmt"
	"html-link-parser/parser"
)

const htmlParserType = "html"

func main() {
	fmt.Println("html-link-parser")
	// n := "example1.html"
	n := "example2.html"
	parser.Parse(htmlParserType, n)
}
