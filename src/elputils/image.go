// Image processing functions, to provide transformation functions

package elputils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"sync"
)

type imageStrip struct {
	image image.Image
	rect  image.Rectangle
}

//list of all available filters
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

// Apply a filter asynchronously to a given image, using 4 goroutines
func ApplyFilterAsync(sourceImg *image.RGBA, filter int, routines int) image.Image {

	width := (*sourceImg).Bounds().Max.X
	height := (*sourceImg).Bounds().Max.Y

	// Temp images to store results
	tmpImages := make([]image.Image, routines)

	// Divide the image in strips to speed up computation
	imagePart := make([]imageStrip, routines)
	for i := range imagePart {
		imagePart[i].rect = image.Rect(i*width/routines, 0, (i+1)*width/routines, height)
	}

	// Compute transformations in different threads
	var wg sync.WaitGroup
	wg.Add(len(imagePart))

	for i := range imagePart {
		go Dispatch(sourceImg, &tmpImages[i], filter, imagePart[i].rect, &wg)
	}

	wg.Wait()

	resultImg := image.NewRGBA((*sourceImg).Bounds())
	for i := range imagePart {
		draw.Draw(resultImg, imagePart[i].rect, tmpImages[i], imagePart[i].rect.Min, draw.Src)
	}

	return resultImg
}

// Write an image object to a file
func ImageToFile(image image.Image, destination string) {
	// Create blank file
	output, err := os.Create(destination)
	if err != nil {
		fmt.Println("Couldn't create file", destination)
	}

	err = jpeg.Encode(output, image, nil)
	if err != nil {
		fmt.Println("Couldn't write image to file", destination)
	}

	_ = output.Close()
}

// Creates an image object from a file
func FileToImage(path string) *image.RGBA {
	input, err := os.Open(path)
	if err != nil {
		fmt.Println("Couldn't open file", path)
	}

	img, err := jpeg.Decode(input)
	if err != nil {
		fmt.Println("Couldn't import image")
	}

	if img != nil {
		imgRes := image.NewRGBA(img.Bounds())

		for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x <= img.Bounds().Max.X; x++ {
				imgRes.Set(x, y, img.At(x, y))
			}
		}
		return imgRes
	}

	return nil
}

// Converts image to grayscale
func GreyScale(img *image.RGBA, res *image.Image, rect image.Rectangle) {
	imgGrey := image.NewGray(rect)

	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			imgGrey.Set(x, y, (*img).At(x, y))
		}
	}
	*res = imgGrey
}

// Converts image to black & white negative
func NegativeBW(img *image.RGBA, res *image.Image, rect image.Rectangle) {
	var imgGris image.Image
	GreyScale(img, &imgGris, rect) //first we need to convert the image in grayscale
	imgNeg := image.NewGray(rect)
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			z, _, _, _ := imgGris.At(x, y).RGBA()
			pix := color.Gray{Y: 255 - uint8(z)}
			imgNeg.Set(x, y, pix)
		}
	}
	*res = imgNeg
}

// Reverses each component of each pixel of the image considered
func NegativeRGB(img *image.RGBA, res *image.Image, rect image.Rectangle) {
	imgNeg := image.NewRGBA(rect)
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			r, v, b, _ := (*img).At(x, y).RGBA()
			pix := color.RGBA{R: 255 - uint8(r/256), G: 255 - uint8(v/256), B: 255 - uint8(b/256), A: 0xff}
			imgNeg.Set(x, y, pix)
		}
	}
	*res = imgNeg

}

// Applies the 3x3 convolution matrix provided on the pixel considered
// We make the sum of each neighbour pixel multiplied by the coefficient and apply the mean in the pixel considered for each component
func Convolution(x int, y int, img *image.RGBA, coefficient *[3][3]float64, somme float64) color.RGBA {
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			rouge, vert, bleu, _ := (*img).At(x+i, y+j).RGBA()
			r += coefficient[i+1][j+1] * float64(rouge)
			v += coefficient[i+1][j+1] * float64(vert)
			b += coefficient[i+1][j+1] * float64(bleu)

		}
	}
	//warning : RGBA method provide each component of the pixel multiplied by the alpha value : that's why we need to divide it by 256
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff}
	return pix
}

