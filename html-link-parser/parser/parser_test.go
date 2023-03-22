package parser

import (
	"fmt"
	// "io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		input string
		want  []Link
	}{
		// success - valid html with 1 Link
		{
			input: `<html>
			<body>
			  <h1>Hello!</h1>
			  <a href="/other-page">A link to another page</a>
			</body>
			</html>`,
			want: []Link{
				{Href: "/other-page", Text: "A link to another page" },
			},
		},
		// incomplete <a> tag - only ignores tag that are without Text data
		{
			input: `<html>
			<body>
			  <h1>Hello!</h1>
			  <a href="/other-page">`,
			want: nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Parse=%d", i), func(t *testing.T) {
			file, err := os.CreateTemp(dir, "test_html_")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("File name: %s\n", file.Name())
			defer os.Remove(file.Name())
			_, err = file.WriteString(tc.input)
			if err != nil {
				t.Log("Failed to write test data to test file")
				t.Fail()
			}
			got, err := Parse(file.Name())
			if err != nil {
				t.Log("Failed to parse test data")
				t.Fail()
				return
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Logf("\ngot: %+v\nwanted: %+v\n", got, tc.want)
				t.Fail()
				return
			}
			t.Log("Success")
		})
	}
}