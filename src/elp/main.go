package main

import (
	"fmt"
	"image"
	"os"
)

func main() {
	// Read image from file that already exists
	existingImageFile, err := os.Open("poke.png")
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	_, imageType, err := image.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	//fmt.Println(imageData)
	fmt.Println(imageType)
}
