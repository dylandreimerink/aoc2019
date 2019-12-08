package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("day8/input.txt")
	if err != nil {
		panic(err)
	}

	layers := decodeImage(content)

	renderImage(layers)
}

const width = 25
const height = 6

func renderImage(image [][][]int) {
	renderedImage := [][]int{}
	for i := 0; i < height; i++ {
		renderedImage = append(renderedImage, make([]int, width))
	}

	for i := len(image) - 1; i >= 0; i-- {
		for y, rows := range image[i] {
			for x, pixel := range rows {
				if pixel != 2 {
					renderedImage[y][x] = pixel
				}
			}
		}
	}

	for _, rows := range renderedImage {
		for _, pixel := range rows {
			if pixel == 0 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func imageChecksum(image [][][]int) int {
	leastZeros := 99999999
	leastZeroOnes := 0
	leastZeroTwos := 0

	for _, layer := range image {
		layerZeroCount := 0
		ones := 0
		twos := 0
		for _, row := range layer {
			for _, pixel := range row {
				if pixel == 0 {
					layerZeroCount++
				}

				if pixel == 1 {
					ones++
				}

				if pixel == 2 {
					twos++
				}
			}
		}

		if layerZeroCount < leastZeros {
			leastZeros = layerZeroCount
			leastZeroOnes = ones
			leastZeroTwos = twos
		}
	}

	return leastZeroOnes * leastZeroTwos
}

func decodeImage(content []byte) [][][]int {
	layers := [][][]int{}
	curLayer := 0
	curX := 0
	curY := 0
	for _, char := range content {
		pixelValue := int(char - 0x30) //ASCII decoding of numbers

		//if layer doesn't exist make it
		if len(layers) <= curLayer {
			layers = append(layers, [][]int{})
			for i := 0; i < height; i++ {
				layers[curLayer] = append(layers[curLayer], make([]int, width))
			}
		}

		layers[curLayer][curY][curX] = pixelValue

		curX++

		if curX == width {
			curX = 0
			curY++
		}

		if curY == height {
			curY = 0
			curLayer++
		}
	}

	return layers
}
