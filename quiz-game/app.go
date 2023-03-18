package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var filename string
var timer int
var shuffle bool

type record struct {
	question string
	answer   string
}

func init() {
	const (
		defaultFilename = "problems.csv"
		defaultTimer    = 30
	)
	flag.StringVar(&filename, "filename", defaultFilename, "a csv file in the format of 'question,answer'")
	flag.IntVar(&timer, "timer", defaultTimer, "timeout per question (in seconds)")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle questions")
}

func main() {
	flag.Parse()
	// open CSV file
	fd, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", filename))
	}
	defer fd.Close()
	lines := readCSV(fd)
	records := parseRecords(lines, shuffle)
	timeout_seconds := time.Duration(timer) * time.Second
	score := quiz(*records, timeout_seconds, shuffle, os.Stdin)
	log.Printf("Total Questions: %d\t\tCorrectly Answered: %d\n", len(*records), score)
}

func readCSV(fd io.Reader) [][]string {
	// read CSV file with ReadAll as initial datasize expected is ~100 rows.
	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()
	if err != nil {
		log.Println(err)
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", filename))
	}
	return records
}

func parseRecords(lines [][]string, shuffle bool) *[]record {
	ret := make([]record, len(lines))
	for i, line := range lines {
		if len(line) < 2 {
			// discard non valid records/lines
			continue
		}
		ret[i] = record{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	if shuffle {
		// seed to get different result every time
		rand.Seed(time.Now().UnixNano())
		time.Sleep(1 * time.Second)
		rand.Shuffle(len(ret), func(i, j int) { ret[i], ret[j] = ret[j], ret[i] })
	}
	return &ret
}

func quiz(records []record, timeout_seconds time.Duration, shuffle bool, reader io.Reader) int {
	log.Println("Starting Quiz ....")
	score := 0
	answersCh := make(chan string)
	go readAnswerInput(answersCh, reader)
	for i, rec := range records {
		fmt.Printf("Problem #%d: %s = ", i+1, rec.question)
		select {
		case answer := <-answersCh:
			if answer == rec.answer {
				fmt.Println("Correct")
				score += 1
			} else {
				fmt.Println("Oops! Incorrect")
			}
		case <-time.After(timeout_seconds):
			fmt.Println("\n Time is over!")
		}
	}
	log.Println("Finished Quiz")
	return score
}

func readAnswerInput(answersCh chan<- string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	var answer string
	for scanner.Scan() {
		answer = scanner.Text()
		answersCh <- strings.TrimSpace(answer)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("reading standard input: %v\n", err)
	}
}

func exit(msg string) {
	log.Println(msg)
	os.Exit(1)
}
