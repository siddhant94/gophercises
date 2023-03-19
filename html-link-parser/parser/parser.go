package parser

import (
	// "bufio"
	// "bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
	"sync"
)

// Link struct with fields Text and Href
type Link struct {
	Href string
	Text string
}

func Parse(parseType, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	// create buffer
	b := make([]byte, 256) // chunk size
	var wg sync.WaitGroup

	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fileContent := string(b[:readTotal])
		// print content from buffer
		fmt.Println("Printing from buffer")
		fmt.Println(fileContent)
		tkn := html.NewTokenizer(strings.NewReader(string(b[:readTotal])))
		wg.Add(1)
		go func() {
			fmt.Printf("After Parse: %+v\n", parseHTML(tkn))
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func parseHTML(tkn *html.Tokenizer) []Link {
	var vals []Link
	var isLink bool
	var oneLink Link
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return vals
		case tt == html.StartTagToken:
			t := tkn.Token()
			isLink = t.Data == "a"
			if !isLink {
				break
			}
			for _, a := range t.Attr {
				if a.Key == "href" {
					oneLink.Href = a.Val
					break
				}
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if isLink {
				oneLink.Text = t.Data
				vals = append(vals, oneLink)
			}
			// isLink = false
		case tt == html.EndTagToken:
			isLink = false
		}
	}
}
