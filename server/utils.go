package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func dumpObject(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(b))
}

func getRandomAnswers(q Question) ([]string, int) {
	out := make([]string, IncorrectAnswerCount, IncorrectAnswerCount+1)
	copy(out, q.IncorrectAnswers)
	out = append(out, q.CorrectAnswer)
	answer := IncorrectAnswerCount

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
		if answer == i {
			answer = j
		} else if answer == j {
			answer = i
		}
	})
	return out, answer
}
