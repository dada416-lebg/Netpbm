package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aquilax/go-perlin"
)

type PPM struct {
	data        [][]Pixel
	width       int
	height      int
	magicNumber string
	max         int
}

type Pixel struct {
	R, G, B uint8
}

// ReadPPM lit une image PPM à partir d'un fichier et renvoie une structure représentant l'image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ppm PPM

	// Lire l'en-tête
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignorer les commentaires
			continue
		}

		switch {
		case ppm.magicNumber == "":
			ppm.magicNumber = line
		case ppm.width == 0:
			// Lire les dimensions
			dimensions := strings.Fields(line)
			if len(dimensions) != 2 {
				return nil, fmt.Errorf("format d'en-tête incorrect")
			}
			ppm.width, _ = strconv.Atoi(dimensions[0])
			ppm.height, _ = strconv.Atoi(dimensions[1])
		case ppm.max == 0:
			// Lire la valeur maximale des couleurs
			ppm.max, _ = strconv.Atoi(line)
			break
		}

		if ppm.width > 0 && ppm.height > 0 && ppm.max > 0 {
			break
		}
	}

	if ppm.magicNumber != "P3" && ppm.magicNumber != "P6" {
		return nil, fmt.Errorf("format ppm non pris en charge : %s", ppm.magicNumber)
	}

	// Initialiser le tableau ppm.data
	ppm.data = make([][]Pixel, ppm.height)

	// Lire les données
	for i := 0; i < ppm.height; i++ {
		ppm.data[i] = make([]Pixel, ppm.width)
		for j := 0; j < ppm.width; j++ {
			var pixel Pixel

			if ppm.magicNumber == "P3" {
				// Format P3 (ASCII)
				_, err := fmt.Fscanf(strings.NewReader(scanner.Text()), "%d %d %d", &pixel.R, &pixel.G, &pixel.B)
				if err != nil {
					return nil, err
				}
			} else if ppm.magicNumber == "P6" {
				// Format P6 (binaire)
				var pixelData [3]uint8
				n, err := file.Read(pixelData[:])
				if err != nil {
					return nil, err
				}

				if n < len(pixelData) {
					return nil, fmt.Errorf("données P6 incomplètes")
				}

				pixel = Pixel{pixelData[0], pixelData[1], pixelData[2]}
			}

			ppm.data[i][j] = pixel
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ppm, nil
}

// Size retourne la largeur et la hauteur de l'image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At retourne la valeur du pixel à la position (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set définit la valeur du pixel à la position (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save enregistre l'image PPM dans un fichier et retourne une erreur en cas de problème.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire l'en-tête PPM
	_, err = fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	if err != nil {
		return err
	}

	// Écrire les données pixel
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			_, err := fmt.Fprintf(file, "%d %d %d\n", pixel.R, pixel.G, pixel.B)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Invert inverse les couleurs de l'image PPM.
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := &ppm.data[i][j]

			// Inverser les composantes de couleur
			pixel.R = 255 - pixel.R
			pixel.G = 255 - pixel.G
			pixel.B = 255 - pixel.B
		}
	}
}

// Flip retourne l'image PPM horizontalement.
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		// Initialisez un tableau temporaire pour stocker la ligne inversée
		reversedRow := make([]Pixel, ppm.width)

		// Inversez la ligne
		for j := 0; j < ppm.width; j++ {
			reversedRow[j] = ppm.data[i][ppm.width-1-j]
		}

		// Copiez la ligne inversée dans l'image d'origine
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j] = reversedRow[j]
		}
	}
}

// Flop retourne l'image PPM verticalement.
func (ppm *PPM) Flop() {
	// Initialisez un tableau temporaire pour stocker l'image inversée
	reversedImage := make([][]Pixel, ppm.height)

	// Inversez l'image
	for i := 0; i < ppm.height; i++ {
		reversedImage[i] = ppm.data[ppm.height-1-i]
	}

	// Copiez l'image inversée dans l'image d'origine
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j] = reversedImage[i][j]
		}
	}
}

