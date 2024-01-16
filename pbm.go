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

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	// Vérifier si les indices x et y sont dans les limites de l'image
	if x >= 0 && x < pbm.width && y >= 0 && y < pbm.height {
		return pbm.data[y][x] // Accéder à la valeur du pixel
	}
	// Si les indices sont hors limites, renvoyer une valeur par défaut (par exemple, false)
	return false
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	// Vérifier si les indices x et y sont dans les limites de l'image
	if x >= 0 && x < pbm.width && y >= 0 && y < pbm.height {
		pbm.data[y][x] = value // Mettre à jour la valeur du pixel
	}
	// Si les indices sont hors limites, ne rien faire (ignorer la mise à jour)
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Écrire le magic number et les dimensions dans le fichier
	_, err = fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)
	if err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	// Écrire les données binaires de l'image dans le fichier
	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				_, err := writer.WriteString("1 ")
				if err != nil {
					return fmt.Errorf("error writing data: %w", err)
				}
			} else {
				_, err := writer.WriteString("0 ")
				if err != nil {
					return fmt.Errorf("error writing data: %w", err)
				}
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return fmt.Errorf("error writing data: %w", err)
		}
	}

	// Vider le tampon pour s'assurer que toutes les données sont écrites dans le fichier
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			// Inverser la valeur du pixel
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	// Parcourir chaque ligne de l'image
	for y := 0; y < pbm.height; y++ {
		// Initialiser une nouvelle ligne pour stocker les pixels inversés
		flippedRow := make([]bool, pbm.width)

		// Parcourir chaque colonne de l'image
		for x := 0; x < pbm.width; x++ {
			// Inverser la position du pixel horizontalement
			flippedRow[x] = pbm.data[y][pbm.width-1-x]
		}

		// Mettre à jour la ligne d'origine avec les pixels inversés
		pbm.data[y] = flippedRow
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	// Initialiser une nouvelle matrice pour stocker les lignes inversées
	flippedData := make([][]bool, pbm.height)

	// Parcourir chaque ligne de l'image
	for y := 0; y < pbm.height; y++ {
		// Inverser la position de la ligne verticalement
		flippedData[pbm.height-1-y] = pbm.data[y]
	}

	// Mettre à jour les données d'origine avec les lignes inversées
	pbm.data = flippedData
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

/*
func main() {
	image, err := ReadPBM("duck.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier PBM:", err)
		return
	}

	// Afficher l'image avant l'inversion
	fmt.Println("Image avant l'inversion :")
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

	// Inverser les couleurs de l'image
	image.Invert()

	// Afficher l'image après l'inversion
	fmt.Println("Image après l'inversion :")
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

	// Retourner l'image horizontalement
	image.Flip()

	// Enregistrement de l'image inversée et retournée dans un nouveau fichier
	err = image.Save("duck_inverted_and_flipped.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image inversée et retournée:", err)
		return
	}

	fmt.Println("L'image inversée et retournée a été enregistrée avec succès.")
	// ...
	// Retourner l'image verticalement
	image.Flop()

	// Enregistrement de l'image inversée, retournée et floppée dans un nouveau fichier
	err = image.Save("duck_inverted_flipped_and_flopped.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image inversée, retournée et floppée:", err)
		return
	}

	fmt.Println("L'image inversée, retournée et floppée a été enregistrée avec succès.")

	// Modifier le nombre magique de l'image
	image.SetMagicNumber("P4")

	// Enregistrement de l'image modifiée dans un nouveau fichier
	err = image.Save("duck_modified.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image modifiée:", err)
		return
	}

	fmt.Println("L'image modifiée a été enregistrée avec succès.")

}
*/
