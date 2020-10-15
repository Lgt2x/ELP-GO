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
	if err != nil {
		fmt.Println("Erreur dans la création du fichier!")
	}
	png.Encode(output, image2) //écriture de l'image
	output.Close()
}

func importImage() image.Image {
	input, err := os.Open("C:/Users/antoi/pictures/hawkeye.jpg")
	img, err := jpeg.Decode(input)

	if err != nil {
		fmt.Println("Erreur dans l'importation de l'image!")
	}

	return img
}

func niveauGris(imge image.Image) image.Image {
	imgGris := image.NewGray(imge.Bounds())

	for y := imge.Bounds().Min.Y; y < imge.Bounds().Max.Y; y++ {
		for x := imge.Bounds().Min.X; x < imge.Bounds().Max.X; x++ {
			imgGris.Set(x, y, imge.At(x, y))

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

func negatifNB(image2 image.Image) image.Image {
	imgGris := niveauGris(image2)
	imgNeg := image.NewGray(image2.Bounds())
	for y := image2.Bounds().Min.Y; y < image2.Bounds().Max.Y; y++ {
		for x := image2.Bounds().Min.X; x < image2.Bounds().Max.X; x++ {
			z, _, _, _ := imgGris.At(x, y).RGBA()
			pix := color.Gray{Y: 255 - uint8(z)}
			imgNeg.Set(x, y, pix)
		}
	}
	return imgNeg
}

func negatifRVB(image2 image.Image) image.Image {
	imgNeg := image.NewRGBA(image2.Bounds())
	for y := image2.Bounds().Min.Y; y < image2.Bounds().Max.Y; y++ {
		for x := image2.Bounds().Min.X; x < image2.Bounds().Max.X; x++ {
			r, v, b, _ := image2.At(x, y).RGBA()
			pix := color.RGBA{R: 255 - uint8(r/255), G: 255 - uint8(v/255), B: 255 - uint8(b/255), A: 0xff}
			imgNeg.Set(x, y, pix)
		}
	}
	return imgNeg
}

func gauss(et float64) ([3][3]float64, float64) {
	var res [3][3]float64
	//var sommeLignes [3]float64
	var somme float64

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			res[x+1][y+1] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			//sommeLignes[x+1] += res[x+1][y+1]
			somme += res[x+1][y+1]
		}
	}

	fmt.Println(somme)
	/*
		for a:=0; a<3; a++{ //on fait la somme des lignes, y'a t-il besoin de la normaliser?
			for b:=0; b<3; b++ {
				sommeLignes[a] += res[a][b]
			}
		}


		for i:=0; i<3; i++{
			for j:=0; j<3; j++ {
				res[i][j] = res[i][j]/sommeLignes[i] //normlisation de la matrice
			}
		}



	*/
	return res, somme //calcul de la matrice de convolution du filtre de gauss en fonction de l'écart-type
}

func convolution(x int, y int, img image.Image, coeff [3][3]float64, somme float64) color.RGBA {
	//fmt.Println(coeff)
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
	pix = color.RGBA{R: uint8(r / (255 * somme)), G: uint8(v / (255 * somme)), B: uint8(b / (255 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x255 car la fct rgba renvoie des uint32 multipliés par alpha qui vaut 255 ici
	//fmt.Println(pix)
	return pix
}

func flouGauss(image2 image.Image) image.Image {
	imgFlou := image.NewRGBA(image2.Bounds())
	coeff := [3][3]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	somme := 9.0
	//coeff, somme := gauss(0.95)
	for y := image2.Bounds().Min.Y + 1; y < image2.Bounds().Max.Y-1; y++ {
		for x := image2.Bounds().Min.X + 1; x < image2.Bounds().Max.X-1; x++ {
			imgFlou.Set(x, y, convolution(x, y, image2, coeff, somme))
		}
	}

	return imgFlou
}

func main() {

	image := importImage()
	//niveauGris(image) //renvoie l'image traitée dans le répertoire courant
	ecritureFichier(negatifRVB(image))
}
