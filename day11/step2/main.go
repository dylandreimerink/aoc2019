package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const DEBUG = false

func dumpIfDebug(prefix string, value interface{}) {
	if DEBUG {
		fmt.Printf("DEBUG - %s: %v\n", prefix, value)
	}
}

func main() {
	content, err := ioutil.ReadFile("day11/input.txt")
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

	const canvasSize = 100

	canvas := make([][]Panel, canvasSize)
	for i := 0; i < canvasSize; i++ {
		canvas[i] = make([]Panel, canvasSize)
	}

	//0 = up, 1 = right, 2 = down, 3 = left
	robotDir := 0
	robotX := (canvasSize / 2)
	robotY := (canvasSize / 2)

	//Make the starting panel white
	canvas[robotY][robotX].color = 1

	//0 = paint mode, 1 = move mode
	outputMode := 0

	comp := &IntComputer{
		memory: initialMemory,
		inputCallback: func() string {
			return strconv.Itoa(canvas[robotY][robotX].color)
		},
		outputCallback: func(out int) {
			if outputMode == 0 {
				canvas[robotY][robotX].color = out
				canvas[robotY][robotX].painted++
				outputMode = 1
			} else {
				if out == 0 {
					robotDir--
					if robotDir < 0 {
						robotDir = 3
					}
				}
				if out == 1 {
					robotDir++
					if robotDir > 3 {
						robotDir = 0
					}
				}

				switch robotDir {
				case 0:
					robotX--
				case 1:
					robotY--
				case 2:
					robotX++
				case 3:
					robotY++
				}

				outputMode = 0
			}
		},
	}

	err = comp.Exec()
	if err != nil {
		panic(err)
	}

	for _, row := range canvas {
		for _, panel := range row {
			if panel.color == 1 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

type Panel struct {
	color   int
	painted int
}

type IntComputer struct {
	memory         []int
	relativeBase   int
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

		dumpIfDebug("opcode", opcode)
		dumpIfDebug("paramModes", fmt.Sprintf("first: %d, sec: %d, third: %d", firstParameterMode, secondParameterMode, thirdParameterMode))

		switch opcode {
		case 1:
			//Addition

			inputPointerA := comp.readValueFromMemory(instructionPointer+1, POSITION_MODE)
			InputPointerB := comp.readValueFromMemory(instructionPointer+2, POSITION_MODE)
			outputPointer := comp.readValueFromMemory(instructionPointer+3, POSITION_MODE)
			instructionPointer += 4

			result := comp.readValueFromMemory(inputPointerA, firstParameterMode) + comp.readValueFromMemory(InputPointerB, secondParameterMode)
			comp.writeValueToMemory(outputPointer, result, thirdParameterMode)
			break
		case 2:
			//Multiply

			inputPointerA := comp.readValueFromMemory(instructionPointer+1, POSITION_MODE)
			InputPointerB := comp.readValueFromMemory(instructionPointer+2, POSITION_MODE)
			outputPointer := comp.readValueFromMemory(instructionPointer+3, POSITION_MODE)
			instructionPointer += 4

			result := comp.readValueFromMemory(inputPointerA, firstParameterMode) * comp.readValueFromMemory(InputPointerB, secondParameterMode)
			comp.writeValueToMemory(outputPointer, result, thirdParameterMode)
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

			writeAddress := comp.readValueFromMemory(instructionPointer+1, POSITION_MODE)
			comp.writeValueToMemory(writeAddress, number, firstParameterMode)

			instructionPointer += 2
			break
		case 4:
			//Print

			output := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)

			if comp.outputCallback == nil {
				fmt.Printf("Output: %d\n", output)
			} else {
				comp.outputCallback(output)
			}

			instructionPointer += 2
			break
		case 5:
			//Jump if true

			cmpValue := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := comp.readValueFromMemory(comp.memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue != 0 {
				instructionPointer = newInstructionPointer
			}
			break
		case 6:
			//Jump if false

			cmpValue := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)
			newInstructionPointer := comp.readValueFromMemory(comp.memory[instructionPointer+2], secondParameterMode)

			instructionPointer += 3

			if cmpValue == 0 {
				instructionPointer = newInstructionPointer
			}
			break
		case 7:
			//Less then

			firstValue := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)
			secondValue := comp.readValueFromMemory(comp.memory[instructionPointer+2], secondParameterMode)
			resultPointer := comp.readValueFromMemory(instructionPointer+3, POSITION_MODE)

			instructionPointer += 4

			comp.writeValueToMemory(resultPointer, 0, thirdParameterMode)
			if firstValue < secondValue {
				comp.writeValueToMemory(resultPointer, 1, thirdParameterMode)
			}
			break
		case 8:
			//equals

			firstValue := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)
			secondValue := comp.readValueFromMemory(comp.memory[instructionPointer+2], secondParameterMode)
			resultPointer := comp.readValueFromMemory(instructionPointer+3, POSITION_MODE)

			instructionPointer += 4

			comp.writeValueToMemory(resultPointer, 0, thirdParameterMode)
			if firstValue == secondValue {
				comp.writeValueToMemory(resultPointer, 1, thirdParameterMode)
			}
			break
		case 9:
			//Adjust relative base

			modValue := comp.readValueFromMemory(comp.memory[instructionPointer+1], firstParameterMode)

			comp.relativeBase += modValue

			instructionPointer += 2
			break
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

const (
	POSITION_MODE int = iota
	DIRECT_MODE
	RELATIVE_MODE
)

func (comp *IntComputer) readValueFromMemory(position, mode int) int {
	var address int

	switch mode {
	case POSITION_MODE:
		address = position

	case DIRECT_MODE:
		return position

	case RELATIVE_MODE:
		address = comp.relativeBase + position

	default:
		panic(fmt.Errorf("Unknown memory access mode: %d", mode))
	}

	dumpIfDebug("Read from memory", fmt.Sprintf("addr: %d, mode: %d", address, mode))

	if address < 0 {
		panic("Can't access a negative memory address")
	}

	if address >= len(comp.memory) {
		extraMemory := address + 16 - len(comp.memory) // +16 for a little extra memory so we don't allocate all the time
		comp.memory = append(comp.memory, make([]int, extraMemory)...)
	}

	return comp.memory[address]
}

func (comp *IntComputer) writeValueToMemory(position, value, mode int) {
	var address int

	switch mode {
	case POSITION_MODE:
		address = position

	case DIRECT_MODE:
		panic("Can't write to memory in direct mode")

	case RELATIVE_MODE:
		address = comp.relativeBase + position

	default:
		panic(fmt.Errorf("Unknown memory access mode: %d", mode))
	}

	dumpIfDebug("Write to memory", fmt.Sprintf("addr: %d, val: %d, mode: %d", address, value, mode))

	if address < 0 {
		panic("Can't access a negative memory address")
	}

	if address >= len(comp.memory) {
		extraMemory := address + 16 - len(comp.memory) // +16 for a little extra memory so we don't allocate all the time
		comp.memory = append(comp.memory, make([]int, extraMemory)...)
	}

	comp.memory[address] = value
}
