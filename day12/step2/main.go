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

	moonsInitialState := []*moon{}

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

		moonsInitialState = append(moonsInitialState, &moon{
			x: x,
			y: y,
			z: z,
		})
	}

	moons := copyMoons(moonsInitialState)

	//Print moons
	fmt.Println("After 0 steps:")
	for _, moon := range moons {
		fmt.Println(moon)
	}

	i := 0
	done := false
	for !done {
		if i%100 == 0 {
			fmt.Printf("Simulating step: %d\n", i)
		}

		simulateStep(moons)

		cmpMoons := copyMoons(moonsInitialState)
		for ii := 1; ii < i; ii++ {
			simulateStep(cmpMoons)

			equal := true
			for index := range cmpMoons {
				if !cmpMoons[index].Equals(moons[index]) {
					equal = false
				}
			}

			if equal {
				fmt.Printf("\nMoons repeated after %d \n", i)
				for _, moon := range moons {
					fmt.Println(moon)
					done = true
					break
				}
			}
		}

		i++
	}
}

func copyMoons(original []*moon) []*moon {
	moons := []*moon{}

	for _, moon := range original {
		//deref moon and copy value
		newMoon := *moon
		moons = append(moons, &newMoon)
	}

	return moons
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

func (m *moon) Equals(cmp *moon) bool {
	return m.x == cmp.x &&
		m.y == cmp.y &&
		m.z == cmp.z &&
		m.vx == cmp.vx &&
		m.vy == cmp.vy &&
		m.vz == cmp.vz
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
