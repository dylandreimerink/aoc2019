package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
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

	const decksize = 10007

	cards := make([]int, decksize)
	for i := 0; i < len(cards); i++ {
		cards[i] = i
	}

	for _, instruction := range instructions {
		instruction.Operation(instruction.N, cards)
	}

	index := -1
	for i, card := range cards {
		if card == 2019 {
			index = i
			break
		}
	}

	fmt.Printf("Card index of card 2019 is %d\n", index)
}

func dealIntoNewStack(N int, cards []int) {
	for i := len(cards)/2 - 1; i >= 0; i-- {
		opp := len(cards) - 1 - i
		cards[i], cards[opp] = cards[opp], cards[i]
	}
}

func cutNCards(N int, cards []int) {
	cut := make([]int, abs(N))

	if N > 0 {
		//Copy the 'cut' part
		copy(cut, cards[0:N])

		//Move last part of slice to first part
		copy(cards, append(cards[N:], make([]int, N)...))

		//Put cut part at the end
		copy(cards[len(cards)-N:], cut)
	} else {
		//Copy the 'cut' part
		copy(cut, cards[len(cards)+N:])

		//Move last part of slice to first part
		copy(cards, append(make([]int, -N), cards[:len(cards)+N]...))

		// //Put cut part at the end
		copy(cards[:-N], cut)
	}
}

func dealWithIncrement(N int, cards []int) {
	newStack := make([]int, len(cards))

	for i := 0; i < len(cards); i++ {
		newStack[i*N%len(cards)] = cards[i]
	}

	copy(cards, newStack)
}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

type Instruction struct {
	N         int
	Operation func(int, []int)
}
