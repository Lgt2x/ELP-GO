package main

import (
	"ELP-GO/src/elputils"
	"fmt"
	"net"
	"os"
)

func usage() {
	fmt.Println("Usage : clientcli <port> <filter id> <source> <dest>")
}
func main() {
	// Get arguments from argc : port number, filter id, source image, destination

	var port string
	var filterId string
	var sourceImg string
	var destImg string

	if len(os.Args) == 5 {
		port = ":" + os.Args[1]
		filterId = os.Args[2]
		sourceImg = os.Args[3]
		destImg = os.Args[4]
	} else {
		elputils.PrintRedLn("Wrong number of arguments")
		usage()
		os.Exit(1)
	}

	// Connecting to the server on port
	fmt.Printf("Connecting to server on port %s\n", port)
	conn, err := net.Dial("tcp", "localhost"+port)

	if err != nil {
		elputils.PrintRedLn("Couldn't listen on port" + port + ". Is the server running ?")
		return
	}

	// Connection successful, close the connection when it's over
	fmt.Printf("Connection established with server on port %s\n\n", port)
	defer conn.Close()

	// Send filter
	_ = elputils.ReceiveArray(conn, ";", '\n') // Ignore filter list sent
	elputils.SendString(conn, filterId+"\n")
	filterResult := elputils.ReceiveString(conn, '\n')

	// Get an error if the filter can't be found
	if filterResult == "0\n" {
		elputils.PrintRedLn("Server couldn't apply requested filter. Is the provided id valid ?")
		os.Exit(1)
	}
	fmt.Printf("Selected filter '%s'\n", filterId)

	// Send file
	fmt.Println("Sending image", sourceImg)
	elputils.UploadFile(conn, sourceImg)

	// Receiving and storing the modified image
	fmt.Println("\nWaiting for the modified image...")
	elputils.ReceiveFile(conn, destImg)
	elputils.PrintGreenLn("Transformation complete, output stored in " + destImg)
}
