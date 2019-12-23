package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("day22/input.txt")
	if err != nil {
		panic(err)
	}

	instructions := []Instruction{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "deal into new stack") {
			instructions = append(instructions, Instruction{Operation: dealIntoNewStack})
		}

		if strings.HasPrefix(line, "deal with increment") {
			n, err := strconv.Atoi(strings.TrimSpace(strings.TrimLeft(line, "deal with increment")))
			if err != nil {
				panic(err)
			}

			instructions = append(instructions, Instruction{Operation: dealWithIncrement, N: n})
		}

		if strings.HasPrefix(line, "cut") {
			n, err := strconv.Atoi(strings.TrimSpace(strings.TrimLeft(line, "cut")))
			if err != nil {
				panic(err)
			}

			instructions = append(instructions, Instruction{Operation: cutNCards, N: n})
		}
	}

	reverseInstructions := make([]Instruction, len(instructions))
	for i, inst := range instructions {
		reverseInstructions[len(reverseInstructions)-1-i] = inst
	}

	const decksize = 119315717514047
	// const decksize = 10

	// cards := make([]int, decksize)
	// for i := 0; i < len(cards); i++ {
	// 	cards[i] = i
	// }

	// cards = []int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6}

	pos := 2020

	for i := 0; i < 101741582076661; i++ {
		if i%101741582076 == 0 {
			spew.Dump(i)
		}
		for _, instruction := range reverseInstructions {
			pos = instruction.Operation(instruction.N, pos, decksize)
		}
	}

	spew.Dump(pos)

	// index := -1
	// for i, card := range cards {
	// 	if card == 2019 {
	// 		index = i
	// 		break
	// 	}
	// }

	// fmt.Printf("Card index of card 2019 is %d\n", index)
}

func dealIntoNewStack(N, pos, decksize int) int {

	return (decksize / 2) - (pos - (decksize / 2)) - 1
}

func cutNCards(N, pos, decksize int) int {
	if N > 0 {
		//If pos is in the cut
		if pos < N {
			return pos + decksize - N
		} else {
			//If not in cut we shift by N to the left
			return pos - N
		}
	} else {
		//If pos in the cut
		if pos >= decksize+N {
			return pos - N - decksize
		} else {
			return pos - N
		}
	}

}

func dealWithIncrement(N, pos, decksize int) int {

	if pos%N == 0 {
		return pos / N
	}

	i := 1
	for {
		possible := decksize*i + pos

		if possible%N == 0 {
			return possible / N
		}

		i++
	}
}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

type Instruction struct {
	N         int
	Operation func(int, int, int) int
}
