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

	components := map[string]int{}
	getComponentes(reactions["FUEL"], 1, reactions, components)

	spew.Dump(components)

	for name, req := range components {
		div := req / reactions[name].OutputAmount
		mod := req % reactions[name].OutputAmount

		if mod != 0 {
			components[name] = (div + 1) * reactions[name].OutputAmount
		}
	}

	oreReq := 0
	for name, req := range components {
		for inName, inAmount := range reactions[name].Inputs {
			if inName == "ORE" {
				oreReq += (req / reactions[name].OutputAmount) * inAmount
			}
		}
	}

	spew.Dump(components)
	spew.Dump(oreReq)
}

func getComponentes(cur Reaction, amount int, reactions map[string]Reaction, components map[string]int) {

	for name, amount := range cur.Inputs {
		if name == "ORE" {
			continue
		}
		components[name] += amount
		for i := 0; i < amount; i++ {
			getComponentes(reactions[name], amount, reactions, components)
		}
	}

}

// func getOreRequirement(cur Reaction, reactions []Reaction, inventory map[string]int) int {
// 	ore := 0

// 	spew.Dump(cur.Output)
// 	spew.Dump(inventory)

// 	if oreAmount, found := cur.Inputs["ORE"]; found {
// 		if len(cur.Inputs) > 1 {
// 			panic("Oeps")
// 		}

// 		inventory[cur.Output] = cur.OutputAmount

// 		return oreAmount
// 	}

// 	for name, amount := range cur.Inputs {
// 		for _, reaction := range reactions {
// 			if reaction.Output == name {
// 				for inventory[name] < amount {
// 					ore += getOreRequirement(reaction, reactions, inventory)
// 				}

// 				inventory[name] = inventory[name] - amount
// 			}
// 		}
// 	}

// 	inventory[cur.Output] = inventory[cur.Output] + cur.OutputAmount

// 	return ore
// }

type Reaction struct {
	Inputs       map[string]int
	Output       string
	OutputAmount int
}
