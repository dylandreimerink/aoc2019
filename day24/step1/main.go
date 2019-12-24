package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("day24/input.txt")
	if err != nil {
		panic(err)
	}

	grid := [][]bool{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		row := []bool{}
		for _, char := range scanner.Text() {
			row = append(row, char == '#')
		}

		grid = append(grid, row)
	}

	uniqueBios := []int{}

	// renderGrid(grid)
	for {
		cur := calcBioDiversity(grid)
		for _, bio := range uniqueBios {
			if cur == bio {
				fmt.Printf("First dup is: %d\n", cur)
				return
			}
		}

		uniqueBios = append(uniqueBios, cur)
		grid = step(grid)
	}
}

func calcBioDiversity(grid [][]bool) int {
	diversity := 0
	for y := len(grid) - 1; y >= 0; y-- {
		for x := len(grid[y]) - 1; x >= 0; x-- {
			diversity = diversity << 1
			if grid[y][x] {
				diversity += 1
			}
		}
	}
	return diversity
}

func step(grid [][]bool) [][]bool {
	newGrid := make([][]bool, len(grid))
	for y, row := range grid {
		newRow := make([]bool, len(row))
		for x, spot := range row {
			adjacent := 0
			if x-1 >= 0 && grid[y][x-1] {
				adjacent++
			}

			if y-1 >= 0 && grid[y-1][x] {
				adjacent++
			}

			if x+1 < len(row) && grid[y][x+1] {
				adjacent++
			}

			if y+1 < len(grid) && grid[y+1][x] {
				adjacent++
			}

			// fmt.Printf("(%d, %d) %t %d\n", x, y, spot, adjacent)

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

		newGrid[y] = newRow
	}

	return newGrid
}

func renderGrid(grid [][]bool) {
	for _, row := range grid {
		for _, spot := range row {
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
