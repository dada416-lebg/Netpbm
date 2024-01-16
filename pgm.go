package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data        [][]uint8
	width       int
	height      int
	magicNumber string
	max         int
}

// PBM représente une image PBM (Portable Bitmap).
type PBM struct {
	width  int
	height int
	data   [][]bool
}

// ReadPGM lit une image PGM à partir d'un fichier et renvoie une structure représentant l'image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pgm PGM

	// Lecture de l'en-tête PGM
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignorer les commentaires
			continue
		}

		if pgm.magicNumber == "" {
			pgm.magicNumber = line
		} else if pgm.width == 0 {
			// Lire les dimensions de l'images
			dimensions := strings.Fields(line)
			if len(dimensions) != 2 {
				return nil, fmt.Errorf("format d'en-tête incorrect")
			}
			pgm.width, _ = strconv.Atoi(dimensions[0])
			pgm.height, _ = strconv.Atoi(dimensions[1])
		} else if pgm.max == 0 {
			// Lire la valeur maximale de l'image
			pgm.max, _ = strconv.Atoi(line)
			break
		}
	}

	if pgm.magicNumber != "P2" && pgm.magicNumber != "P5" {
		return nil, fmt.Errorf("format pgm non pris en charge: %s", pgm.magicNumber)
	}

	// Lecture des données de l'image
	// Lecture des données de l'image
	pgm.data = make([][]uint8, pgm.height)
	for i := 0; i < pgm.height; i++ {
		scanner.Scan()
		line := scanner.Text()
		values := strings.Fields(line)
		if len(values) != pgm.width {
			return nil, fmt.Errorf("nombre incorrect de valeurs dans la ligne")
		}

		pgm.data[i] = make([]uint8, pgm.width)
		for j := 0; j < pgm.width; j++ {
			value, err := strconv.Atoi(values[j])
			if err != nil {
				return nil, err
			}
			pgm.data[i][j] = uint8(value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &pgm, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At retourne la valeur du pixel à la position (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	// Vérifier les limites
	if x < 0 || x >= pgm.width || y < 0 || y >= pgm.height {
		fmt.Println("Coordonnées hors limites")
		return 0
	}

	return pgm.data[y][x]
}

// Set définit la valeur du pixel à la position (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	// Vérifier les limites
	if x < 0 || x >= pgm.width || y < 0 || y >= pgm.height {
		fmt.Println("Coordonnées hors limites")
		return
	}

	pgm.data[y][x] = value
}

// Save enregistre l'image PGM dans un fichier et renvoie une erreur en cas de problème.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire l'en-tête PGM
	header := fmt.Sprintf("%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	// Écrire les données de l'image
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			_, err := file.WriteString(fmt.Sprintf("%d ", pgm.data[i][j]))
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Invert inverse les couleurs de l'image PGM.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip retourne l'image PGM horizontalement.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			// Échanger les pixels symétriques par rapport à l'axe vertical
			pgm.data[i][j], pgm.data[i][pgm.width-1-j] = pgm.data[i][pgm.width-1-j], pgm.data[i][j]
		}
	}
}

// Flop renverse l'image PGM verticalement.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		// Échanger les lignes symétriques par rapport à l'axe horizontal
		pgm.data[i], pgm.data[pgm.height-1-i] = pgm.data[pgm.height-1-i], pgm.data[i]
	}
}

// SetMagicNumber définit le nombre magique de l'image PGM.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PGM.
func (pgm *PGM) SetMaxValue(maxValue int) {
	pgm.max = maxValue
}

// Rotate90CW fait pivoter l'image PGM de 90 degrés dans le sens des aiguilles d'une montre.
func (pgm *PGM) Rotate90CW() {
	// Créer une nouvelle matrice pour stocker les données pivotées
	rotatedData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	// Remplir la nouvelle matrice pivotée en effectuant la rotation
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < pgm.height; j++ {
			rotatedData[i][j] = pgm.data[pgm.height-1-j][i]
		}
	}

	// Mettre à jour les dimensions de l'image après la rotation
	pgm.width, pgm.height = pgm.height, pgm.width

	// Mettre à jour les données de l'image avec les données pivotées
	pgm.data = rotatedData
}

