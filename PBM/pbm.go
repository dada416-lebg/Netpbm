package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
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

	reader := bufio.NewReader(file)

	// Read magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	data := make([][]bool, height)

	for i := range data {
		data[i] = make([]bool, width)
	}

	if magicNumber == "P1" {
		// Read P1 format (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				data[y][x] = field == "1"
			}
		}
	} else if magicNumber == "P4" {
		// Read P4 format (binary)
		expectedBytesPerRow := (width + 7) / 8
		for y := 0; y < height; y++ {
			row := make([]byte, expectedBytesPerRow)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}
			if n < expectedBytesPerRow {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
			}

			for x := 0; x < width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)

				// Convert ASCII to decimal and extract the bit
				decimalValue := int(row[byteIndex])
				bitValue := (decimalValue >> bitIndex) & 1

				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
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
	if pbm.magicNumber == "P1" {
		// Si le format est P1 (ASCII), écrire les données comme des caractères '0' et '1'
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
	} else if pbm.magicNumber == "P4" {
		// Si le format est P4 (binaire), écrire les données sous forme binaire
		for _, row := range pbm.data {
			bytes := make([]byte, (pbm.width+7)/8)
			for x, value := range row {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				if value {
					bytes[byteIndex] |= (1 << uint(bitIndex))
				}
			}
			_, err := writer.Write(bytes)
			if err != nil {
				return fmt.Errorf("error writing binary data: %w", err)
			}
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
