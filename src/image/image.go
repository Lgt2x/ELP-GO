package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

func ecritureFichier(image2 image.Image) {
	output, err := os.Create("img_modif.png") //création du fichier image de sortie

	err = png.Encode(output, image2) //écriture de l'image
	err = output.Close()
	if err != nil {
		fmt.Println("Erreur dans la création du fichier!")
	}
}

func importImage() image.Image {
	input, err := os.Open("C:/Users/antoi/pictures/platja_swag.jpg")
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

func negatifNB(img image.Image) image.Image { //image négatif en noir en blanc
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

func negatifRVB(img *image.Image, res *image.Image, wg *sync.WaitGroup, minX int, minY int, maxX int, maxY int) { //négatif couleurs
	//mu.Lock()
	imgNeg := image.NewRGBA((*img).Bounds())
	//mu.Unlock()
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			//mu.Lock()
			r, v, b, _ := (*img).At(x, y).RGBA()
			//mu.Unlock()
			pix := color.RGBA{R: 255 - uint8(r/256), G: 255 - uint8(v/256), B: 255 - uint8(b/256), A: 0xff}
			imgNeg.Set(x, y, pix)
		}
	}
	*res = imgNeg
	defer wg.Done()
}

func convolution(x int, y int, img image.Image, coeff [3][3]float64, somme float64) color.RGBA { //filtre à l'aide d'une matrice de convolution 3x3 un pixel
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

func calculGauss(n int) ([][]float64, float64) {

	coeff := make([][]float64, n) //création d'une matrice vide de taille n*n
	for i := 0; i < n; i++ {
		coeff[i] = make([]float64, n)
	}
	var somme float64
	et := 0.75 //écart-type

	for x := -n / 2; x <= n/2; x++ { //calcul de la matrice de convolution de gauss, c'est débile ça recalcule à chaque fois
		//mais j'arrive pas un faire une fonction modulable qui renvoie une matrice n*n
		for y := -n / 2; y <= n/2; y++ {
			coeff[x+n/2][y+n/2] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			somme += coeff[x+n/2][y+n/2]
		}
	}

	return coeff, somme
}

func convolutionGauss(x int, y int, img image.Image, n int, coeff *[][]float64, somme float64) color.RGBA { //convolution avec la fonction de gauss avec une matrice de taille n
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0

	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			rouge, vert, bleu, _ := (img).At(x+i, y+j).RGBA()
			r += (*coeff)[i+n/2][j+n/2] * float64(rouge)
			v += (*coeff)[i+n/2][j+n/2] * float64(vert)
			b += (*coeff)[i+n/2][j+n/2] * float64(bleu)

		}
	}
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff} //retour de la couleur du pixel normalisé x255 car la fct rgba renvoie des uint32 multipliés par alpha qui vaut 255 ici
	return pix
}

func flouGauss(img *image.Image, res *image.Image, n int, wg *sync.WaitGroup, minX int, minY int, maxX int, maxY int) { //flou de gauss avec matrice de taille n (plus n est grand plus l'effet est important), peut-etre faire plusieurs itérations
	imgFlou := image.NewRGBA((*img).Bounds())
	coeff, somme := calculGauss(n)

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			imgFlou.Set(x, y, convolutionGauss(x, y, *img, n, &coeff, somme))
		}
	}

	*res = imgFlou
	wg.Done()
}

func flouUniforme(img image.Image) image.Image { //applique un filtre uniforme 3x3 à chaque pixel
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

func contours(img image.Image, puissance int) image.Image { //filtre laplacien : fonctionne mieux, utiliser un autre filtre? sobel, prewitt...
	imgCont := image.NewRGBA(img.Bounds())
	coeff := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //laplacien
	//coeff := [3][3]float64{{-2,-2,0},{-2,0,2},{0,2,2}} //sobel
	somme := float64(puissance) //pose problème entre nv détails , et efficacité ==> eventuellement le proposer en réglage à l'utilisateur
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont.Set(x, y, convolution(x, y, img, coeff, somme))
		}
	}

	return negatifNB(imgCont) //applique le filtre négatif pour que ça soit "plus joli"
}

