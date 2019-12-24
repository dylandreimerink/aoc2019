package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

const levels = 201
const size = 5

func main() {
	content, err := ioutil.ReadFile("day24/input.txt")
	if err != nil {
		panic(err)
	}

	grid := make([][][]bool, levels)
	for i := range grid {
		grid[i] = make([][]bool, size)
		for ii := range grid[i] {
			grid[i][ii] = make([]bool, size)
		}
	}

	layer := [][]bool{}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		row := []bool{}
		for _, char := range scanner.Text() {
			row = append(row, char == '#')
		}

		layer = append(layer, row)
	}

	//Put the parsed layer in the middle
	grid[levels/2] = layer

	for i := 0; i < 200; i++ {
		grid = step(grid)
	}

	// renderGrid(grid)

	count := 0

	for _, layer := range grid {
		for _, row := range layer {
			for _, spot := range row {
				if spot {
					count++
				}
			}
		}
	}

	fmt.Println(count)
}

func step(grid [][][]bool) [][][]bool {

	newGrid := make([][][]bool, len(grid))

	for z, layer := range grid {

		newLayer := make([][]bool, len(layer))
		for y, row := range layer {

			newRow := make([]bool, len(row))
			for x, spot := range row {
				if x == 2 && y == 2 {
					continue
				}

				adjacent := 0

				//If one spot left on the same level has a bug
				if x-1 >= 0 && grid[z][y][x-1] {
					//Unless that spot is the middle
					if !(y == 2 && x-1 == 2) {
						adjacent++
					}
				}

				//If one spot up on the same level has a bug
				if y-1 >= 0 && grid[z][y-1][x] {
					//Unless that spot is the middle
					if !(y-1 == 2 && x == 2) {
						adjacent++
					}
				}

				//If one spot right on the same level has a bug
				if x+1 < len(row) && grid[z][y][x+1] {
					//Unless that spot is the middle
					if !(y == 2 && x+1 == 2) {
						adjacent++
					}
				}

				//If one spot down on the same level has a bug
				if y+1 < len(layer) && grid[z][y+1][x] {
					//Unless that spot is the middle
					if !(y+1 == 2 && x == 2) {
						adjacent++
					}
				}

				//If left size and there is a level above check spot 12
				if x == 0 && z-1 >= 0 && grid[z-1][2][1] {
					adjacent++
				}

				//If right size and there is a level above check spot 14
				if x == size-1 && z-1 >= 0 && grid[z-1][2][3] {
					adjacent++
				}

				//If top size and there is a level above check spot 8
				if y == 0 && z-1 >= 0 && grid[z-1][1][2] {
					adjacent++
				}

				//If bottom size and there is a level above check spot 18
				if y == size-1 && z-1 >= 0 && grid[z-1][3][2] {
					adjacent++
				}

				//If spot 12 and sub level exists
				if y == 2 && x == 1 && z+1 < len(grid) {
					//for every spot on the left side
					for i := 0; i < size; i++ {
						if grid[z+1][i][0] {
							adjacent++
						}
					}
				}

				//If spot 14 and sub level exists
				if y == 2 && x == 3 && z+1 < len(grid) {
					//for every spot on the right side
					for i := 0; i < size; i++ {
						if grid[z+1][i][size-1] {
							adjacent++
						}
					}
				}

				//If spot 8 and sub level exists
				if y == 1 && x == 2 && z+1 < len(grid) {
					//for every spot on the top side
					for i := 0; i < size; i++ {
						if grid[z+1][0][i] {
							adjacent++
						}
					}
				}

				//If spot 18 and sub level exists
				if y == 3 && x == 2 && z+1 < len(grid) {
					//for every spot on the bottom side
					for i := 0; i < size; i++ {
						if grid[z+1][size-1][i] {
							adjacent++
						}
					}
				}

				if spot {
					if adjacent == 1 {
						newRow[x] = true
					}
				} else {
					if adjacent == 1 || adjacent == 2 {
						newRow[x] = true
					}
				}
			}

			newLayer[y] = newRow
		}
		newGrid[z] = newLayer
	}

	return newGrid
}

func renderGrid(grid [][][]bool) {
	for z, layer := range grid {
		fmt.Printf("Layer %d\n", -(levels/2)+z)
		for y, row := range layer {
			for x, spot := range row {
				if x == 2 && y == 2 {
					fmt.Print("?")
					continue
				}
				if spot {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Print("\n")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