// SetMagicNumber définit le nombre magique de l'image PPM.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PPM.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PPM de 90° dans le sens des aiguilles d'une montre.
func (ppm *PPM) Rotate90CW() {
	// Initialisez un tableau temporaire pour stocker l'image pivotée
	rotatedImage := make([][]Pixel, ppm.width)
	for i := 0; i < ppm.width; i++ {
		rotatedImage[i] = make([]Pixel, ppm.height)
	}

	// Faites la rotation
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotatedImage[j][ppm.height-1-i] = ppm.data[i][j]
		}
	}

	// Mettez à jour les dimensions de l'image
	ppm.width, ppm.height = ppm.height, ppm.width

	// Copiez l'image pivotée dans l'image d'origine
	ppm.data = rotatedImage
}

// PGM structure for grayscale images
type PGM struct {
	data   [][]uint8
	width  int
	height int
	max    int
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Initialize a new PGM structure
	pgm := &PGM{
		data:   make([][]uint8, ppm.height),
		width:  ppm.width,
		height: ppm.height,
		max:    255, // PGM uses a max value of 255 for grayscale
	}

	// Initialize the data array for PGM
	for i := 0; i < ppm.height; i++ {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	// Convert RGB to grayscale using the weighted sum (commonly used weights)
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			grayValue := uint8(0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B))
			pgm.data[i][j] = grayValue
		}
	}

	return pgm
}

// Structure PBM pour les images binaires
type PBM struct {
	data   [][]bool
	width  int
	height int
}

// ToPBM convertit l'image PPM en PBM.
func (ppm *PPM) ToPBM(seuil uint8) *PBM {
	// Initialisez une nouvelle structure PBM
	pbm := &PBM{
		data:   make([][]bool, ppm.height),
		width:  ppm.width,
		height: ppm.height,
	}

	// Initialisez le tableau de données pour PBM
	for i := 0; i < ppm.height; i++ {
		pbm.data[i] = make([]bool, ppm.width)
	}

	// Convertissez les valeurs de niveaux de gris en binaire en fonction du seuil
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			valeurGris := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)

			// Vérifiez si la valeur de niveaux de gris est supérieure au seuil
			if uint8(valeurGris) >= seuil {
				pbm.data[i][j] = true
			} else {
				pbm.data[i][j] = false
			}
		}
	}

	return pbm
}

// Structure PBM pour les images binaires
type PBM struct {
	data   [][]bool
	width  int
	height int
}

// ToPBM convertit l'image PPM en PBM.
func (ppm *PPM) ToPBM(seuil uint8) *PBM {
	// Initialisez une nouvelle structure PBM
	pbm := &PBM{
		data:   make([][]bool, ppm.height),
		width:  ppm.width,
		height: ppm.height,
	}

	// Initialisez le tableau de données pour PBM
	for i := 0; i < ppm.height; i++ {
		pbm.data[i] = make([]bool, ppm.width)
	}

	// Convertissez les valeurs de niveaux de gris en binaire en fonction du seuil
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			pixel := ppm.data[i][j]
			valeurGris := 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)

			// Vérifiez si la valeur de niveaux de gris est supérieure au seuil
			if uint8(valeurGris) >= seuil {
				pbm.data[i][j] = true
			} else {
				pbm.data[i][j] = false
			}
		}
	}

	return pbm
}

type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// ...
}

// DrawRectangle dessine un rectangle dans l'image PPM.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Vérifier que les dimensions du rectangle sont valides
	if width <= 0 || height <= 0 {
		fmt.Println("Les dimensions du rectangle ne sont pas valides.")
		return
	}

	// Coordonnées du coin supérieur gauche du rectangle
	x1, y1 := p1.X, p1.Y
	// Coordonnées du coin inférieur droit du rectangle
	x2, y2 := x1+width-1, y1+height-1

	// Assurez-vous que les coordonnées du coin inférieur droit sont dans les limites de l'image
	if x2 >= ppm.width || y2 >= ppm.height {
		fmt.Println("Les dimensions du rectangle dépassent les limites de l'image.")
		return
	}

	// Dessiner le rectangle en définissant la couleur des pixels correspondants
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			ppm.data[y][x] = color
		}
	}
}

