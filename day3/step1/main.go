package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

//May be tweaked depending on input
const gridSize = 50000

type coordinate struct {
	X, Y int
}

func main() {

	content, err := ioutil.ReadFile("day3/input.txt")
	if err != nil {
		panic(err)
	}

	wires := [][]string{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		wires = append(wires, strings.Split(scanner.Text(), ","))
	}

	wireGrid := make([][]int, gridSize)
	for index := range wireGrid {
		wireGrid[index] = make([]int, gridSize)
	}

	intersections := []coordinate{}

	for index, wire := range wires {
		wireNum := index + 1

		x := gridSize / 2
		y := x

		for _, v := range wire {
			action := v[0]
			steps, err := strconv.Atoi(v[1:])
			if err != nil {
				panic(err)
			}

			compareAndSet := func() {
				defer func() {
					if err := recover(); err != nil {
						panic(fmt.Errorf("Grid size is to small, err: %v", err))
					}
				}()

				if wireGrid[x][y] != wireNum && wireGrid[x][y] != 0 {
					intersections = append(intersections, coordinate{x, y})
				}

				wireGrid[x][y] = wireNum
			}

			compareAndSet()
			for i := 0; i < steps; i++ {

				switch action {
				case 'U':
					y--
					break
				case 'D':
					y++
					break
				case 'L':
					x--
					break
				case 'R':
					x++
					break
				default:
					panic(fmt.Errorf("Unknown action: %v", action))
				}

				compareAndSet()
			}
		}
	}

	smallest := math.MaxInt64

	for _, intersection := range intersections {
		//We intersect at the central port, don't count that one
		if intersection.X == (gridSize/2) && intersection.Y == (gridSize/2) {
			continue
		}

		deltaX := intersection.X - (gridSize / 2)
		if deltaX < 0 {
			deltaX = -deltaX
		}

		deltaY := intersection.Y - (gridSize / 2)
		if deltaY < 0 {
			deltaY = -deltaY
		}

		distance := deltaX + deltaY
		if distance < smallest {
			smallest = distance
		}
	}

	fmt.Printf("Intersection closes to the center point is at a manhattan distance of: %d", smallest)
}
