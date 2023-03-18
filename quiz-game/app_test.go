package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseRecords(t *testing.T) {
	tests := []struct {
		input [][]string
		want  []record
	}{
		// accepted CSV fomat (each row = `question,answer`)- success
		{[][]string{{"5+5", "10"}, {"what's 2+2", "4"}}, []record{{question: "5+5", answer: "10"}, {question: "what's 2+2", answer: "4"}}},
		// more than 3 column - success
		{[][]string{{"5+5", "10", "dummy-column"}}, []record{{question: "5+5", answer: "10"}}},
		// less than 2 column - failure
		{[][]string{{"5+5"}}, []record{{}}}, // or make([]record, 1)
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("parseRecords=%d", i), func(t *testing.T) {
			problems := parseRecords(tc.input, false)
			if reflect.DeepEqual(*problems, tc.want) {
				t.Logf("Success !")
			} else {
				t.Fatalf("got %v; want %v", problems, tc.want)
			}

		})
	}
	// test parseRecords with shuffle true
	t.Run(fmt.Sprintln("parseRecords=shuffle true"), func(t *testing.T) {
		shuffled := parseRecords(tests[0].input, true)
		// check if shuffle successful
		if reflect.DeepEqual(*shuffled, tests[0].want) {
			t.Fatalf("got %v; want- should not equal %v", *shuffled, tests[0].want)
		} else {
			t.Logf("Success !")
		}

	})

}

func TestQuiz(t *testing.T) {
	tests := []struct {
		input []record
		want  int
	}{
		{[]record{
			{"5+5", "10"},
			{"what's 6+4", "10"},
		}, 2},
		{[]record{
			{"12-2", "10"},
			{"what's 3+2", "5"},
		}, 1},
	}
	for i, tc := range tests {
		reader := strings.NewReader("10\n10\n")
		t.Run(fmt.Sprintf("quiz=%d", i), func(t *testing.T) {
			// timeout=1 so that <-timeAfter is skipped in first run
			score := quiz(tc.input, 1*time.Second, false, reader)
			if score == tc.want {
				t.Logf("Success !")
			} else {
				t.Fatalf("got %v; want %v", score, tc.want)
			}

		})
	}
}
