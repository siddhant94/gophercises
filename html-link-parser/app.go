package main

import (
	"fmt"
	"html-link-parser/parser"
)


func main() {
	fmt.Println("html-link-parser")
	files := []string {"example1.html", "example2.html", "example3.html", "example4.html"}
	for _, v := range files {
		_, _ = parser.Parse(v)
	}
}
