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

	memory := []int{}

	for _, value := range strings.Split(string(content), ",") {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}

		memory = append(memory, intValue)
	}

	memory[1] = 12
	memory[2] = 2

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
			panic(fmt.Errorf("Unknown opcode: %d\n", opcode))
		}

	}

	fmt.Printf("Value at location 0 is: %d\n", memory[0])
}
