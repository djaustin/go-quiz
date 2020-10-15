package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Question struct {
	question string
	answer   string
}

func main() {
	fileName := flag.String("f", "", "name of the quiz question csv file")
	duration := flag.Int("d", 0, "the duration of the quiz in seconds - must be greater than 1")

	flag.Parse()
	if *fileName == "" || *duration < 1 {
		flag.Usage()
		return
	}

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalln("cannot open question file", err.Error())
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	questions := []Question{}
	for _, line := range lines {
		if len(line) != 2 {
			continue
		}
		question := Question{
			question: line[0],
			answer:   line[1],
		}
		questions = append(questions, question)
	}

	var score int
	c := make(chan bool)
	go func() {
		for i, question := range questions {
			var answerGiven string
			fmt.Printf("%d) %s:  ", i+1, question.question)
			fmt.Scanln(&answerGiven)
			if answerGiven == question.answer {
				score++
			}
		}
		c <- true
	}()

	timer := time.NewTimer(time.Second * time.Duration(*duration))

	wait := true

	for wait {
		select {
		case <-timer.C:
			wait = false
		case <-c:
			wait = false
		default:
			continue
		}
	}

	fmt.Printf("\nYou scored %d/%d\n", score, len(questions))
}