// DrawFilledRectangle dessine un rectangle rempli dans l'image PPM.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Vérifier que les dimensions du rectangle sont valides
	if width <= 0 || height <= 0 {
		fmt.Println("Les dimensions du rectangle ne sont pas valides.")
		return
	}

	// Coordonnées du coin supérieur gauche du rectangle
	x1, y1 := p1.X, p1.Y
	// Coordonnées du coin inférieur droit du rectangle
	x2, y2 := x1+width-1, y1+height-1

	// Assurez-vous que les coordonnées du coin inférieur droit sont dans les limites de l'image
	if x2 >= ppm.width || y2 >= ppm.height {
		fmt.Println("Les dimensions du rectangle dépassent les limites de l'image.")
		return
	}

	// Dessiner le rectangle rempli en définissant la couleur des pixels correspondants
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			ppm.data[y][x] = color
		}
	}
}

// DrawCircle dessine un cercle dans l'image PPM.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// Vérifier que le rayon est positif
	if radius <= 0 {
		fmt.Println("Le rayon du cercle doit être positif.")
		return
	}

	// Coordonnées du centre du cercle
	x0, y0 := center.X, center.Y

	// Initialiser les coordonnées
	x := radius
	y := 0
	err := 0

	// Dessiner le cercle en utilisant l'algorithme de tracé de cercle basé sur la méthode de balayage
	for x >= y {
		// Définir les pixels symétriques correspondants sur les huit octants du cercle
		ppm.Set(x0+x, y0-y, color)
		ppm.Set(x0+y, y0-x, color)
		ppm.Set(x0-y, y0-x, color)
		ppm.Set(x0-x, y0-y, color)
		ppm.Set(x0-x, y0+y, color)
		ppm.Set(x0-y, y0+x, color)
		ppm.Set(x0+y, y0+x, color)
		ppm.Set(x0+x, y0+y, color)

		// Mettre à jour les coordonnées
		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// DrawFilledCircle dessine un cercle plein dans l'image PPM.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Vérifier que le rayon est positif
	if radius <= 0 {
		fmt.Println("Le rayon du cercle doit être positif.")
		return
	}

	// Coordonnées du centre du cercle
	x0, y0 := center.X, center.Y

	// Initialiser les coordonnées
	x := radius
	y := 0
	err := 0

	// Dessiner le cercle en utilisant l'algorithme de tracé de cercle basé sur la méthode de balayage
	for x >= y {
		// Dessiner une ligne horizontale entre les deux points symétriques du cercle
		for i := x0 - x; i <= x0+x; i++ {
			ppm.Set(i, y0+y, color)
			ppm.Set(i, y0-y, color)
		}

		// Mettre à jour les coordonnées
		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// DrawTriangle dessine un triangle dans l'image PPM.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// Trier les points par ordonnée croissante
	points := []Point{p1, p2, p3}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Coordonnées des sommets triés
	y1, y2, y3 := points[0].Y, points[1].Y, points[2].Y
	x1, x2, x3 := points[0].X, points[1].X, points[2].X

	// Calculer les pentes des côtés du triangle
	slope1 := float64(x2-x1) / float64(y2-y1)
	slope2 := float64(x3-x1) / float64(y3-y1)
	slope3 := float64(x3-x2) / float64(y3-y2)

	// Initialiser les coordonnées du bord gauche et du bord droit
	leftX := x1
	rightX := x1

	// Boucler sur chaque ligne scanline du triangle
	for y := y1; y <= y3; y++ {
		// Dessiner la ligne horizontale entre les bords gauche et droit
		for x := leftX; x <= rightX; x++ {
			ppm.Set(x, y, color)
		}

		// Mettre à jour les bords gauche et droit
		if y < y2 {
			leftX += int(slope1)
			rightX += int(slope2)
		} else {
			leftX += int(slope1)
			rightX += int(slope3)
		}
	}
}

// DrawFilledTriangle dessine un triangle rempli dans l'image PPM.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Trier les points par ordonnée croissante
	points := []Point{p1, p2, p3}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Coordonnées des sommets triés
	y1, y2, y3 := points[0].Y, points[1].Y, points[2].Y
	x1, x2, x3 := points[0].X, points[1].X, points[2].X

	// Calculer les pentes des côtés du triangle
	slope1 := float64(x2-x1) / float64(y2-y1)
	slope2 := float64(x3-x1) / float64(y3-y1)
	slope3 := float64(x3-x2) / float64(y3-y2)

	// Initialiser les coordonnées du bord gauche et du bord droit
	leftX := x1
	rightX := x1

	// Boucler sur chaque ligne scanline du triangle
	for y := y1; y <= y3; y++ {
		// Dessiner la ligne horizontale entre les bords gauche et droit
		for x := leftX; x <= rightX; x++ {
			ppm.Set(x, y, color)
		}

		// Mettre à jour les bords gauche et droit
		if y < y2 {
			leftX += int(slope1)
			rightX += int(slope2)
		} else {
			leftX += int(slope1)
			rightX += int(slope3)
		}
	}
}