//Computes gauss coefficients and make a NxN gauss matrix
func GaussMatrix(n int) ([][]float64, float64) {
	coeff := make([][]float64, n) //creates empty matrix
	for i := 0; i < n; i++ {
		coeff[i] = make([]float64, n)
	}
	var somme float64
	et := 0.75 //standard deviation used there

	for x := -n / 2; x <= n/2; x++ {
		for y := -n / 2; y <= n/2; y++ {
			coeff[x+n/2][y+n/2] = 100 * math.Exp(-(math.Pow(float64(x), 2)+math.Pow(float64(y), 2))/2*math.Pow(et, 2)) / (2 * math.Pi * math.Pow(et, 2))
			somme += coeff[x+n/2][y+n/2]
		}
	}

	return coeff, somme
}

//Same as Convolution but for a NxN matrix
func ConvolutionGauss(x int, y int, img *image.RGBA, n int, coeff *[][]float64, somme float64) color.RGBA {
	var pix color.RGBA
	var r, v, b float64
	r = 0
	v = 0
	b = 0

	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			rouge, vert, bleu, _ := (*img).At(x+i, y+j).RGBA()
			r += (*coeff)[i+n/2][j+n/2] * float64(rouge)
			v += (*coeff)[i+n/2][j+n/2] * float64(vert)
			b += (*coeff)[i+n/2][j+n/2] * float64(bleu)

		}
	}
	//warning : RGBA method provide each component of the pixel multiplied by the alpha value : that's why we need to divide it by 256
	pix = color.RGBA{R: uint8(r / (256 * somme)), G: uint8(v / (256 * somme)), B: uint8(b / (256 * somme)), A: 0xff}
	return pix
}

//applies a gaussian blur using a NxN convolution matrix (more N is higher, more the blurry effect is important)
//If the image is very big, it can be necessary to apply this filter several times
func GaussBlur(img *image.RGBA, res *image.Image, n int, rect image.Rectangle) {
	imgFlou := image.NewRGBA(rect)
	coeff, somme := GaussMatrix(n)

	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			imgFlou.Set(x, y, ConvolutionGauss(x, y, img, n, &coeff, somme))
		}
	}

	*res = imgFlou

}

//applies an uniform blur on each pixel using a 3x3 convolution matrix
func UniformBlur(img *image.RGBA, res *image.Image, rect image.Rectangle) {
	imgFlou := image.NewRGBA((*img).Bounds())
	coeff := [3][3]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	somme := 9.0
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			imgFlou.Set(x, y, Convolution(x, y, img, &coeff, somme))
		}
	}

	*res = imgFlou

}

//Boundaries detection using a Laplacian filter
//we apply a 3x3 convolution matrix on the source image which determines changes of gradient intensity
func Boundaries(img *image.RGBA, res *image.Image, puissance int, rect image.Rectangle) {
	imgCont := image.NewRGBA(rect)
	coeff := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //laplacien
	somme := float64(puissance)                                     //This value influences the power of the filter (lower it is, better is the boundaries detection but can create a lot of noise)
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			imgCont.Set(x, y, Convolution(x, y, img, &coeff, somme))
		}
	}
	NegativeBW(imgCont, res, rect) //applies a negative filter in order to be "prettier"
}

//Boundaries detection using the Prewitt filter
//we apply 2 3x3 convolution matrices in two different directions (0° and 90°) on the source image and then we combine these two new images
func PrewittBorders(img *image.RGBA, res *image.Image, puissance int, rect image.Rectangle) {
	imgCont0 := image.NewRGBA(rect)
	imgCont90 := image.NewRGBA(rect)
	imgRes := image.NewRGBA(rect)

	coeff1 := [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}} //Prewitt 0°
	coeff2 := [3][3]float64{{-2, -2, 0}, {-2, 0, 2}, {0, 2, 2}}      //Prewitt 90°

	somme := float64(puissance) //This value influences the power of the filter (lower it is, better is the boundaries detection but can create a lot of noise)

	//Applies the first filter 0°
	for y := img.Bounds().Min.Y + 1; y <= img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x <= img.Bounds().Max.X-1; x++ {
			imgCont0.Set(x, y, Convolution(x, y, img, &coeff1, somme))
		}
	}

	//the second one
	for y := img.Bounds().Min.Y + 1; y <= img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x <= img.Bounds().Max.X-1; x++ {
			imgCont90.Set(x, y, Convolution(x, y, img, &coeff2, somme))
		}
	}

	//combine both filters in one image which is returned
	for y := img.Bounds().Min.Y + 1; y < img.Bounds().Max.Y-1; y++ {
		for x := img.Bounds().Min.X + 1; x < img.Bounds().Max.X-1; x++ {
			pix := uint8(math.Sqrt(math.Pow(float64(imgCont0.RGBAAt(x, y).R), 2) + math.Pow(float64(imgCont90.RGBAAt(x, y).R), 2)))
			imgRes.Set(x, y, color.Gray{Y: pix})
		}
	}

	NegativeBW(imgRes, res, rect) //applies a negative filter in order to be "prettier"

}

