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

const DEBUG = false

func dumpIfDebug(prefix string, value interface{}) {
	if DEBUG {
		fmt.Printf("DEBUG - %s: %v\n", prefix, value)
	}
}

func main() {
	content, err := ioutil.ReadFile("day17/input.txt")
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

	//Wake up robot
	initialMemory[0] = 2

	const size = 60

	exterior := [][]int{}
	for i := 0; i < size; i++ {
		exterior = append(exterior, make([]int, size))
	}

	mainPathTodo := true
	mainPath := "A,C,A,B,A,C,B,C,B,C\n"
	mainPathOffset := 0

	// subRoutineCalc := true
	subRoutineTodo := true
	subRoutine := "L,10,R,8,L,6,R,6\nR,8,L,6,L,10,L,10\nL,8,L,8,R,8\n"
	subRoutineOffset := 0

	displayQ := 0

	var comp *IntComputer
	comp = &IntComputer{
		memory: initialMemory,
		inputCallback: func() string {
			if mainPathTodo {
				ret := mainPath[mainPathOffset]
				mainPathOffset++
				if mainPathOffset >= len(mainPath) {
					mainPathTodo = false
				}
				return strconv.Itoa(int(ret))
			}

			if subRoutineTodo {
				ret := subRoutine[subRoutineOffset]
				subRoutineOffset++
				if subRoutineOffset >= len(subRoutine) {
					subRoutineTodo = false
				}
				return strconv.Itoa(int(ret))
			}

			if displayQ == 0 {
				displayQ++
				return strconv.Itoa(int('n'))
			}

			if displayQ == 1 {
				displayQ++
				return strconv.Itoa(10) //Newline
			}

			spew.Dump("UHMMM")

			return "0"
		},
		outputCallback: func(out int) {

			if out > 255 {
				fmt.Println(out)
			} else {
				fmt.Print(string(rune(out)))
			}
		},
	}

	err = comp.Exec()
	if err != nil {
		panic(err)
	}
}

func getDroidPath(exterior [][]int, droidX, droidY, droidDir int) []string {
	path := []string{}

	rotL := func() {
		droidDir--
		if droidDir < 0 {
			droidDir = 3
		}
	}

	rotR := func() {
		droidDir++
		if droidDir > 3 {
			droidDir = 0
		}
	}

	done := false
	for !done {
		stepLength := 0
		switch droidDir {
		case 0:
			for {
				if droidX-1 >= 0 && exterior[droidY][droidX-1] == '#' {
					droidX--
					stepLength++
					continue
				}

				if exterior[droidY+1][droidX] == '#' {
					path = append(path, strconv.Itoa(stepLength), "L")
					rotL()
					break
				}

				if droidY-1 >= 0 && exterior[droidY-1][droidX] == '#' {
					path = append(path, strconv.Itoa(stepLength), "R")
					rotR()
					break
				}

				done = true
				break
			}
		case 1:
			for {
				if droidY-1 >= 0 && exterior[droidY-1][droidX] == '#' {
					droidY--
					stepLength++
					continue
				}

				if droidX-1 >= 0 && exterior[droidY][droidX-1] == '#' {
					path = append(path, strconv.Itoa(stepLength), "L")
					rotL()
					break
				}

				if exterior[droidY][droidX+1] == '#' {
					path = append(path, strconv.Itoa(stepLength), "R")
					rotR()
					break
				}

				done = true
				break
			}
		case 2:
			for {
				if exterior[droidY][droidX+1] == '#' {
					droidX++
					stepLength++
					continue
				}

				if droidY-1 >= 0 && exterior[droidY-1][droidX] == '#' {
					path = append(path, strconv.Itoa(stepLength), "L")
					rotL()
					break
				}

				if exterior[droidY+1][droidX] == '#' {
					path = append(path, strconv.Itoa(stepLength), "R")
					rotR()
					break
				}

				done = true
				break
			}
		case 3:
			for {
				if exterior[droidY+1][droidX] == '#' {
					droidY++
					stepLength++
					continue
				}

				if exterior[droidY][droidX+1] == '#' {
					path = append(path, strconv.Itoa(stepLength), "L")
					rotL()
					break
				}

				if droidX-1 >= 0 && exterior[droidY][droidX-1] == '#' {
					path = append(path, strconv.Itoa(stepLength), "R")
					rotR()
					break
				}

				done = true
				break
			}
		}
	}

	return path
}

func renderExterior(exterior [][]int) {
	for _, row := range exterior {
		for _, column := range row {
			fmt.Print(string(rune(column)))
		}
		fmt.Print("\n")
	}
}

type IntComputer struct {
	memory         []int
	relativeBase   int
	inputCallback  func() string
	outputCallback func(int)
	extHalt        bool
}

func (comp *IntComputer) Exec() error {

	instructionPointer := 0

	haled := false
	for !haled && !comp.extHalt {
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
