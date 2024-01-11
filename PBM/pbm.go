package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct { //image BPM
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

	// Lire la première ligne pour obtenir le nombre magique
	if scanner.Scan() {
		magicNumber = scanner.Text()
	} else {
		return nil, fmt.Errorf("erreur lors de la lecture du nombre magique %w", err)
	}

	// Ignorer les commentaires
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			break
		}
	}

	// Lire les dimensions
	if scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("erreur lors de la lecture des dimensions %w", err)
	}

	// Lire les données binaires
	for scanner.Scan() {
		line := scanner.Text()
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

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

func main() {
	filename := "duck.pbm"

	image, err := ReadPBM(filename)
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

}
