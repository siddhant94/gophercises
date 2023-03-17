package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"flag"
	"math/rand"
	"log"
	"os"
	"strings"
	"time"
)

var filename string
var timer int
var shuffle bool
func init() {
	const (
		defaultFilename = "problems.csv"
		filenameUsage   = "filename for quest-ans CSV"
		defaultTimer = 30
		timerUsage = "timeout in seconds per question"
	)
	flag.StringVar(&filename, "filename", defaultFilename, filenameUsage)
	flag.IntVar(&timer, "timer", defaultTimer, timerUsage)
	flag.BoolVar(&shuffle, "shuffle", false, "")
}

func main() {
	flag.Parse()
	// open CSV file
	fd, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	// read CSV file with ReadAll as initial datasize expected is ~100 rows.
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
  	if error != nil {
    	log.Fatalln(error)
  	}
  	fmt.Println(records)
	timeout_seconds := time.Duration(timer) * time.Second
	quiz(records, timeout_seconds, shuffle)
}

func quiz(rows [][]string, timeout_seconds time.Duration, shuffle bool) {
	log.Println("Starting Quiz ....")
	count := 0
	score := 0
	// seed to get different result every time
	rand.Seed(time.Now().UnixNano())
	num_rows := len(rows)
	userInput := make(chan string)
	quit := make(chan bool)
	go readInput(userInput, quit)
	if shuffle {
		rand.Shuffle(num_rows, func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })
	}
	for _,v := range rows {
		count += 1
		ques := v[0]
		correctAnswer := v[1]
		log.Printf("Ques-%d: %s", count, ques)
		select {
        	case inputAnswer := <-userInput:
				inputAnswer = strings.TrimSpace(inputAnswer)
				if checkAnswer(correctAnswer, inputAnswer) {
					log.Println("Correct Answer")
					score += 1
				} else {
					log.Println("Oops! Incorrect")
				}
            	
        	case <-time.After(timeout_seconds):
            	fmt.Println("\n Time is over!")
        }
	}
	// select {
    // case quit <- true:
    //     fmt.Println("sent message")
    // default:
    //     fmt.Println("no message sent")
    // }
	log.Println("Finished Quiz")
	log.Printf("\nTotal Questions: %d\t\tCorrectly Answered: %d\n", count, score)
}

func readInput(userInput chan<- string, quit <-chan bool) {
	for {
		select {
        case <-quit:
			// close(userInput)
            return
        default:
            var userAnswer string
			// _, err := fmt.Scanln(&userAnswer)
			r := bufio.NewReader(os.Stdin)
			userAnswer, err := r.ReadString('\n')
    		if err != nil {
    		    log.Println(err)
    		}
			userAnswer = strings.TrimSpace(userAnswer)
			userInput <- userAnswer
        }
	}
}

func checkAnswer(correctAnswer, inputAnswer string) bool {
	return correctAnswer == inputAnswer
}