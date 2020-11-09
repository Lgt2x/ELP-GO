package elputils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Input a string from stdin
func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	filtre, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return filtre
}

func InputFilter(conn net.Conn) {
	filter := InputString()
	SendString(conn, filter)

	validationServer := ReceiveString(conn, '\n')
	filtre_valide := strings.Compare(validationServer[0:1], "1")

	if filtre_valide != 0 {
		fmt.Println("Saisie invalide")
		InputFilter(conn)
	}
}

func InputImagePath() (string, string) {
	fmt.Print("Saisie du chemin de l'image_utils (relatif): ")
	imagePath := InputString()
	imagePathAbs, _ := filepath.Abs(imagePath[:len(imagePath)-1])

	if FileExists(imagePathAbs) {
		return imagePath[:len(imagePath)-1], imagePathAbs
	}
	return InputImagePath()
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func ValideFiltre(c net.Conn) {
	validationFiltre := false
	for validationFiltre != true {
		choixFiltre := ReceiveString(c, '\n')
		fmt.Println(choixFiltre)

		switch choixFiltre {
		// enumeration des choix valides
		case "1":
			// validation et go pour machin truc
			fmt.Println("Choix valide")
			validationFiltre = true
			SendString(c, "1\n")
		default:
			fmt.Println("Choix invalide")
			SendString(c, "0\n")
		}
	}
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

func DeleteFile(filename string) {
	// delete the file filename
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Suppression impossible")
	}
}
