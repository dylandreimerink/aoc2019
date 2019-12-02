package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("day2/input.txt")
	if err != nil {
		panic(err)
	}

	initialMemory := []int{}

	for _, value := range strings.Split(string(content), ",") {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}

		initialMemory = append(initialMemory, intValue)
	}

	//I suspect this value changes for everyone
	const wantedResult = 19690720

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			memory := make([]int, len(initialMemory))

			//Important since a go slice is modified if passed by value without cloning
			copy(memory, initialMemory)

			memory[1] = noun
			memory[2] = verb

			executedMemory, err := execIntComputer(memory)
			if err != nil {
				fmt.Printf("Noun = %d and verb = %d resulted in an error: %v\n", noun, verb, err)
			}

			if executedMemory[0] == wantedResult {
				fmt.Printf("Got wanted result when noun = %d and verb = %d, 100 * noun + verb = %d\n", noun, verb, 100*noun+verb)
				return
			}
		}
	}

	fmt.Printf("Did not find result\n")
}

func execIntComputer(memory []int) ([]int, error) {

	instructionPointer := 0

	haled := false
	for !haled {
		opcode := memory[instructionPointer]

		switch opcode {
		case 1:
			inputPointerA := memory[instructionPointer+1]
			InputPointerB := memory[instructionPointer+2]
			outputPointer := memory[instructionPointer+3]
			instructionPointer += 4

			memory[outputPointer] = memory[inputPointerA] + memory[InputPointerB]
			break
		case 2:
			inputPointerA := memory[instructionPointer+1]
			InputPointerB := memory[instructionPointer+2]
			outputPointer := memory[instructionPointer+3]
			instructionPointer += 4

			memory[outputPointer] = memory[inputPointerA] * memory[InputPointerB]
			break
		case 99:
			haled = true
			break
		default:
			return memory, fmt.Errorf("Unknown opcode: %d\n", opcode)
		}
	}

	return memory, nil
}
