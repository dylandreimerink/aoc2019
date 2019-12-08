package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("day6/input.txt")
	if err != nil {
		panic(err)
	}

	objects := map[string]*Object{}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ")")
		if len(parts) != 2 {
			panic("Error while parsing")
		}

		var objectA *Object
		var found bool
		if objectA, found = objects[parts[0]]; !found {
			objectA = &Object{ID: parts[0], Orbits: []*Object{}}
			objects[objectA.ID] = objectA
		}

		var objectB *Object
		if objectB, found = objects[parts[1]]; !found {
			objectB = &Object{ID: parts[1], Orbits: []*Object{}}
			objects[objectB.ID] = objectB
		}

		objectB.Orbits = append(objectB.Orbits, objectA)
	}

	totalOrbits := 0
	for _, object := range objects {
		totalOrbits += object.getAllOrbits()
	}

	fmt.Printf("Total orbits: %d\n", totalOrbits)
}

type Object struct {
	ID     string
	Orbits []*Object
}

func (o *Object) getAllOrbits() int {
	orbits := 0
	for _, orbit := range o.Orbits {
		orbits++
		orbits += orbit.getAllOrbits()
	}

	return orbits
}
