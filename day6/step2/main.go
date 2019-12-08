package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
			objectA = &Object{
				ID:         parts[0],
				Orbits:     []*Object{},
				OrbittedBy: []*Object{},
			}

			objects[objectA.ID] = objectA
		}

		var objectB *Object
		if objectB, found = objects[parts[1]]; !found {
			objectB = &Object{
				ID:         parts[1],
				Orbits:     []*Object{},
				OrbittedBy: []*Object{},
			}

			objects[objectB.ID] = objectB
		}

		objectB.Orbits = append(objectB.Orbits, objectA)
		objectA.OrbittedBy = append(objectA.OrbittedBy, objectB)
	}

	const from = "YOU"
	const to = "SAN"

	next := "COM"
	for {
		_, toNext := objects[next].IsOrbittedBy(to)
		_, fromNext := objects[next].IsOrbittedBy(from)
		if toNext != fromNext {
			toLength := 0
			for {
				toLength++
				if _, toNext = objects[toNext].IsOrbittedBy(to); toNext == to {
					break
				}
			}

			fromLength := 0
			for {
				fromLength++
				if _, fromNext = objects[fromNext].IsOrbittedBy(from); fromNext == from {
					break
				}
			}

			spew.Dump(toLength + fromLength)

			break
		}

		next = toNext
	}
}

type Object struct {
	ID         string
	Orbits     []*Object
	OrbittedBy []*Object
}

func (o *Object) getAllOrbits() int {
	orbits := 0
	for _, orbit := range o.Orbits {
		orbits++
		orbits += orbit.getAllOrbits()
	}

	return orbits
}

func (o *Object) IsOrbittedBy(id string) (bool, string) {
	for _, orbit := range o.OrbittedBy {
		if orbit.ID == id {
			return true, orbit.ID
		}

		if ok, _ := orbit.IsOrbittedBy(id); ok {
			return true, orbit.ID
		}
	}

	return false, ""
}
