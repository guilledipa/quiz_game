// https://github.com/gophercises/quiz
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	quiz        = flag.String("quiz", "./problems.csv", "Quiz filename.")
	gameTime    = flag.Duration("time", 10*time.Second, "Timer.")
	shuffleQuiz = flag.Bool("shuffle_quiz", false, "Randomize quiz.")
	correct     int
	ch          chan string
)

func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userAnswer := scanner.Text()
	userAnswer = strings.TrimSpace(userAnswer)
	userAnswer = strings.ToLower(userAnswer)
	return userAnswer
}

func shuffle(qSlice [][]string) [][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(qSlice), func(i, j int) {
		qSlice[i], qSlice[j] = qSlice[j], qSlice[i]
	})
	return qSlice
}

func askQuestions() (chan bool, int) {
	ch := make(chan bool)
	data, err := ioutil.ReadFile(*quiz)
	if err != nil {
		log.Fatal(err)
	}
	// You can assume that quizzes will be relatively
	// short (< 100 questions) and will have single word/number answers.
	// load all in memory.
	questions, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if *shuffleQuiz {
		shuffle(questions)
	}

	go func() {
		for _, q := range questions {
			fmt.Printf("%v: ", q[0])
			userAnswer := getInput()
			ch <- q[1] == userAnswer
		}
		close(ch)
	}()
	return ch, len(questions)
}

func main() {
	flag.Parse()
	rChan, totalQuestions := askQuestions()
	// I need this because cant creak outer loops from within case statement
	// https://golang.org/ref/spec#Break_statements
OuterLoop:
	for {
		select {
		case r, ok := <-rChan:
			if !ok {
				// All questions answered.
				break OuterLoop
			}
			if r { // Ignore wrong (false) responses
				correct++
			}
		case <-time.After(*gameTime):
			fmt.Println("Timeout!")
			break OuterLoop
		}
	}
	fmt.Printf("Got %d/%d\n", correct, totalQuestions)
}
