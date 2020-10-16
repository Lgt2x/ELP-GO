package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func ecritureFichier(image2 image.Image) {
	output, err := os.Create("img_modif.png") //création du fichier image de sortie

	err = png.Encode(output, image2) //écriture de l'image
	err = output.Close()
	if err != nil {
		fmt.Println("Erreur dans la création du fichier!")
	}
}

func importImage() image.Image {
	input, err := os.Open("C:/Users/antoi/pictures/hawkeye.jpg")
	img, err := jpeg.Decode(input)

	if err != nil {
		fmt.Println("Erreur dans l'importation de l'image!")
	}

	return img
}

func niveauGris(img image.Image) image.Image {
	imgGris := image.NewGray(img.Bounds())

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			imgGris.Set(x, y, img.At(x, y))

			/* //si on veut vraiment flex
			R, G, B, _ := img.At(x, y).RGBA()
			        //Luma: Y = 0.2126*R + 0.7152*G + 0.0722*B
			        Y := (0.2126*float64(R) + 0.7152*float64(G) + 0.0722*float64(B)) * (255.0 / 65535)
			        grayPix := color.Gray{uint8(Y)}
			*/
		}
	}
	return imgGris
}

func negatifNB(img image.Image) image.Image {
	imgGris := niveauGris(img)
	imgNeg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			z, _, _, _ := imgGris.At(x, y).RGBA()
			pix := color.Gray{Y: 255 - uint8(z)}
			imgNeg.Set(x, y, pix)
		}
	}
	return imgNeg
}

func negatifRVB(img image.Image) image.Image {
	imgNeg := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, v, b, _ := img.At(x, y).RGBA()
			pix := color.RGBA{R: 255 - uint8(r/256), G: 255 - uint8(v/256), B: 255 - uint8(b/256), A: 0xff}
			imgNeg.Set(x, y, pix)
		}
	}
	return imgNeg
}

func gauss3(et float64) ([3][3]float64, float64) {
	var res [3][3]float64
	var somme float64

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			res[x+1][y+1] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			somme += res[x+1][y+1]
		}
	}
	fmt.Println(res)
	fmt.Println(somme)
	return res, somme //calcul de la matrice de convolution du filtre de gauss en fonction de l'écart-type
}

func gauss5(et float64) ([5][5]float64, float64) {
	var res [5][5]float64
	var somme float64

	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			res[x+2][y+2] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			somme += res[x+2][y+2]
		}
	}
	fmt.Println(res)
	fmt.Println(somme)
	return res, somme //calcul de la matrice de convolution du filtre de gauss en fonction de l'écart-type
}

func convolution(x int, y int, img image.Image, coeff [3][3]float64, somme float64) color.RGBA {
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			rouge, vert, bleu, _ := img.At(x+i, y+j).RGBA()
			r += coeff[i+1][j+1] * float64(rouge)
			v += coeff[i+1][j+1] * float64(vert)
			b += coeff[i+1][j+1] * float64(bleu)

		}
	}
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x255 car la fct rgba renvoie des uint32 multipliés par alpha qui vaut 255 ici
	return pix
}

func convolution5(x int, y int, img image.Image, coeff [5][5]float64, somme float64) color.RGBA {
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			rouge, vert, bleu, _ := img.At(x+i, y+j).RGBA()
			r += coeff[i+2][j+2] * float64(rouge)
			v += coeff[i+2][j+2] * float64(vert)
			b += coeff[i+2][j+2] * float64(bleu)

		}
	}
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x255 car la fct rgba renvoie des uint32 multipliés par alpha qui vaut 255 ici
	return pix
}

func flouGauss3(img image.Image) image.Image {
	imgFlou := image.NewRGBA(img.Bounds())
	coeff, somme := gauss3(0.75)
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgFlou.Set(x, y, convolution(x, y, img, coeff, somme))
		}
	}

	return imgFlou
}

func flouGauss5(img image.Image) image.Image {
	imgFlou := image.NewRGBA(img.Bounds())
	coeff, somme := gauss5(0.75)
	for y := img.Bounds().Min.Y + 2; y < img.Bounds().Max.Y-2; y++ {
		for x := img.Bounds().Min.X + 2; x < img.Bounds().Max.X-2; x++ {
			imgFlou.Set(x, y, convolution5(x, y, img, coeff, somme))
		}
	}

	return imgFlou
}

func flouUniforme(img image.Image) image.Image {
	imgFlou := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	somme := 9.0
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgFlou.Set(x, y, convolution(x, y, img, coeff, somme))
		}
	}

	return imgFlou
}

func contours(img image.Image) image.Image { //fonctionne mais pas ultra ouf, utiliser un autre filtre? canny/sobel/laplacien...
	imgCont := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}}
	somme := 32.0
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont.Set(x, y, convolution(x, y, img, coeff, somme))
		}
	}

	return negatifNB(imgCont)
}

func nettete(img image.Image) image.Image { //pb à voir fait de la merde
	imgNet := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}}
	somme := 1.0
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgNet.Set(x, y, convolution(x, y, img, coeff, somme))
		}
	}

	return imgNet
}

func dispatch(img image.Image, n int) image.Image { //permet de sélectionner quelle transfo en fonction de l'entrée utilisateur du programme principal
	var res image.Image

	switch n {

	case 1:
		res = negatifNB(img)
		break

	case 2:
		res = negatifRVB(img)
		break

	case 3:
		res = niveauGris(img)
		break

	case 4:
		res = flouUniforme(img)
		break

	case 5:
		res = flouGauss3(img)
		break

	case 6:
		res = flouGauss5(img)
		break

	case 7:
		res = contours(img)
		break
	}
	return res
}

func main() {

	imageTest := importImage()
	//niveauGris(image) //renvoie l'image traitée dans le répertoire courant
	ecritureFichier(flouGauss5(imageTest))
}
