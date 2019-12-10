package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
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

	x, y, _ := bestMonitorSpot(astroidMap)
	nx, ny := nthAstroidStanding(astroidMap, x, y, 200)
	fmt.Printf("200th @ (%d, %d), Awnser: %d\n", nx, ny, nx*100+ny)
}

func nthAstroidStanding(astroidMap [][]int, x, y, n int) (int, int) {

	astroidMap[y][x] = 2

	astroidsPerAngle := countSeenAstroids(astroidMap, x, y)
	angles := make([]float64, 0, len(astroidsPerAngle))

	for angle := range astroidsPerAngle {
		angles = append(angles, angle)
	}

	sort.Float64s(angles)

	zeroIndex := 0
	for index, angle := range angles {
		if angle >= 0 {
			zeroIndex = index
			break
		}
	}

	astroids := 0
	i := zeroIndex
	for {
		angle := angles[i%len(angles)]

		if len(astroidsPerAngle[angle]) > 0 {
			astroids++

			closest := 0
			distance := float64(99999999999)
			for index, astroid := range astroidsPerAngle[angle] {
				if astroid.distance < distance {
					distance = astroid.distance
					closest = index
				}
			}

			if astroids >= n {
				winner := astroidsPerAngle[angle][closest]
				return winner.x, winner.y
			}

			astroidMap[astroidsPerAngle[angle][closest].y][astroidsPerAngle[angle][closest].x] = 3

			astroidsPerAngle[angle] = append(astroidsPerAngle[angle][:closest], astroidsPerAngle[angle][closest+1:]...)
		}

		i++
	}
}

func bestMonitorSpot(astroidMap [][]int) (int, int, int) {

	mostAstroids := 0
	bestX := 0
	bestY := 0

	for iy := 0; iy < len(astroidMap); iy++ {
		for ix := 0; ix < len(astroidMap[iy]); ix++ {
			if astroidMap[iy][ix] == 1 {
				astroidsPerAngle := countSeenAstroids(astroidMap, ix, iy)
				if len(astroidsPerAngle) > mostAstroids {
					mostAstroids = len(astroidsPerAngle)
					bestX = ix
					bestY = iy
				}
			}
		}
	}

	return bestX, bestY, mostAstroids
}

type astroid struct {
	x, y     int
	distance float64
}

func countSeenAstroids(astroidMap [][]int, x, y int) map[float64][]astroid {
	astroidsPerAngle := map[float64][]astroid{}

	for iy := 0; iy < len(astroidMap); iy++ {
		for ix := 0; ix < len(astroidMap[iy]); ix++ {
			if iy == y && ix == x {
				continue
			}

			//An astroid
			if astroidMap[iy][ix] == 1 {
				dx := float64(x - ix)
				dy := float64(y - iy)

				theta := math.Atan2(dy, dx) - (math.Pi / 2)

				_, found := astroidsPerAngle[theta]
				if !found {
					astroidsPerAngle[theta] = []astroid{}
				}

				distance := math.Sqrt(math.Pow(math.Abs(dx), 2) + math.Pow(math.Abs(dy), 2))

				astroidsPerAngle[theta] = append(astroidsPerAngle[theta], astroid{x: ix, y: iy, distance: distance})
			}
		}
	}

	return astroidsPerAngle
}

func renderAstroidMap(astroidMap [][]int) {
	for iy := 0; iy < len(astroidMap); iy++ {
		for ix := 0; ix < len(astroidMap[iy]); ix++ {
			switch astroidMap[iy][ix] {
			case 0:
				fmt.Print(".")
			case 1:
				fmt.Print("#")
			case 2:
				fmt.Print("%")
			case 3:
				fmt.Print("*")
			}
		}
		fmt.Print("\n")
	}
}
