package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	filename := flag.String("csv", "problem.csv", "csv file containing the questions and answers")
	timelimit := flag.Int("timelimit", 30, "time limit to answer all the questions")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("File with the name '%s' cannot be found", *filename)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)
	problems, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Malformed file '%s'", *filename)
		os.Exit(1)
	}
	var correct uint32 = 0
	answerc := make(chan bool)
	fmt.Printf("A total of %d questions will be asked and you have %d seconds to answer all of them. Press enter if you're ready", len(problems), timelimit)
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(*timelimit) * time.Second)
problemloop:
	for _, problem := range problems {
		go ask(problem, answerc)
		select {
		case _ = <-timer.C:
			fmt.Printf("\nTime is up!!!")
			break problemloop
		case answer := <-answerc:
			if answer == true {
				correct++
			}
			break
		}
	}
	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
}

func ask(problem []string, c chan bool) {
	fmt.Printf("The result of %s is: ", problem[0])
	var answer string
	fmt.Scanln(&answer)
	if answer == problem[1] {
		c <- true
	} else {
		c <- false
	}

}

func countdown(timer int, tc chan bool) {
	time.Sleep(time.Second * 5)
	tc <- true
}
