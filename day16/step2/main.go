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
		for i := 0; i < 10000; i++ {
			for _, numberChar := range signalString {
				signal = append(signal, int(numberChar-'0'))
			}
		}
	}

	firstPart := make([]int, 7)
	copy(firstPart, signal[0:7])

	offset := 0

	for _, digit := range firstPart {
		offset *= 10
		offset += digit
	}

	for i := 0; i < 100; i++ {
		signal = calcPhase(signal, []int{0, 1, 0, -1})
	}

	for _, digit := range signal[offset : offset+8] {
		fmt.Print(digit)
	}
	fmt.Print("\n")
}

func calcPhase(signal, pattern []int) []int {
	newOut := make([]int, len(signal))

	sum := 0

	for i := len(signal) - 1; i >= 0; i-- {

		sum += signal[i]
		newOut[i] = sum % 10
	}

	return newOut
}
