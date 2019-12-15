package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const DEBUG = false

func dumpIfDebug(prefix string, value interface{}) {
	if DEBUG {
		fmt.Printf("DEBUG - %s: %v\n", prefix, value)
	}
}

func main() {
	content, err := ioutil.ReadFile("day15/input.txt")
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

	mapSize := 50

	shipMap := [][]int{}
	for i := 0; i < mapSize; i++ {
		shipMap = append(shipMap, make([]int, mapSize))
	}

	droidX := 25
	droidY := 25

	nextDir := 1

	moves := 0

	var comp *IntComputer
	comp = &IntComputer{
		memory: initialMemory,
		inputCallback: func() string {
			return strconv.Itoa(nextDir)
		},
		outputCallback: func(out int) {
			if out == 0 {
				switch nextDir {
				case 1:
					shipMap[droidY-1][droidX] = 1
				case 2:
					shipMap[droidY+1][droidX] = 1
				case 3:
					shipMap[droidY][droidX-1] = 1
				case 4:
					shipMap[droidY][droidX+1] = 1
				}
			} else if out == 1 || out == 2 {
				switch nextDir {
				case 1:
					droidY--
				case 2:
					droidY++
				case 3:
					droidX--
				case 4:
					droidX++
				}

				shipMap[droidY][droidX] = out + 1
			}

			moves++
			if moves%10000 == 0 {

				undiscovered := false

				for y, row := range shipMap {
					for x, column := range row {
						if column == 2 {
							if shipMap[y+1][x] == 0 || shipMap[y-1][x] == 0 ||
								shipMap[y][x+1] == 0 || shipMap[y][x-1] == 0 {
								undiscovered = true
							}
						}
					}
				}

				if !undiscovered {
					comp.extHalt = true
					minutes := BreadthFirstOxygen(shipMap, droidX, droidY)
					fmt.Printf("It takes %d minutes\n", minutes)
					return
				}

				renderMap(shipMap, droidX, droidY)
			}

			rand.Seed(time.Now().UnixNano())
			nextDir = (rand.Int() % 4) + 1
		},
	}

	err = comp.Exec()
	if err != nil {
		panic(err)
	}
}

func BreadthFirstOxygen(shipMap [][]int, dx, dy int) int {
	depth := 0

	for {
		seenAir := false
		for y, row := range shipMap {
			for x, column := range row {
				if column == 2 {
					seenAir = true
				}
				if column == 4 {
					shipMap[y][x] = 3
				}
			}
		}

		if !seenAir {
			return depth
		}

		for y, row := range shipMap {
			for x, column := range row {
				if column == 3 {
					if shipMap[y+1][x] == 2 {
						shipMap[y+1][x] = 4
					}
					if shipMap[y-1][x] == 2 {
						shipMap[y-1][x] = 4
					}
					if shipMap[y][x+1] == 2 {
						shipMap[y][x+1] = 4
					}
					if shipMap[y][x-1] == 2 {
						shipMap[y][x-1] = 4
					}
				}
			}
		}
		renderMap(shipMap, dx, dy)
		depth++
	}
}

func renderMap(shipMap [][]int, dx, dy int) {
	for y, row := range shipMap {
		if y == 0 {
			for i := 0; i < len(row); i++ {
				fmt.Print("-")
			}
			fmt.Print("\n")
		}
		for x, column := range row {
			if x == 0 {
				fmt.Print("|")
			}

			if x == dx && y == dy {
				fmt.Print("D")
				continue
			}

			if x == 25 && y == 25 {
				fmt.Print("B")
				continue
			}

			switch column {
			case 0:
				fmt.Print(" ")
			case 1:
				fmt.Print("#")
			case 2:
				fmt.Print(".")
			case 3:
				fmt.Print("O")
			case 4:
				fmt.Print("o")
			}

			if x == len(row)-1 {
				fmt.Print("|")
			}

		}

		if y == len(shipMap)-1 {
			fmt.Print("\n")
			for i := 0; i < len(row); i++ {
				fmt.Print("-")
			}
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