// PBM représente une image PBM (Portable Bitmap).

// ToPBM convertit l'image PGM en PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		width:  pgm.width,
		height: pgm.height,
		data:   make([][]bool, pgm.height),
	}

	// Convertir les valeurs des pixels de l'image PGM en valeurs booléennes pour l'image PBM
	for i := 0; i < pgm.height; i++ {
		pbm.data[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			pbm.data[i][j] = pgm.data[i][j] > uint8(pgm.max)/2
		}
	}

	return pbm
}

// SavePBM enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème.
func (pbm *PBM) SavePBM(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire l'en-tête PBM
	_, err = file.WriteString(fmt.Sprintf("P1\n%d %d\n", pbm.width, pbm.height))
	if err != nil {
		return err
	}

	// Écrire les données de l'image PBM
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			bit := 0
			if pbm.data[i][j] {
				bit = 1
			}
			_, err := file.WriteString(fmt.Sprintf("%d ", bit))
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	filename := "duck.pgm" // Remplacez cela par le chemin de votre fichier PGM
	pgm, err := ReadPGM("duck.pgm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}
	fmt.Println("filename:", filename)
	width, height := pgm.Size()
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)
	fmt.Println(pgm.data)
	/*
		// Exemple d'utilisation de la fonction Set pour définir la valeur du pixel à la position (2, 3)
		x, y := 2, 3
		newValue := uint8(100)
		pgm.Set(x, y, newValue)

		// Vérifier la nouvelle valeur avec la fonction At
		fmt.Printf("Nouvelle valeur du pixel à la position (%d, %d): %d\n", x, y, pgm.At(x, y))

		// Exemple d'utilisation de la fonction SetMagicNumber pour définir le nombre magique
		newMagicNumber := "P2" // Remplacez cela par le nouveau nombre magique souhaité
		pgm.SetMagicNumber(newMagicNumber)

		// Définir la nouvelle valeur maximale
		newMaxValue := 255
		pgm.SetMaxValue(newMaxValue)

		// Enregistrer l'image PGM après modification
		err = pgm.Save("modified_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}

		// Enregistrer l'image PGM après modification
		err = pgm.Save("modified_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}
		// Exemple d'utilisation de la fonction Invert pour inverser les couleurs
		pgm.Invert()

		// Enregistrer l'image PGM après inversion des couleurs
		err = pgm.Save("inverted_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}

		fmt.Println("Image inversée et enregistrée avec succès.")
		// Exemple d'utilisation de la fonction Flip pour inverser horizontalement l'image
		pgm.Flip()

		// Enregistrer l'image PGM après inversion horizontale
		err = pgm.Save("flipped_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}

		fmt.Println("Image inversée horizontalement et enregistrée avec succès.")
		// Exemple d'utilisation de la fonction Flop pour inverser verticalement l'image
		pgm.Flop()

		// Enregistrer l'image PGM après inversion verticale
		err = pgm.Save("flopped_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}

		fmt.Println("Image inversée verticalement et enregistrée avec succès.")

		// Exemple d'utilisation de la fonction Rotate90CW pour pivoter l'image de 90° dans le sens des aiguilles d'une montre
		pgm.Rotate90CW()

		// Enregistrer l'image PGM après la rotation
		err = pgm.Save("90CW_duck.pgm")
		if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier:", err)
			return
		}

		fmt.Println("Image pivotée de 90° dans le sens des aiguilles d'une montre et enregistrée avec succès.")
	*/
	// Exemple d'utilisation de la fonction ToPBM pour convertir l'image PGM en PBM
	pbm := pgm.ToPBM()

	// Enregistrer l'image PBM après la conversion
	err = pbm.SavePBM("converted_duck.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement du fichier PBM:", err)
		return
	}

	fmt.Println("Image PGM convertie en PBM et enregistrée avec succès.")

}
