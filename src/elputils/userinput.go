// Utility functions for user input processing
// used for both client and server
package elputils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Input a string from stdin and handle errors
func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	inputstr, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return inputstr
}

// Input filter, client-side
// Tries again until it succeeds
func InputFilter(conn net.Conn, filterList []string) int {
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
	filtre_valide := strings.Compare(validationServer[0:1], "1")

	if filtre_valide != 0 {
		fmt.Println("Invalid choice")
		// Try again if invalid
		InputFilter(conn, filterList)
	}
	num, _ := strconv.Atoi(strings.Trim(filter, "\n"))
	return num
}

// Receive and process filter choice, server-side
func ReceiveFilter(conn net.Conn, maxFil int) int {
	filterValidated := false
	for filterValidated != true {
		filterStr := strings.Trim(ReceiveString(conn, '\n'), "\n")
		filter, err := strconv.Atoi(filterStr)
		if err != nil {
			panic(err)
		}

		if filter > 0 && filter <= maxFil {
			fmt.Printf("Received filter : %s\n", filterStr)
			SendString(conn, "1\n")
			filterValidated = true
			return filter
		} else {
			fmt.Printf("Invalid filter : %s\n", filterStr)
			SendString(conn, "1\n")
		}
	}

	return 0
}

// Inputs an image path
func InputImagePath() string {
	fmt.Print("Relative path to the image ")
	imagePath := strings.Trim(InputString(), "\n")
	imagePathAbs, _ := filepath.Abs(imagePath)
	//fmt.Println(imagepath)
	//filter .jpg image only
	if !strings.HasSuffix(imagePath,".jpg")&&!strings.HasSuffix(imagePath, ".jpeg"){
		fmt.Println("Unsupported format")
		return InputImagePath()
	}
	if FileExists(imagePathAbs) {
		return imagePath
	}
	// To it again if it fails
	fmt.Println("File not found")
	return InputImagePath()
}

// Check if a file exists in the current directory
func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func FillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func NewName(filename string) string {
	// find the extension, store it then delete it
	indexExt := strings.Index(filename, ".")
	return filename[:indexExt] + "_modified" + filename[indexExt:]
}

// Delete the file filename
func DeleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Suppression impossible")
	}
}
