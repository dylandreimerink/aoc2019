package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("day16/input.txt")
	if err != nil {
		panic(err)
	}

	signal := []int{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		signalString := scanner.Text()
		for _, numberChar := range signalString {
			signal = append(signal, int(numberChar-'0'))
		}
	}

	for i := 0; i < 100; i++ {
		signal = calcPhase(signal, []int{0, 1, 0, -1})
	}

	for _, digit := range signal {
		fmt.Print(digit)
	}
	fmt.Print("\n")
}

func calcPhase(signal, pattern []int) []int {
	newOut := make([]int, len(signal))

	for i := 0; i < len(signal); i++ {
		sum := 0
		for ii := 0; ii < len(signal); ii++ {
			sum += signal[ii] * pattern[(ii+1)/(i+1)%len(pattern)]
		}
		newOut[i] = sum
	}

	for i := 0; i < len(newOut); i++ {
		newOut[i] = newOut[i] % 10
		if newOut[i] < 0 {
			newOut[i] = -newOut[i]
		}
	}

	return newOut
}
