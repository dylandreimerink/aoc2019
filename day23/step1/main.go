package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const DEBUG = false

func dumpIfDebug(prefix string, value interface{}) {
	if DEBUG {
		fmt.Printf("DEBUG - %s: %v\n", prefix, value)
	}
}

func main() {
	content, err := ioutil.ReadFile("day23/input.txt")
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

	networkMap := map[int]chan packet{}

	compMap := map[int]*IntComputer{}

	for addr := 0; addr < 50; addr++ {

		outCount := 0
		outAddr := -1
		outX := -1

		initAddr := false
		compAddr := addr

		var curIn *packet

		//Make a buffered channel with a queue size of 100
		compChan := make(chan packet, 100)

		networkMap[addr] = compChan

		memCopy := make([]int, len(initialMemory))
		copy(memCopy, initialMemory)

		comp := &IntComputer{
			memory: memCopy,
			inputCallback: func() string {
				if !initAddr {
					initAddr = true
					return strconv.Itoa(compAddr)
				}

				if curIn != nil {
					y := curIn.y
					curIn = nil
					return strconv.Itoa(y)
				}

				select {
				case in := <-compChan:
					curIn = &in
					return strconv.Itoa(curIn.x)
				default:
					return strconv.Itoa(-1)
				}
			},
			outputCallback: func(out int) {
				switch outCount {
				case 0:
					outAddr = out
					outCount++
					break
				case 1:
					outX = out
					outCount++
					break
				case 2:
					if outAddr == 255 {
						fmt.Printf("Y of packet to address 255 is: %d\n", out)
						os.Exit(0)
					}

					fmt.Printf("From %d to %d - x: %d, y: %d\n", compAddr, outAddr, outX, out)

					networkMap[outAddr] <- packet{x: outX, y: out}
					outCount = 0
				}
			},
		}

		compMap[addr] = comp
	}

	wg := &sync.WaitGroup{}

	for addr := 0; addr < 50; addr++ {
		go func(address int) {
			wg.Add(1)
			defer wg.Done()
			err := compMap[address].Exec()
			if err != nil {
				panic(err)
			}
		}(addr)
	}

	wg.Wait()
}

type packet struct {
	x, y int
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