//computes the mean and standard deviation of each pixel in the considered pixel's area and then compare it with the pixel considered
//if the pixel is outside [mean-stdev; mean+stdev], we apply a 3x3 gaussian blur on the pixel considered
func DespeckleBW(img *image.RGBA, x int, y int, n int, coeffGauss *[][]float64, sommeGauss float64) color.Gray {
	var moyenne, stdev float64
	moyenne, stdev = 0, 0

	//mean
	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on évite le pixel central
				t, _, _, _ := img.At(x+i, y+j).RGBA()
				moyenne += float64(t)

			}
		}
	}

	moyenne = moyenne / (math.Pow(float64(n), 2) - 1)

	//stdev
	for i := -n / 2; i <= n/2; i++ {
		for j := -n / 2; j <= n/2; j++ {
			if i != 0 || j != 0 { //on exclut le pixel central
				t, _, _, _ := img.At(x+i, y+j).RGBA()
				stdev += math.Pow(float64(t)-moyenne, 2)
			}
		}
	}

	stdev = math.Sqrt(stdev / (math.Pow(float64(n), 2) - 1))

	//tests if the pixel is in the interval [mean-stdev; mean+stdev] and applies 3x3 gaussian blur if necessary
	t, _, _, _ := img.At(x, y).RGBA()

	if float64(t) <= moyenne-stdev || float64(t) >= moyenne+stdev {

		t = uint32(ConvolutionGauss(x, y, img, 3, coeffGauss, sommeGauss).R) //We take here the red component but it doesn't care because we are in grayscale
	}

	pix := color.Gray{Y: uint8(t)} //builds the new pixel in grayscale
	return pix

}

//This function applies DespeckleBW on each pixel of the image 2 times in order to have great results
func NoiseReductionBW(img *image.RGBA, res *image.Image, nbIterations int, n int, rect image.Rectangle) {
	imgDebruit := image.NewRGBA(rect)
	coeffGauss, sommeGauss := GaussMatrix(3)

	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			imgDebruit.Set(x, y, DespeckleBW(img, x, y, n, &coeffGauss, sommeGauss))
		}
	}

	for k := 0; k < nbIterations-1; k++ {
		for y := imgDebruit.Bounds().Min.Y + n/2; y <= imgDebruit.Bounds().Max.Y-n/2; y++ {
			for x := imgDebruit.Bounds().Min.X + n/2; x <= imgDebruit.Bounds().Max.X-n/2; x++ {
				imgDebruit.Set(x, y, DespeckleBW(imgDebruit, x, y, n, &coeffGauss, sommeGauss))
			}
		}
	}

	*res = imgDebruit

}

// Apply a filter to a source image following a filter id
// Filter id <=> filter matching follows the order defined in filterList
func Dispatch(source *image.RGBA, dest *image.Image, filter int, rect image.Rectangle, wg *sync.WaitGroup) {
	switch filter {
	case 1:
		NegativeBW(source, dest, rect)
		break
	case 2:
		NegativeRGB(source, dest, rect)
		break
	case 3:
		GreyScale(source, dest, rect)
		break
	case 4:
		UniformBlur(source, dest, rect)
		break
	case 5:
		GaussBlur(source, dest, 7, rect)
		break
	case 6:
		NoiseReductionBW(source, dest, 2, 5, rect)
		break
	case 7:
		Boundaries(source, dest, 8, rect)
		break
	case 8:
		PrewittBorders(source, dest, 32, rect) //à vérifier niveau de puissance selon image
		break
	}

	wg.Done()
}