func contoursPrewitt(img image.Image, puissance int) image.Image { //filtre prewitt
	imgCont0 := image.NewRGBA(img.Bounds())
	imgCont90 := image.NewRGBA(img.Bounds())
	imgRes := image.NewRGBA(img.Bounds())

	coeff1 := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //Prewitt 0°
	coeff2 := [3][3]float64{{-2, -2, 0}, {-2, 0, 2}, {0, 2, 2}}      //Prewitt 90°

	somme := float64(puissance) //pose problème entre nv détails , et efficacité ==> eventuellement le proposer en réglage à l'utilisateur
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont0.Set(x, y, convolution(x, y, img, coeff1, somme))
		}
	}

	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			imgCont90.Set(x, y, convolution(x, y, img, coeff2, somme))
		}
	}

	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			pix := uint8(math.Sqrt(math.Pow(float64(imgCont0.RGBAAt(x, y).R), 2) + math.Pow(float64(imgCont90.RGBAAt(x, y).R), 2)))
			imgRes.Set(x, y, color.Gray{Y: pix})
		}
	}

	return negatifNB(imgRes) //applique le filtre négatif pour que ça soit "plus joli"
}

func nettete(img image.Image) image.Image { //pb à voir fait de la merde, n'est pas référencé dans le dispatch
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

func separation(img image.Image) (image.Image, image.Image, image.Image) { //sert à séparer les 3 composantes de l'image, pas très utile
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

func despeckleNB(img image.Image, x int, y int, n int, coeffGauss [][]float64, sommeGauss float64) color.Gray { //réduction de bruit Noir et Blanc : fonctionnel
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

		t = uint32(convolutionGauss(x, y, img, 3, &coeffGauss, sommeGauss).R) //on prend la composante rouge ici mais peu importe vu qu'on est en niveaux de gris
	}

	pix := color.Gray{Y: uint8(t)} //on reconstruit le pixel modifié
	return pix

}

//MARCHE PAS en couleurs
func despeckleRVB(img image.Image, x int, y int, n int, coeffGauss [][]float64, sommeGauss float64) color.RGBA { //déparasitage de chaque composante, n taille matrice
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
		r = uint32(convolutionGauss(x, y, img, 3, &coeffGauss, sommeGauss).R)
	}
	if float64(v) <= moyenneV-stdevV || float64(v) >= moyenneV+stdevV {
		v = uint32(convolutionGauss(x, y, img, 3, &coeffGauss, sommeGauss).G)
	}
	if float64(b) <= moyenneB-stdevB || float64(b) >= moyenneR+stdevB {
		b = uint32(convolutionGauss(x, y, img, 3, &coeffGauss, sommeGauss).B)
	}

	pix := color.RGBA{R: uint8(r), G: uint8(v), B: uint8(b), A: 0xff} //on reconstruit le pixel modifié
	return pix
}

func debruitageNB(img image.Image, nbIterations int, n int) image.Image {
	imgDebruit := image.NewGray(img.Bounds())
	coeffGauss, sommeGauss := calculGauss(3)

	for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
		for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
			imgDebruit.Set(x, y, despeckleNB(img, x, y, n, coeffGauss, sommeGauss))
		}
	}

	for k := 0; k < nbIterations-1; k++ {
		for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
			for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
				imgDebruit.Set(x, y, despeckleNB(imgDebruit, x, y, n, coeffGauss, sommeGauss))
			}
		}
	}

	return imgDebruit
}

