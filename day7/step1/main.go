package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("day7/input.txt")
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

	maxOutput := 0

	for p0 := 0; p0 <= 4; p0++ {
		for p1 := 0; p1 <= 4; p1++ {
			if p1 == p0 {
				continue
			}

			for p2 := 0; p2 <= 4; p2++ {
				if p2 == p0 {
					continue
				}
				if p2 == p1 {
					continue
				}
				for p3 := 0; p3 <= 4; p3++ {
					if p3 == p0 {
						continue
					}
					if p3 == p1 {
						continue
					}
					if p3 == p2 {
						continue
					}
					for p4 := 0; p4 <= 4; p4++ {
						if p4 == p0 {
							continue
						}
						if p4 == p1 {
							continue
						}
						if p4 == p2 {
							continue
						}
						if p4 == p3 {
							continue
						}

						phasesSettings := []int{p0, p1, p2, p3, p4}

						ampInput := 0

						for _, phaseSetting := range phasesSettings {
							inputCount := 0

							comp := IntComputer{
								memory: append([]int(nil), initialMemory...),
								inputCallback: func() string {
									if inputCount == 0 {
										inputCount++
										return strconv.Itoa(phaseSetting)
									}

									return strconv.Itoa(ampInput)
								},
								outputCallback: func(output int) {
									ampInput = output
								},
							}

							err = comp.Exec()
							if err != nil {
								panic(err)
							}
						}

						if ampInput > maxOutput {
							maxOutput = ampInput
						}
					}
				}
			}
		}
	}

	spew.Dump(maxOutput)
}

type IntComputer struct {
	memory         []int
	inputCallback  func() string
	outputCallback func(int)
}

func (comp *IntComputer) Exec() error {

	instructionPointer := 0

	haled := false
	for !haled {
		opcode := comp.memory[instructionPointer]

		thirdParameterMode := opcode / 10000
		opcode -= 10000 * thirdParameterMode

		secondParameterMode := opcode / 1000
		opcode -= 1000 * secondParameterMode

		firstParameterMode := opcode / 100
		opcode -= 100 * firstParameterMode

		switch opcode {
		case 1:
			//Addition

			inputPointerA := comp.memory[instructionPointer+1]
			InputPointerB := comp.memory[instructionPointer+2]
			outputPointer := comp.memory[instructionPointer+3]
			instructionPointer += 4

			comp.memory[outputPointer] = readValueFromMemory(comp.memory, inputPointerA, firstParameterMode) + readValueFromMemory(comp.memory, InputPointerB, secondParameterMode)
			break
		case 2:
			//Multiply

			inputPointerA := comp.memory[instructionPointer+1]
			InputPointerB := comp.memory[instructionPointer+2]
			outputPointer := comp.memory[instructionPointer+3]
			instructionPointer += 4

			comp.memory[outputPointer] = readValueFromMemory(comp.memory, inputPointerA, firstParameterMode) * readValueFromMemory(comp.memory, InputPointerB, secondParameterMode)
			break
		case 3:
			//Input
			text := ""

			if comp.inputCallback == nil {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter input: ")

				var err error

				text, err = reader.ReadString('\n')
				if err != nil {
					panic(err)
				}

				text = text[0 : len(text)-1]
			} else {
				text = comp.inputCallback()
			}

			number, err := strconv.Atoi(text)
			if err != nil {
				panic(err)
			}

			writeAddress := comp.memory[instructionPointer+1]
			comp.memory[writeAddress] = number

			instructionPointer += 2
		case 4:
			//Print

			output := readValueFromMemory(comp.memory, comp.memory[instructionPointer+1], firstParameterMode)

			if comp.outputCallback == nil {
				fmt.Printf("Output: %d\n", output)
			} else {
				comp.outputCallback(output)
			}

			instructionPointer += 2
		case 5:
			//Jump if true

			cmpValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := readValueFromMemory(comp.memory, comp.memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue != 0 {
				instructionPointer = newInstructionPointer
			}
		case 6:
			//Jump if false

			cmpValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := readValueFromMemory(comp.memory, comp.memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue == 0 {
				instructionPointer = newInstructionPointer
			}
		case 7:
			//Less then

			firstValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+1], firstParameterMode)
			secondValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+2], secondParameterMode)
			resultPointer := comp.memory[instructionPointer+3]

			instructionPointer += 4

			comp.memory[resultPointer] = 0
			if firstValue < secondValue {
				comp.memory[resultPointer] = 1
			}
		case 8:
			//equals

			firstValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+1], firstParameterMode)
			secondValue := readValueFromMemory(comp.memory, comp.memory[instructionPointer+2], secondParameterMode)
			resultPointer := comp.memory[instructionPointer+3]

			instructionPointer += 4

			comp.memory[resultPointer] = 0
			if firstValue == secondValue {
				comp.memory[resultPointer] = 1
			}

		case 99:
			//Halt

			haled = true
			break
		default:
			return fmt.Errorf("Unknown opcode: %d\n", opcode)
		}
	}

	return nil
}
func readValueFromMemory(memory []int, position, mode int) int {
	if mode == 1 {
		return position
	}

	return memory[position]
}
