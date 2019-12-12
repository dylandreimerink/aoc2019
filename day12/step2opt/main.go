package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/davecgh/go-spew/spew"
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

	repX := 0
	uniqueXs := [][][2]int{}
	repY := 0
	uniqueYs := [][][2]int{}
	repZ := 0
	uniqueZs := [][][2]int{}

	for i := 1; i <= 1000000; i++ {
		if i%100 == 0 {
			fmt.Printf("Simulating step: %d\n", i)
		}

		if repX == 0 {
			uniqueXs = append(uniqueXs, [][2]int{
				[2]int{moons[0].x, moons[0].vx},
				[2]int{moons[1].x, moons[1].vx},
				[2]int{moons[2].x, moons[2].vx},
				[2]int{moons[3].x, moons[3].vx},
			})
		}
		if repY == 0 {
			uniqueYs = append(uniqueYs, [][2]int{
				[2]int{moons[0].y, moons[0].vy},
				[2]int{moons[1].y, moons[1].vy},
				[2]int{moons[2].y, moons[2].vy},
				[2]int{moons[3].y, moons[3].vy},
			})
		}
		if repZ == 0 {
			uniqueZs = append(uniqueZs, [][2]int{
				[2]int{moons[0].z, moons[0].vz},
				[2]int{moons[1].z, moons[1].vz},
				[2]int{moons[2].z, moons[2].vz},
				[2]int{moons[3].z, moons[3].vz},
			})
		}

		if repX != 0 && repY != 0 && repZ != 0 {
			break
		}

		simulateStep(moons)

		if repX == 0 {
			for _, test := range uniqueXs {
				equal := true
				for index, moon := range moons {
					if !(test[index][0] == moon.x && test[index][1] == moon.vx) {
						equal = false
					}
				}
				if equal {
					spew.Dump("Found x")
					repX = i
				}
			}
		}

		if repY == 0 {
			for _, test := range uniqueYs {
				equal := true
				for index, moon := range moons {
					if !(test[index][0] == moon.y && test[index][1] == moon.vy) {
						equal = false
					}
				}
				if equal {
					spew.Dump("Found y")
					repY = i
				}
			}
		}

		if repZ == 0 {
			for _, test := range uniqueZs {
				equal := true
				for index, moon := range moons {
					if !(test[index][0] == moon.z && test[index][1] == moon.vz) {
						equal = false
					}
				}
				if equal {
					spew.Dump("Found z")
					repZ = i
				}
			}
		}
	}

	spew.Dump(repX, repY, repZ)
	spew.Dump(LCM(repX, repY, repZ))
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

// greatest common divisor (GCD) via Euclidean algorithm
// Source: https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
// Source: https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
