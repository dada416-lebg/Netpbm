package main

import (
	"bufio"
	"fmt"
	"os"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var magicNumber string
	var width, height int
	var data [][]bool

	// Read the magic number
	if scanner.Scan() {
		magicNumber = scanner.Text()
		fmt.Println("Magic Number:", magicNumber)
	} else {
		return nil, fmt.Errorf("error reading magic number: %w", err)
	}

	// Read dimensions
	if scanner.Scan() {
		dimensionLine := scanner.Text()
		fmt.Println("Dimensions Line:", dimensionLine)

		_, err := fmt.Sscanf(dimensionLine, "%d %d", &width, &height)
		if err != nil {
			return nil, fmt.Errorf("error reading dimensions: %w", err)
		}

		fmt.Println("Dimensions:", width, height)
	} else {
		return nil, fmt.Errorf("error reading dimensions: no data found")
	}

	// Read binary data based on the magic number (only P1 supported)
	for scanner.Scan() {
		line := scanner.Text()

		if magicNumber == "P1" { // PBM Plain (ASCII)
			var row []bool
			for _, char := range line {
				if char == '1' {
					row = append(row, true)
				} else if char == '0' {
					row = append(row, false)
				}
			}
			data = append(data, row)
		}
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func main() {
	image, err := ReadPBM("duck.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier PBM:", err)
		return
	}

	for _, row := range image.data {
		for _, value := range row {
			if value {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println() // Nouvelle ligne pour chaque ligne de l'image
	}

	fmt.Println(image.Size())
}
