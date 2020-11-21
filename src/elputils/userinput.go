// Utility functions for user input processing
// used for both client and server
package elputils

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"net"
	"os"
	"strconv"
	"strings"
)

// Input a string from stdin and handle errors
func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return input
}

// Input filter, client-side
// Tries again until it succeeds
func InputFilter(conn net.Conn, filterList []string) int {
	var num int

	// Display available choices
	fmt.Println("Available filters :")
	for i := 0; i < len(filterList); i++ {
		fmt.Printf("%d. : %s\n", i+1, filterList[i])
	}
	fmt.Print("Enter a valid filter id ")
	filter := InputString()

	// Send and validate filter id
	SendString(conn, filter)
	validationServer := ReceiveString(conn, '\n')
	validFilter := strings.Compare(strings.Trim(validationServer, "\n"), "1")

	if validFilter != 0 {
		fmt.Println("Invalid choice")
		// Try again if invalid
		num = InputFilter(conn, filterList)
	} else{
		num, _ = strconv.Atoi(strings.Trim(filter, "\n"))
	}

	return num
}

// Receive and process filter choice, server-side
func ReceiveFilter(conn net.Conn, maxFil int) int {
	filterValidated := false
	for filterValidated != true {
		filterStr := strings.Trim(ReceiveString(conn, '\n'), "\n")
		filter, err := strconv.Atoi(filterStr)
		if err != nil {
			fmt.Println("Error - Received filter can't be converted into an int")
			panic(err)
		}

		if filter > 0 && filter <= maxFil {
			fmt.Printf("Received filter : %s\n", filterStr)
			SendString(conn, "1\n")
			filterValidated = true
			return filter
		} else {
			fmt.Printf("Invalid filter : %s\n", filterStr)
			SendString(conn, "0\n")
		}
	}

	return 0
}

func InputImageVerification(imagePath string) bool {
	// Filter .jpg images only
	if !strings.HasSuffix(imagePath, ".jpg") && !strings.HasSuffix(imagePath, ".jpeg") {
		fmt.Println("Unsupported image format. Use a jpeg image")
		return false
	}

	// check if the file can be decoded as jpeg or not
	input, err := os.Open(imagePath)
	_, err = jpeg.DecodeConfig(input)
	if err != nil {
		fmt.Println("Image has format .jpeg or .jpg but cannot be interpreted as such")
		return false
	}
	input.Close()

	// check if the file exists or not
	// NB: this is already done by the previous lines (os.Open then jpeg.DecodeConfig)
	if !FileExists(imagePath) {
		return false
	}
	return true
}

// Inputs an image path
func InputImagePath() string {
	fmt.Print("Relative path to the image ")
	imagePath := strings.Trim(InputString(), "\n")

	// check the validity of the image path given
	if !InputImageVerification(imagePath) {
		imagePath = InputImagePath()
	}
	return imagePath
}

// Check if a file exists in the current directory
func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func FillString(returnString string, toLength int) string {
	for {
		stringLength := len(returnString)
		if stringLength < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}

// Delete the file filename
func DeleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Suppression impossible")
	}
}
