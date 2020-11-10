// Image processing functions, to provide transformation functions

package elputils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// Write an image object to a file
func ImageToFile(image image.Image, destination string) {
	// Create blank file
	output, err := os.Create(destination)
	if err != nil {
		fmt.Println("Couldn't create file", destination)
	}

	err = png.Encode(output, image)
	if err != nil {
		fmt.Println("Couldn't write image to file", destination)
	}

	_ = output.Close()
}

// Creates an image object from a file
func FileToImage(path string) image.Image {
	input, err := os.Open(path)
	if err != nil {
		fmt.Println("Coudln't open file", path)
	}

	img, err := jpeg.Decode(input)

	if err != nil {
		fmt.Println("Coudln't import image")
	}

	return img
}

// Converts image to Black & White
func GreyScale(img image.Image) image.Image {
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

// Converts image to black & white negative
func NegativeBW(img image.Image) image.Image {
	imgGris := GreyScale(img)
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

// Negative (RGB)
func NegativeRGB(img image.Image) image.Image {
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

// Filter using a 3x3 convolution matrix
func Convolution(x int, y int, img image.Image, coeff [3][3]float64, somme float64) color.RGBA {
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
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x256 car la fct rgba renvoie des uint32 multipliés par alpha
	return pix
}

func GaussMatrix(n int) ([][]float64, float64) {
	coeff := make([][]float64, n) //création d'une matrice vide de taille n*n
	for i := 0; i < n; i++ {
		coeff[i] = make([]float64, n)
	}
	var somme float64
	et := 0.75 //écart-type

	for x := -n / 2; x <= n/2; x++ { //calcul de la matrice de Convolution de gauss, c'est débile ça recalcule à chaque fois
		//mais j'arrive pas un faire une fonction modulable qui renvoie une matrice n*n
		for y := -n / 2; y <= n/2; y++ {
			coeff[x+n/2][y+n/2] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			somme += coeff[x+n/2][y+n/2]
		}
	}

	return coeff, somme
}

func ConvolutionGauss(x int, y int, img image.Image, n int, coeff [][]float64, somme float64) color.RGBA { //Convolution avec la fonction de gauss avec une matrice de taille n
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0

	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			rouge, vert, bleu, _ := img.At(x+i, y+j).RGBA()
			r += coeff[i+n/2][j+n/2] * float64(rouge)
			v += coeff[i+n/2][j+n/2] * float64(vert)
			b += coeff[i+n/2][j+n/2] * float64(bleu)

		}
	}
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x255 car la fct rgba renvoie des uint32 multipliés par alpha qui vaut 255 ici
	return pix
}

func GaussBlur(img image.Image, n int) image.Image { //flou de gauss avec matrice de taille n (plus n est grand plus l'effet est important), peut-etre faire plusieurs itérations
	imgFlou := image.NewRGBA(img.Bounds())
	coeff, somme := GaussMatrix(n)

	for y := img.Bounds().Min.Y + 2; y < img.Bounds().Max.Y-2; y++ {
		for x := img.Bounds().Min.X + 2; x < img.Bounds().Max.X-2; x++ {
			imgFlou.Set(x, y, ConvolutionGauss(x, y, img, n, coeff, somme))
		}
	}

	return imgFlou
}

func UniformBlur(img image.Image) image.Image { //applique un filtre uniforme 3x3 à chaque pixel
	imgFlou := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	somme := 9.0
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgFlou.Set(x, y, Convolution(x, y, img, coeff, somme))
		}
	}

	return imgFlou
}

func Boundaries(img image.Image, puissance int) image.Image { //filtre laplacien : fonctionne mieux, utiliser un autre filtre? sobel, prewitt...
	imgCont := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //laplacien
	//coeff := [3][3]float64{{-2,-2,0},{-2,0,2},{0,2,2}} //sobel
	somme := float64(puissance) //pose problème entre nv détails , et efficacité ==> eventuellement le proposer en réglage à l'utilisateur
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont.Set(x, y, Convolution(x, y, img, coeff, somme))
		}
	}

	return NegativeBW(imgCont) //applique le filtre négatif pour que ça soit "plus joli"
}

