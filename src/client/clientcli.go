package main

import (
	"ELP-GO/src/elputils"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Println("Usage : clientcl <port> <filter id> <source> <dest>")
}
func main() {
	// Get arguments from argc : port number, filter id, source image, destination

	var port string
	var filter_id int
	var source_img string
	var dest_img string

	if len(os.Args) != 5 {
		port := ":" + os.Args[1]
		err, filter_id := strconv.Atoi(os.Args[2])
		source_img := os.Args[3]
		dest_img := os.Args[4]
	} else {
		fmt.Println("Wrong number of arguments")
		usage()
	}

	// Connecting to the server on port
	fmt.Printf("Connecting to server on port %s. You can change the port used by specifying it as an argc value eg 'client 5000'\n", port)
	conn, err := net.Dial("tcp", "localhost"+port)

	if err != nil {
		fmt.Printf("Couldn't listen on port %s. Is the server running ?\n", port)
		return
	}

	// Connection successful, close the connection when it's over
	fmt.Printf("Connection established with server on port %s\n\n", port)
	defer conn.Close()

	// Input filter
	filterList := elputils.ReceiveArray(conn, ";", '\n')
	filterNum := elputils.InputFilter(conn, filterList)
	fmt.Printf("Selected filter '%s'\n", filterList[filterNum-1])

	// Send file
	fmt.Println("Sending image", source_img)
	elputils.UploadFile(conn, source_img)

	// Receiving the modified image
	fmt.Println("\nWaiting for the modified image...")
	elputils.ReceiveFile(conn, dest_img)
	fmt.Println("Transformation complete, output stored in ", dest_img)
}
