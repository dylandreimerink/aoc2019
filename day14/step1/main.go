package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("day14/input.txt")
	if err != nil {
		panic(err)
	}

	reactions := map[string]Reaction{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		reaction := Reaction{
			Inputs: map[string]int{},
		}

		parts := strings.Split(scanner.Text(), "=>")
		inputs := strings.Split(strings.TrimSpace(parts[0]), ",")
		for _, input := range inputs {

			inputParts := strings.Split(strings.TrimSpace(input), " ")

			var err error
			reaction.Inputs[inputParts[1]], err = strconv.Atoi(inputParts[0])
			if err != nil {
				panic("Invalid input")
			}
		}

		outputParts := strings.Split(strings.TrimSpace(parts[1]), " ")

		var err error
		reaction.Output = outputParts[1]
		reaction.OutputAmount, err = strconv.Atoi(outputParts[0])
		if err != nil {
			panic("Invalid input")
		}

		reactions[reaction.Output] = reaction
	}

	required := map[string]int{}
	producing := map[string]int{}

	required["FUEL"] = 1

	var produce func(name string)
	produce = func(name string) {
		amount := required[name] - producing[name]
		if amount < 0 {
			amount = 0
		}

		amount = int(math.Ceil(float64(amount) / float64(reactions[name].OutputAmount)))

		producing[name] += amount * reactions[name].OutputAmount
		for inName, inCount := range reactions[name].Inputs {
			required[inName] += amount * inCount
		}

		for inName := range reactions[name].Inputs {
			if inName != "ORE" {
				produce(inName)
			}
		}
	}

	produce("FUEL")
	spew.Dump(required["ORE"])
}

type Reaction struct {
	Inputs       map[string]int
	Output       string
	OutputAmount int
}
