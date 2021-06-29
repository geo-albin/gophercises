package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Problems list of problems
type Problems struct {
	q     string
	a     string
	asked bool
}

func main() {
	rand.Seed(time.Now().UnixNano())
	probFile := flag.String("file", "problems.csv", "csv file containing the problems")
	timerVal := flag.Int("timer", 5, "timer value in seconds")
	flag.Parse()
	problems, err := readProblemFile(*probFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var correct int32
	var totalQ int
	ans := make(chan string)
	e := make(chan error)
	for {
		if totalQ >= len(problems) {
			break
		}
		n := rand.Intn(len(problems))
		problem := problems[n]
		if problem.asked {
			continue
		}
		totalQ++
		problem.asked = true
		go func(p Problems) {
			fmt.Printf("%s : ", p.q)
			reader := bufio.NewReader(os.Stdin)
			a, err := reader.ReadString('\n')
			if err != nil {
				e <- err
			} else {
				ans <- strings.TrimSuffix(a, "\n")
			}
		}(problem)
		select {
		case a := <-ans:
			if a != problem.a {
				fmt.Printf("Answer is wrong. Correct answer is %s, your answer %s\n", problem.a, a)
			} else {
				correct++
			}
		case err := <-e:
			fmt.Println(err.Error())
		case <-time.After(time.Duration(*timerVal) * time.Second):
			fmt.Println("timeout!!")
		}
	}
	close(ans)
	close(e)
	fmt.Printf("Total correct answers = %d / %d\n", correct, len(problems))
}

func readProblemFile(profFile string) ([]Problems, error) {
	result := make([]Problems, 0)
	recordFile, err := os.Open(profFile)
	defer recordFile.Close()
	if err != nil {
		return result, err
	}
	reader := csv.NewReader(recordFile)
	r, err := reader.ReadAll()
	for _, q := range r {
		p := Problems{
			q:     q[0],
			a:     q[1],
			asked: false,
		}
		result = append(result, p)
	}
	return result, nil
}
