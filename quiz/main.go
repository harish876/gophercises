package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type Quiz struct {
	Question string
	Answer   int64
}

func GetQuizData(data [][]string) []Quiz {
	var list []Quiz
	for i, line := range data {
		if i > 0 {
			var rec Quiz
			for j, field := range line {
				if j == 0 {
					rec.Question = field
				} else if j == 1 {
					value, _ := strconv.Atoi(field)
					rec.Answer = int64(value)
				}
			}
			list = append(list, rec)
		}
	}
	return list
}

func shuffle(arr []int) {
	// Start from the last element and swap with a randomly selected element before it
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func main() {

	var filename string
	var timeout int
	var shuffle int

	fs := pflag.NewFlagSet("quiz cli", pflag.ExitOnError)
	fs.StringVarP(&filename, "file", "f", "", "Specify a file (long form: --file, short form: -f)")
	fs.IntVarP(&timeout, "timeout", "t", 10, "Specify a timeout value in seconds (long form: --timeout, short form: -t)")
	fs.IntVarP(&shuffle, "shuffle", "s", 0, "Specify a random number for shuffling as integer(long form: --shuffle, short form: -s)")

	fs.Parse(os.Args)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("File does not exist %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(os.Stdin)
	csvReader := csv.NewReader(f)
	data, _ := csvReader.ReadAll()
	updatedData := GetQuizData(data)

	correctAnswers := 0
	wrongAnswers := 0

	startTest := false
	for !startTest {
		fmt.Println("Start Your test. Answer in yes or no")
		start, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		start = strings.TrimSpace(start)
		if start == "yes" || start == "y" || start == "Yes" {
			startTest = true
		}

	}

	if startTest {
		rand.Seed(int64(shuffle)) // deprecated replace
		timer := time.NewTimer(time.Duration(timeout) * time.Second)
		done := make(chan bool)

		go func() {
			for _, val := range updatedData {
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("The Question is %s \n", val.Question)
				fmt.Print("Enter your Answer: ")

				answer, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading input:", err)
					return
				}

				answer = strings.TrimSpace(answer)

				userAnswer, err := strconv.Atoi(answer)

				if err != nil {
					fmt.Printf("Incorrect Answer format. Should be a number. This will counted as a wrong answer\n")
				}

				if userAnswer == int(val.Answer) {
					correctAnswers = correctAnswers + 1
				} else {
					wrongAnswers = wrongAnswers + 1
				}

				fmt.Printf("You entered answer as, %s!\n", answer)
			}
			done <- true
		}()

		select {
		case <-timer.C:
			fmt.Println("\n Oops, your Timer expired. Test ended")
		case <-done:
			fmt.Println("Test completed")
		}

		fmt.Printf("The number of correct answers is %d, and the number of wrong answers is %d, number of unattempted questions %d", correctAnswers, wrongAnswers, len(updatedData)-correctAnswers-wrongAnswers)
	}
}