func PrewittBorders(img image.Image, puissance int) image.Image { //filtre prewitt
	imgCont0 := image.NewRGBA(img.Bounds())
	imgCont90 := image.NewRGBA(img.Bounds())
	imgRes := image.NewRGBA(img.Bounds())

	coeff1 := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //Prewitt 0°
	coeff2 := [3][3]float64{{-2, -2, 0}, {-2, 0, 2}, {0, 2, 2}}      //Prewitt 90°

	somme := float64(puissance) //pose problème entre nv détails , et efficacité ==> eventuellement le proposer en réglage à l'utilisateur
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont0.Set(x, y, Convolution(x, y, img, coeff1, somme))
		}
	}

	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont90.Set(x, y, Convolution(x, y, img, coeff2, somme))
		}
	}

	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			pix := uint8(math.Sqrt(math.Pow(float64(imgCont0.RGBAAt(x, y).R), 2) + math.Pow(float64(imgCont90.RGBAAt(x, y).R), 2)))
			imgRes.Set(x, y, color.Gray{Y: pix})
		}
	}

	return NegativeBW(imgRes) //applique le filtre négatif pour que ça soit "plus joli"
}

func neat(img image.Image) image.Image { //pb à voir fait de la merde, n'est pas référencé dans le Dispatch
	imgNet := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}}
	somme := 1.0
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgNet.Set(x, y, Convolution(x, y, img, coeff, somme))
		}
	}

	return imgNet
}

func separation(img image.Image) (image.Image, image.Image, image.Image) { //sert à séparer les 3 composantes de l'image_utils, pas très utile
	imgR := image.NewRGBA(img.Bounds())
	imgV := image.NewRGBA(img.Bounds())
	imgB := image.NewRGBA(img.Bounds())

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			r, v, b, _ := img.At(x, y).RGBA()
			imgR.Set(x, y, color.RGBA{R: uint8(r)})
			imgV.Set(x, y, color.RGBA{G: uint8(v)})
			imgB.Set(x, y, color.RGBA{B: uint8(b)})
		}
	}
	return imgR, imgV, imgB
}

func DespeckleBW(img image.Image, x int, y int, n int, coeffGauss [][]float64, sommeGauss float64) color.Gray { //réduction de bruit Noir et Blanc : fonctionnel
	var moyenne, stdev float64
	moyenne, stdev = 0, 0

	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on évite le pixel central
				t, _, _, _ := img.At(x+i, y+j).RGBA()
				moyenne += float64(t)

			}
		}
	}

	moyenne = moyenne / (math.Pow(float64(n), 2) - 1)

	//on calcule l'écart-type
	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on exclut le pixel central
				t, _, _, _ := img.At(x+i, y+j).RGBA()
				stdev += math.Pow(float64(t)-moyenne, 2)
			}
		}
	}

	stdev = math.Sqrt(stdev / (math.Pow(float64(n), 2) - 1))

	//on teste ensuite le pixel central s'il est dans l'intervalle moyenne +/- écart-type sinon on applique un filtre de gauss 3x3 à ce pixel
	t, _, _, _ := img.At(x, y).RGBA()

	if float64(t) <= moyenne-stdev || float64(t) >= moyenne+stdev {

		t = uint32(ConvolutionGauss(x, y, img, 3, coeffGauss, sommeGauss).R) //on prend la composante rouge ici mais peu importe vu qu'on est en niveaux de gris
	}

	pix := color.Gray{Y: uint8(t)} //on reconstruit le pixel modifié
	return pix

}