// DrawPolygon dessine un polygone rempli dans l'image PPM.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Vérifier si le nombre de points est suffisant pour former un polygone
	if len(points) < 3 {
		fmt.Println("Le polygone doit avoir au moins 3 points.")
		return
	}

	// Trier les points par ordonnée croissante
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Trouver les coordonnées du polygone dans la boîte englobante
	minY := points[0].Y
	maxY := points[len(points)-1].Y

	// Initialiser le tableau des bords gauche et droit
	leftX := make([]int, maxY-minY+1)
	rightX := make([]int, maxY-minY+1)

	// Remplir le tableau des bords gauche et droit en utilisant l'algorithme de balayage
	for i := range leftX {
		leftX[i] = ppm.width
		rightX[i] = 0
	}

	for i := 0; i < len(points); i++ {
		curr := points[i]
		next := points[(i+1)%len(points)]

		// Dessiner une ligne entre les points actuels et suivants
		ppm.drawLineDDA(curr.X, curr.Y, next.X, next.Y, color, &leftX, &rightX)
	}

	// Remplir le polygone en utilisant les bords gauche et droit
	for y := minY; y <= maxY; y++ {
		for x := leftX[y-minY]; x <= rightX[y-minY]; x++ {
			ppm.Set(x, y, color)
		}
	}
}

// drawLineDDA utilise l'algorithme DDA pour dessiner une ligne entre deux points.
func (ppm *PPM) drawLineDDA(x1, y1, x2, y2 int, color Pixel, leftX, rightX *[]int) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))

	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)

	x := float64(x1)
	y := float64(y1)

	for i := 0; i <= steps; i++ {
		// Mettre à jour les bords gauche et droit
		(*leftX)[int(y)] = int(math.Min(float64((*leftX)[int(y)]), x))
		(*rightX)[int(y)] = int(math.Max(float64((*rightX)[int(y)]), x))

		x += xIncrement
		y += yIncrement
	}
}

// DrawFilledPolygon dessine un polygone rempli dans l'image PPM.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Vérifier si le nombre de points est suffisant pour former un polygone
	if len(points) < 3 {
		fmt.Println("Le polygone doit avoir au moins 3 points.")
		return
	}

	// Trier les points par ordonnée croissante
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Trouver les coordonnées du polygone dans la boîte englobante
	minY := points[0].Y
	maxY := points[len(points)-1].Y

	// Initialiser le tableau des bords gauche et droit
	leftX := make([]int, maxY-minY+1)
	rightX := make([]int, maxY-minY+1)

	// Remplir le tableau des bords gauche et droit en utilisant l'algorithme de balayage
	for i := range leftX {
		leftX[i] = ppm.width
		rightX[i] = 0
	}

	for i := 0; i < len(points); i++ {
		curr := points[i]
		next := points[(i+1)%len(points)]

		// Dessiner une ligne entre les points actuels et suivants
		ppm.drawLineDDA(curr.X, curr.Y, next.X, next.Y, color, &leftX, &rightX)
	}

	// Remplir le polygone en utilisant les bords gauche et droit
	for y := minY; y <= maxY; y++ {
		for x := leftX[y-minY]; x <= rightX[y-minY]; x++ {
			ppm.Set(x, y, color)
		}
	}
}

// DrawKochSnowflake dessine un flocon de neige de Koch récursif.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, length int, color Pixel) {
	// Définir les angles pour les segments du flocon de neige
	angles := []float64{0, 120, -120, 0}

	// Dessiner la première section du flocon de neige
	ppm.drawKochSegment(n, start, length, angles, color)
}

