package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/RyanCarrier/dijkstra"
)

func main() {
	content, err := ioutil.ReadFile("day20/input.txt")
	if err != nil {
		panic(err)
	}

	maze := [][]rune{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		row := []rune{}
		for _, char := range scanner.Text() {
			row = append(row, char)
		}
		maze = append(maze, row)
	}

	portals := findPortals(maze)
	passages := findPassages(maze, portals)

	uniquePortals := []string{}
	for _, portal := range portals {
		unique := true
		for _, uniquePortal := range uniquePortals {
			if uniquePortal == portal.label {
				unique = false
			}
		}
		if unique {
			uniquePortals = append(uniquePortals, portal.label)
		}
	}

	startID := -1
	endID := -1

	graph := dijkstra.NewGraph()
	for id, label := range uniquePortals {
		if label == "AA" {
			startID = id
		}
		if label == "ZZ" {
			endID = id
		}
		graph.AddVertex(id)
	}

	for _, passage := range passages {
		aid := -1
		for id, portal := range uniquePortals {
			if portal == passage.portalA.label {
				aid = id
				break
			}
		}

		bid := -1
		for id, portal := range uniquePortals {
			if portal == passage.portalB.label {
				bid = id
				break
			}
		}

		//Add arc from a to b
		err := graph.AddArc(aid, bid, int64(passage.length))
		if err != nil {
			panic(err)
		}
	}

	path, err := graph.Shortest(startID, endID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shortest path is %d long\n", path.Distance-1)

	// renderMaze(maze)
}

type coord struct {
	x, y int
}

type portal struct {
	label        string
	firstLetter  coord
	secondLetter coord
	passage      coord
}

type passage struct {
	portalA portal
	portalB portal
	length  int
}

func findPassages(maze [][]rune, portals []portal) []passage {
	passages := []passage{}

	for _, portalA := range portals {

		//Make a copy so we can change it
		mazeCopy := [][]rune{}
		for _, row := range maze {
			rowCopy := make([]rune, len(row))
			copy(rowCopy, row)
			mazeCopy = append(mazeCopy, rowCopy)
		}

		mazeCopy[portalA.passage.y][portalA.passage.x] = 'x'

		steps := 1

		for {
			progress := false
			for y, row := range mazeCopy {
				for x, tile := range row {
					if tile == 'x' {
						progress = true
						mazeCopy[y][x] = '*'
					}
				}
			}

			steps++

			if !progress {
				break
			}

			checkPortal := func(x, y int) {
				for _, portalB := range portals {
					if x == portalB.passage.x && y == portalB.passage.y {
						passages = append(passages, passage{portalA: portalA, portalB: portalB, length: steps})
					}
				}
			}

			for y, row := range mazeCopy {
				for x, tile := range row {
					if tile == '*' {

						if mazeCopy[y][x+1] == '.' {
							mazeCopy[y][x+1] = 'x'
							checkPortal(x+1, y)
						}

						if mazeCopy[y+1][x] == '.' {
							mazeCopy[y+1][x] = 'x'
							checkPortal(x, y+1)
						}

						if mazeCopy[y][x-1] == '.' {
							mazeCopy[y][x-1] = 'x'
							checkPortal(x-1, y)
						}

						if mazeCopy[y-1][x] == '.' {
							mazeCopy[y-1][x] = 'x'
							checkPortal(x, y-1)
						}

					}
				}
			}

			// time.Sleep(200 * time.Millisecond)
			// renderMaze(mazeCopy)
		}
	}

	return passages
}

func findPortals(maze [][]rune) []portal {
	portals := []portal{}

	for y, row := range maze {
		for x, tile := range row {
			//If label
			if tile >= 'A' && tile <= 'Z' {

				newPortal := portal{
					label:       string(tile),
					firstLetter: coord{x: x, y: y},
				}

				//Below
				if y+1 < len(maze) && maze[y+1][x] >= 'A' && maze[y+1][x] <= 'Z' {
					newPortal.secondLetter = coord{x: x, y: y + 1}
					newPortal.label = newPortal.label + string(maze[y+1][x])

					//If below the second letter is a passage
					if y+2 < len(maze) && maze[y+2][x] == '.' {
						newPortal.passage = coord{x: x, y: y + 2}
						portals = append(portals, newPortal)
						continue
					}

					//If above the first letter is a passage
					if y-1 >= 0 && maze[y-1][x] == '.' {
						newPortal.passage = coord{x: x, y: y - 1}
						portals = append(portals, newPortal)
						continue
					}
				}

				//Right
				if x+1 < len(maze[y]) && maze[y][x+1] >= 'A' && maze[y][x+1] <= 'Z' {
					newPortal.secondLetter = coord{x: x + 1, y: y}
					newPortal.label = newPortal.label + string(maze[y][x+1])

					//If right of the second letter is a passage
					if x+2 < len(maze[y]) && maze[y][x+2] == '.' {
						newPortal.passage = coord{x: x + 2, y: y}
						portals = append(portals, newPortal)
						continue
					}

					//If left of the first letter is a passage
					if x-1 >= 0 && maze[y][x-1] == '.' {
						newPortal.passage = coord{x: x - 1, y: y}
						portals = append(portals, newPortal)
						continue
					}
				}
			}
		}
	}

	return portals
}

func renderMaze(maze [][]rune) {
	for _, row := range maze {
		for _, tile := range row {
			fmt.Print(string(tile))
		}
		fmt.Print("\n")
	}
}