// MARCHE PAS en couleurs
func DespeckleRGB(img image.Image, x int, y int, n int, coeffGauss [][]float64, sommeGauss float64) color.RGBA { //déparasitage de chaque composante, n taille matrice
	var moyenneR, moyenneV, moyenneB float64
	moyenneB, moyenneR, moyenneV = 0, 0, 0

	var stdevR, stdevV, stdevB float64 //écart-type
	stdevR, stdevV, stdevB = 0, 0, 0

	//on calcule la moyenne de chaque pixel voisin au pixel en question pour chaque composante
	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on exclut le pixel central
				r, v, b, _ := img.At(x+i, y+j).RGBA()
				moyenneR += float64(r)
				moyenneV += float64(v)
				moyenneB += float64(b)
			}
		}
	}

	moyenneR, moyenneV, moyenneB = moyenneR/(math.Pow(float64(n), 2)-1), moyenneV/(math.Pow(float64(n), 2)-1), moyenneB/(math.Pow(float64(n), 2)-1)

	//on calcule l'écart-type pour chaque composante
	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on exclut le pixel central
				r, v, b, _ := img.At(x+i, y+j).RGBA()
				stdevR += math.Pow(float64(r)-moyenneR, 2)
				stdevV += math.Pow(float64(v)-moyenneV, 2)
				stdevB += math.Pow(float64(b)-moyenneB, 2)
			}
		}
	}

	stdevR, stdevV, stdevB = math.Sqrt(stdevR/(math.Pow(float64(n), 2)-1)), math.Sqrt(stdevV/(math.Pow(float64(n), 2)-1)), math.Sqrt(stdevB/(math.Pow(float64(n), 2)-1))

	//on teste ensuite le pixel central s'il est dans l'intervalle moyenne +/- écart-type sinon il prend on applique un filtre gaussien sur le pixel
	r, v, b, _ := img.At(x, y).RGBA()

	if float64(r) <= moyenneR-stdevR || float64(r) >= moyenneR+stdevR {
		r = uint32(ConvolutionGauss(x, y, img, 3, coeffGauss, sommeGauss).R)
	}
	if float64(v) <= moyenneV-stdevV || float64(v) >= moyenneV+stdevV {
		v = uint32(ConvolutionGauss(x, y, img, 3, coeffGauss, sommeGauss).G)
	}
	if float64(b) <= moyenneB-stdevB || float64(b) >= moyenneR+stdevB {
		b = uint32(ConvolutionGauss(x, y, img, 3, coeffGauss, sommeGauss).B)
	}

	pix := color.RGBA{R: uint8(r), G: uint8(v), B: uint8(b), A: 0xff} //on reconstruit le pixel modifié
	return pix
}

func NoiseReductionBW(img image.Image, nbIterations int, n int) image.Image {
	imgDebruit := image.NewGray(img.Bounds())
	coeffGauss, sommeGauss := GaussMatrix(3)

	for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
		for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
			imgDebruit.Set(x, y, DespeckleBW(img, x, y, n, coeffGauss, sommeGauss))
		}
	}

	for k := 0; k < nbIterations-1; k++ {
		for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
			for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
				imgDebruit.Set(x, y, DespeckleBW(imgDebruit, x, y, n, coeffGauss, sommeGauss))
			}
		}
	}

	return imgDebruit
}

func NoiseReductionRGB(img image.Image, nbIterations int, n int) image.Image { //ne marche pas bien, n'est pas référencé dans le Dispatch
	imgDebruit := image.NewRGBA(img.Bounds())
	coeffGauss, sommeGauss := GaussMatrix(3)

	for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
		for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
			imgDebruit.Set(x, y, DespeckleRGB(img, x, y, n, coeffGauss, sommeGauss))
		}
	}

	for k := 0; k < nbIterations-1; k++ {
		for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
			for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
				imgDebruit.Set(x, y, DespeckleRGB(imgDebruit, x, y, n, coeffGauss, sommeGauss))
			}
		}
	}

	return imgDebruit
}

var FilterList = []string{
	"Negative Black & white",
	"Negative RGB",
	"Grey scale",
	"Uniform Blur",
	"Gauss Blur",
	"Noise reduction",
	"Boundary detection",
	"Boundaries with Prewitt",
}

func Dispatch(img image.Image, n int) image.Image { //permet de sélectionner quelle transfo en fonction de l'entrée utilisateur du programme principal
	//variable param pour les traitement nécessitant un niveau de puissance selectionné par l'utilisateur
	var res image.Image

	switch n {

	case 1:
		res = NegativeBW(img)
		break

	case 2:
		res = NegativeRGB(img)
		break

	case 3:
		res = GreyScale(img)
		break

	case 4:
		res = UniformBlur(img)
		break

	case 5:
		res = GaussBlur(img, 5) //param1=taille matrice /!\ nb impairs= puissance du flou en fct de la résolution et de l'envie de l'utilisateur
		break

	case 6:
		res = NoiseReductionBW(img, 1, 5) //param1 = nbIteration du filtre 1 à 2 conseillé
		//param2 = taille matrice 5 conseillé (nombres impairs /!\)
		break

	case 7:
		res = Boundaries(img, 8) //param1 = puissance de séparation en paramètre 8 voire 16, plus c'est élevé plus seuls les "gros" Boundaries seront visibles
		break

	case 8:
		res = PrewittBorders(img, 8) //pareil que Boundaries mais utilise le filtre de Prewitt à la place
		break

	}

	return res
}