// drawKochSegment dessine un segment du flocon de neige de Koch récursif.
func (ppm *PPM) drawKochSegment(n int, start Point, length int, angles []float64, color Pixel) {
	if n == 0 {
		// Cas de base : dessiner un segment
		end := Point{start.X + length, start.Y}
		ppm.DrawLine(start, end, color)
	} else {
		// Calculer les nouveaux points pour diviser le segment en trois parties égales
		third := length / 3
		p1 := start
		p2 := Point{start.X + third, start.Y}
		p3 := Point{start.X + third, start.Y + int(float64(third)*math.Sqrt(3))}
		p4 := Point{start.X + 2*third, start.Y + int(float64(third)*math.Sqrt(3))}
		p5 := Point{start.X + 2*third, start.Y}
		p6 := Point{start.X + 3*third, start.Y}

		// Dessiner les segments récursivement
		ppm.drawKochSegment(n-1, p1, third, angles, color)
		ppm.drawKochSegment(n-1, p2, third, angles, color)
		ppm.drawKochSegment(n-1, p3, third, angles, color)
		ppm.drawKochSegment(n-1, p4, third, angles, color)
		ppm.drawKochSegment(n-1, p5, third, angles, color)
		ppm.drawKochSegment(n-1, p6, third, angles, color)
	}
}

// DrawSierpinskiTriangle dessine un triangle de Sierpinski récursif.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	if n == 0 {
		// Cas de base : dessiner un triangle
		p1 := start
		p2 := Point{start.X + width, start.Y}
		p3 := Point{start.X + width/2, start.Y - int(float64(width)*math.Sqrt(3)/2)}

		ppm.DrawLine(p1, p2, color)
		ppm.DrawLine(p2, p3, color)
		ppm.DrawLine(p3, p1, color)
	} else {
		// Calculer les points pour diviser le triangle en trois parties égales
		thirdWidth := width / 2
		p1 := start
		p2 := Point{start.X + width, start.Y}
		p3 := Point{start.X + thirdWidth, start.Y - int(float64(width)*math.Sqrt(3)/2)}

		// Calculer le point au sommet du triangle intérieur
		topPoint := Point{start.X + thirdWidth/2, start.Y - int(float64(thirdWidth)*math.Sqrt(3)/2)}

		// Dessiner les triangles récursivement
		ppm.DrawSierpinskiTriangle(n-1, p1, thirdWidth, color)
		ppm.DrawSierpinskiTriangle(n-1, topPoint, thirdWidth, color)
		ppm.DrawSierpinskiTriangle(n-1, p2, thirdWidth, color)
		ppm.DrawSierpinskiTriangle(n-1, p3, thirdWidth, color)
	}
}

// DrawPerlinNoise draws Perlin noise on the image.
func (ppm *PPM) DrawPerlinNoise(color1, color2 Pixel) {
	// Set up Perlin noise generator
	p := perlin.NewPerlin(3, 3, 1, 42)

	// Loop through each pixel in the image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Generate Perlin noise value for the current pixel
			noiseValue := p.Noise2D(float64(x)/50, float64(y)/50)

			// Map the noise value to a color between color1 and color2
			color := interpolateColor(color1, color2, noiseValue)

			// Set the color for the current pixel
			ppm.Set(x, y, color)
		}
	}
}

// interpolateColor linearly interpolates between two colors based on a t value.
func interpolateColor(color1, color2 Pixel, t float64) Pixel {
	// Ensure t is within the [0, 1] range
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Linear interpolation for each color component
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)

	return Pixel{r, g, b}
}
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// Calculer les ratios de redimensionnement
	widthRatio := float64(ppm.width) / float64(newWidth)
	heightRatio := float64(ppm.height) / float64(newHeight)

	// Initialiser la nouvelle image redimensionnée
	resizedImage := make([][]Pixel, newHeight)
	for i := 0; i < newHeight; i++ {
		resizedImage[i] = make([]Pixel, newWidth)
	}

	// Appliquer l'algorithme des k-voisins les plus proches
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Coordonnées correspondantes dans l'image d'origine
			sourceX := int(float64(x) * widthRatio)
			sourceY := int(float64(y) * heightRatio)

			// S'assurer que les coordonnées restent dans les limites de l'image d'origine
			if sourceX >= ppm.width {
				sourceX = ppm.width - 1
			}
			if sourceY >= ppm.height {
				sourceY = ppm.height - 1
			}

			// Copier la valeur du pixel correspondant
			resizedImage[y][x] = ppm.data[sourceY][sourceX]
		}
	}

	// Mettre à jour les dimensions de l'image
	ppm.width = newWidth
	ppm.height = newHeight

	// Mettre à jour les données de l'image avec l'image redimensionnée
	ppm.data = resizedImage
}