func debruitageRVB(img image.Image, nbIterations int, n int) image.Image { //ne marche pas bien, n'est pas référencé dans le dispatch
	imgDebruit := image.NewRGBA(img.Bounds())
	coeffGauss, sommeGauss := calculGauss(3)

	for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
		for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
			imgDebruit.Set(x, y, despeckleRVB(img, x, y, n, coeffGauss, sommeGauss))
		}
	}

	for k := 0; k < nbIterations-1; k++ {
		for y := img.Bounds().Min.Y + n/2; y < img.Bounds().Max.Y-n/2; y++ {
			for x := img.Bounds().Min.X + n/2; x < img.Bounds().Max.X-n/2; x++ {
				imgDebruit.Set(x, y, despeckleRVB(imgDebruit, x, y, n, coeffGauss, sommeGauss))
			}
		}
	}

	return imgDebruit
}

func dispatch(img *image.Image, n int, param1 int, param2 int) image.Image { //permet de sélectionner quelle transfo en fonction de l'entrée utilisateur du programme principal
	//variable param pour les traitements nécessitant un niveau de puissance selectionné par l'utilisateur
	var img1, img2, img3, img4 image.Image
	imgRes := image.NewRGBA((*img).Bounds())
	var wg sync.WaitGroup
	imgp1, imgp2, imgp3, imgp4 := &img1, &img2, &img3, &img4

	maxX := (*img).Bounds().Max.X
	maxY := (*img).Bounds().Max.Y

	wg.Add(4)

	switch n {

	case 1:
		//negatifNB(img)
		break

	case 2:

		go negatifRVB(img, imgp1, &wg, 0, 0, maxX/2, maxY/2)
		go negatifRVB(img, imgp2, &wg, maxX/2, 0, maxX, maxY/2)
		go negatifRVB(img, imgp3, &wg, 0, maxY/2, maxX/2, maxY)
		go negatifRVB(img, imgp4, &wg, maxX/2, maxY/2, maxX, maxY)

		break

	case 3:
		//niveauGris(img)
		break

	case 4:
		//flouUniforme(img)
		break

	case 5:
		go flouGauss(img, imgp1, 7, &wg, 2, 2, maxX/2, maxY/2) //param1=taille matrice /!\ nb impairs= puissance du flou en fct de la résolution et de l'envie de l'utilisateur
		go flouGauss(img, imgp2, 7, &wg, maxX/2, 2, maxX-2, maxY/2)
		go flouGauss(img, imgp3, 7, &wg, 2, maxY/2, maxX/2, maxY-2)
		go flouGauss(img, imgp4, 7, &wg, maxX/2, maxY/2, maxX-2, maxY-2)
		break

	case 6:
		//debruitageNB(img, param1, param2) //param1 = nbIteration du filtre 1 à 2 conseillé
		//param2 = taille matrice 5 conseillé (nombres impairs /!\)
		break

	case 7:
		//contours(img, param1) //param1 = puissance de séparation en paramètre 8 voire 16, plus c'est élevé plus seuls les "gros" contours seront visibles
		break

	case 8:
		//contoursPrewitt(img, param1) //pareil que contours mais utilise le filtre de Prewitt à la place
		break

	}
	wg.Wait()

	draw.Draw(imgRes, image.Rect(0, 0, maxX/2, maxY/2), img1, image.Point{}, draw.Src)
	draw.Draw(imgRes, image.Rect(maxX/2, 0, maxX, maxY/2), img2, image.Point{X: maxX / 2}, draw.Src)
	draw.Draw(imgRes, image.Rect(0, maxY/2, maxX/2, maxY), img3, image.Point{Y: maxY / 2}, draw.Src)
	draw.Draw(imgRes, image.Rect(maxX/2, maxY/2, maxX, maxY), img4, image.Point{X: maxX / 2, Y: maxY / 2}, draw.Src)

	return imgRes
}

func main() {
	fmt.Println(time.Now())
	imageTest := importImage()
	ecritureFichier(dispatch(&imageTest, 5, 0, 0))
	fmt.Println(time.Now())
}
