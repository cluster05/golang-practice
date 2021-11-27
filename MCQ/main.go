package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type QueAnswer struct {
	Quetion string
	Answer  string
}

func parseCSV(filename string) ([]QueAnswer, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	queAnswer := make([]QueAnswer, len(records))

	for i, v := range records {
		queAnswer[i] = QueAnswer{
			Quetion: v[0],
			Answer:  strings.TrimSpace(v[1]),
		}
	}

	return queAnswer, nil

}

func exitGame(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {

	filename := flag.String("csv", "problem.csv", "a csv file in the format of 'question,answer")

	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	queAnswers, err := parseCSV(*filename)

	if err != nil {
		fmt.Println("Error : ", err)
		return
	}

	totalCorrectAnswer := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemLoop:
	for i, traverse := range queAnswers {

		fmt.Printf("# Quetion : %d >> %s = ", i+1, traverse.Quetion)

		answerCh := make(chan string)

		go func() {

			var userAnswer string
			fmt.Scanf("%s\n", &userAnswer)
			answerCh <- userAnswer

		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case userAnswer := <-answerCh:
			if userAnswer == traverse.Answer {
				totalCorrectAnswer++
			}
		}
	}

	exitGame(fmt.Sprintf("Your Score is : %v / %v", totalCorrectAnswer, len(queAnswers)))

}
