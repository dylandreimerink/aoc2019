package main

import (
	"fmt"
)

const CODE_LENGTH = 6

const CODE_START = 271973
const CODE_STOP = 785961

func main() {

	validPermutations := 0

	for code := CODE_START; code <= CODE_STOP; code++ {
		//Create zero padded code
		codeString := fmt.Sprintf("%06d\n", code)

		decreasing := false
		doubleDigit := false

		for i := 1; i < CODE_LENGTH; i++ {
			prevDigit := codeString[i-1] - 0x30
			thisDigit := codeString[i] - 0x30

			if thisDigit < prevDigit {
				decreasing = true
				break
			}

			if prevDigit == thisDigit {
				doubleDigit = true
			}
		}

		if !decreasing && doubleDigit {
			fmt.Print(codeString)
			validPermutations++
		}
	}

	fmt.Printf("\n\nAmount of valid codes: %d\n", validPermutations)
}
