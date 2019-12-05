package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("day5/input.txt")
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

	_, err = execIntComputer(initialMemory)
	if err != nil {
		panic(err)
	}
}

func execIntComputer(memory []int) ([]int, error) {

	instructionPointer := 0

	haled := false
	for !haled {
		opcode := memory[instructionPointer]

		thirdParameterMode := opcode / 10000
		opcode -= 10000 * thirdParameterMode

		secondParameterMode := opcode / 1000
		opcode -= 1000 * secondParameterMode

		firstParameterMode := opcode / 100
		opcode -= 100 * firstParameterMode

		switch opcode {
		case 1:
			//Addition

			inputPointerA := memory[instructionPointer+1]
			InputPointerB := memory[instructionPointer+2]
			outputPointer := memory[instructionPointer+3]
			instructionPointer += 4

			memory[outputPointer] = readValueFromMemory(memory, inputPointerA, firstParameterMode) + readValueFromMemory(memory, InputPointerB, secondParameterMode)
			break
		case 2:
			//Multiply

			inputPointerA := memory[instructionPointer+1]
			InputPointerB := memory[instructionPointer+2]
			outputPointer := memory[instructionPointer+3]
			instructionPointer += 4

			memory[outputPointer] = readValueFromMemory(memory, inputPointerA, firstParameterMode) * readValueFromMemory(memory, InputPointerB, secondParameterMode)
			break
		case 3:
			//Input

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter input: ")

			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			number, err := strconv.Atoi(text[0 : len(text)-1])
			if err != nil {
				panic(err)
			}

			writeAddress := memory[instructionPointer+1]
			memory[writeAddress] = number

			instructionPointer += 2
		case 4:
			//Print

			output := readValueFromMemory(memory, memory[instructionPointer+1], firstParameterMode)
			fmt.Printf("Output: %d\n", output)

			instructionPointer += 2
		case 5:
			//Jump if true

			cmpValue := readValueFromMemory(memory, memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := readValueFromMemory(memory, memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue != 0 {
				instructionPointer = newInstructionPointer
			}
		case 6:
			//Jump if false

			cmpValue := readValueFromMemory(memory, memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := readValueFromMemory(memory, memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue == 0 {
				instructionPointer = newInstructionPointer
			}
		case 7:
			//Less then

			firstValue := readValueFromMemory(memory, memory[instructionPointer+1], firstParameterMode)
			secondValue := readValueFromMemory(memory, memory[instructionPointer+2], secondParameterMode)
			resultPointer := memory[instructionPointer+3]

			instructionPointer += 4

			memory[resultPointer] = 0
			if firstValue < secondValue {
				memory[resultPointer] = 1
			}
		case 8:
			//equals

			firstValue := readValueFromMemory(memory, memory[instructionPointer+1], firstParameterMode)
			secondValue := readValueFromMemory(memory, memory[instructionPointer+2], secondParameterMode)
			resultPointer := memory[instructionPointer+3]

			instructionPointer += 4

			memory[resultPointer] = 0
			if firstValue == secondValue {
				memory[resultPointer] = 1
			}

		case 99:
			//Halt

			haled = true
			break
		default:
			return memory, fmt.Errorf("Unknown opcode: %d\n", opcode)
		}
	}

	return memory, nil
}

func readValueFromMemory(memory []int, position, mode int) int {
	if mode == 1 {
		return position
	}

	return memory[position]
}
