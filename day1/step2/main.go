package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("day1/input.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	totalFuel := 0
	for scanner.Scan() {
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}

		totalFuel += getFuelForMass(mass)
	}

	spew.Dump(totalFuel)
}

func getFuelForMass(mass int) int {
	totalFuel := 0
	fuel := int(math.Floor(float64(mass/3))) - 2
	totalFuel += fuel

	for fuel > 0 {
		fuel = int(math.Floor(float64(fuel/3))) - 2
		if fuel > 0 {
			totalFuel += fuel
		}
	}

	return totalFuel
}
