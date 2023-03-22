package parser

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Link struct with fields Text and Href
type Link struct {
	Href string
	Text string
}

func Parse(filename string) ([]Link, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()

	// create buffer
	b := make([]byte, 256) // chunk size
	var wg sync.WaitGroup
	var extractedLinks []Link
	var fileContent string
	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		fileContent = fmt.Sprintf("%s%s", fileContent, string(b[:readTotal]))
	}
	tkn := html.NewTokenizer(strings.NewReader(fileContent))
	wg.Add(1)
	go func() {
		extractedLinks = append(extractedLinks, parseHTML(tkn)...)
		wg.Done()
	}()
	wg.Wait()
	log.Printf("Extracted Links\n%+v\n", extractedLinks)
	return extractedLinks, nil
}

func parseHTML(tkn *html.Tokenizer) []Link {
	var vals []Link
	var isLink bool
	var oneLink Link
Loop:
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			err := tkn.Err()
        	if err == io.EOF {
        	    //end of the file, break out of the loop
        	    break Loop
        	}
			if err != nil {
				log.Printf("\nhtml error token %v\n", err)
			}
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
		case tt == html.EndTagToken:
			if isLink {
				isLink = false
			}
		}
	}
	return vals
}
