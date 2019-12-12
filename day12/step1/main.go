package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

var coordinateRegex = regexp.MustCompile(`<x=(-?\d+),\s*y=(-?\d+),\s*z=(-?\d+)>`)

func main() {
	content, err := ioutil.ReadFile("day12/input.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	moons := []*moon{}

	for scanner.Scan() {
		match := coordinateRegex.FindStringSubmatch(scanner.Text())
		if match == nil {
			panic("Error parsing input")
		}

		x, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(match[2])
		if err != nil {
			panic(err)
		}
		z, err := strconv.Atoi(match[3])
		if err != nil {
			panic(err)
		}

		moons = append(moons, &moon{
			x: x,
			y: y,
			z: z,
		})
	}

	//Print moons
	fmt.Println("After 0 steps:")
	for _, moon := range moons {
		fmt.Println(moon)
	}

	for i := 1; i <= 1000; i++ {
		simulateStep(moons)

		if i%100 == 0 {
			//Print moons
			fmt.Printf("\nAfter %d steps:\n", i)
			for _, moon := range moons {
				fmt.Println(moon)
			}
		}
	}

	totalEnergy := 0

	for _, moon := range moons {
		totalEnergy += moon.kineticEnergy() * moon.potentialEnergy()
	}

	fmt.Printf("\nTotal energy is: %d\n", totalEnergy)
}

func simulateStep(moons []*moon) {
	applyGravity(moons)
	applyVelocity(moons)
}

func applyGravity(moons []*moon) {
	//Apply gravity
	for aIndex, moonA := range moons {
		for bIndex, moonB := range moons {
			if aIndex == bIndex {
				continue
			}

			if moonA.x > moonB.x {
				moonA.vx--
			}
			if moonA.x < moonB.x {
				moonA.vx++
			}
			if moonA.y > moonB.y {
				moonA.vy--
			}
			if moonA.y < moonB.y {
				moonA.vy++
			}
			if moonA.z > moonB.z {
				moonA.vz--
			}
			if moonA.z < moonB.z {
				moonA.vz++
			}
		}
	}
}

func applyVelocity(moons []*moon) {
	//Modify position using velocity
	for _, moon := range moons {
		moon.x += moon.vx
		moon.y += moon.vy
		moon.z += moon.vz
	}
}

type moon struct {
	//Position
	x, y, z int

	//Velocity
	vx, vy, vz int
}

func (m moon) String() string {
	return fmt.Sprintf("pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>", m.x, m.y, m.z, m.vx, m.vy, m.vz)
}

func (m *moon) potentialEnergy() int {
	return abs(m.x) + abs(m.y) + abs(m.z)
}

func (m *moon) kineticEnergy() int {
	return abs(m.vx) + abs(m.vy) + abs(m.vz)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}
