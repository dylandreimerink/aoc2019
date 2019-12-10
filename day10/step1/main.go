package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("day10/input.txt")
	if err != nil {
		panic(err)
	}

	var astroidMap [][]int

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]int, len(line))
		astroidMap = append(astroidMap, row)

		for i, char := range line {
			row[i] = 0
			if char == '#' {
				row[i] = 1
			}
		}
	}

	spew.Dump(bestMonitorSpot(astroidMap))
}

func bestMonitorSpot(astroidMap [][]int) (int, int, int) {

	mostAstroids := 0
	bestX := 0
	bestY := 0

	for iy := 0; iy < len(astroidMap); iy++ {
		for ix := 0; ix < len(astroidMap[iy]); ix++ {
			if astroidMap[iy][ix] == 1 {
				count := countSeenAstroids(astroidMap, ix, iy)
				if count > mostAstroids {
					mostAstroids = count
					bestX = ix
					bestY = iy
				}
			}
		}
	}

	return bestX, bestY, mostAstroids
}

func countSeenAstroids(astroidMap [][]int, x, y int) int {
	var seenAngles []float64

	astroids := 0

	for iy := 0; iy < len(astroidMap); iy++ {
		for ix := 0; ix < len(astroidMap[iy]); ix++ {
			if iy == y && ix == x {
				continue
			}

			//An astroid
			if astroidMap[iy][ix] == 1 {
				dx := float64(x - ix)
				dy := float64(y - iy)

				theta := math.Atan2(dy, dx)

				found := false
				for _, angle := range seenAngles {
					if angle == theta {
						found = true
						break
					}
				}

				if !found {
					seenAngles = append(seenAngles, theta)
					astroids++
				}
			}
		}
	}

	return astroids
}
