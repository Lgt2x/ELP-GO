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

func negatifRVB(image2 image.Image) image.Image { //pb fait de la merde
	imgNeg := image.NewRGBA(image2.Bounds())
	for y := image2.Bounds().Min.Y; y < image2.Bounds().Max.Y; y++ {
		for x := image2.Bounds().Min.X; x < image2.Bounds().Max.X; x++ {
			r, v, b, _ := image2.At(x, y).RGBA()
			pix := color.RGBA{R: 255 - uint8(r), G: 255 - uint8(v), B: 255 - uint8(b), A: 0xff}
			imgNeg.Set(x, y, pix)
		}
	}
	return imgNeg
}

func gauss(et float64) [3][3]uint8 {
	var res [3][3]uint8

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			res[x][y] = uint8(100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2)))
		}
	}
	return res //calcul de la matrice de convolution du filtre de gauss en fonction de l'écart-type
}

func convolution(x int, y int, img image.Image, coeff [3][3]uint8) color.RGBA {
	var pix color.RGBA
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			rouge, vert, bleu, _ := img.At(y+i, x+j).RGBA()
			r := coeff[i+1][j+1] * uint8(rouge) / 9.0
			v := coeff[i+1][j+1] * uint8(vert) / 9.0
			b := coeff[i+1][j+1] * uint8(bleu) / 9.0
			pix = color.RGBA{r, v, b, 0xff}
		}
	}
	return pix
}

func flouGauss(image2 image.Image) image.Image { //marche pas
	imgFlou := image.NewRGBA(image2.Bounds())
	coeff := gauss(1.5)
	for y := image2.Bounds().Min.Y + 2; y < image2.Bounds().Max.Y-2; y++ {
		for x := image2.Bounds().Min.X + 2; x < image2.Bounds().Max.X-2; x++ {
			imgFlou.Set(x, y, convolution(x, y, image2, coeff))
		}
	}

	return imgFlou
}

func main() {

	image := importImage()
	//niveauGris(image) //renvoie l'image traitée dans le répertoire courant
	ecritureFichier(flouGauss(image))
}
