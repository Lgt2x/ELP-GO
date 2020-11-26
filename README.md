# ELP-GO

A Golang project of a TCP server for image processing

Shout-outs to PFR

## Authors
Louis Gombert, Constantin Thebaudeau, Safae Hariri, and Antoine Merle

## Usage
```
go run src/server/server.go

# In another terminal
go run src/client/client.go

# Or
go run src/client/clientcli.go <port> <filter id> <source> <dest>
```
If package `elputils` isn't found, add this project's directory to `$GOPATH` env variable
If standard packages can't be found, check your `$GOROOT` en var.

## List of available filters
1) Negative - Black and White
2) Negative - RGB
3) Grayscale
4) Uniform blur
5) Gauss blur
6) Boundaries detection (laplacian)
7) Boundaries detection (Prewitt)
8) Noise reduction - Black and White


## Client-Server communication
1) Establishing connection
2) Server sends the list of the available filters (using a String and ‘;’ separations)
3) Client sends the filter_id (from 1 to 8)
4) Server checks if the id is correct and return 0 or 1
5) If 0, client can re-ask a good filter ID to the user or terminate the connexion
6) Then, the client sends the image (name (64bits), size (10bits), and the content with a buffer system) [before the client checks if the image is a JPEG-type]
7) Server sends back the modified image using the same way


## Resources

 * https://www.devdungeon.com/content/working-images-go
 * https://yourbasic.org/golang/create-image/
 * https://golang.org/pkg/image/#pkg-examples