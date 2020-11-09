# ELP-GO

Golang project of a TCP server for image processing

Shoutouts to PFR

## Usage
```
go run src/server/server.go

# In another terminal
go run src/client/client.go
```
If package elputils isn't found, add this project's directory to $GOPATH env variable
If standard packages can't be found, check your $GOROOT en var.

Resources :
 * https://www.devdungeon.com/content/working-images-go
 * https://yourbasic.org/golang/create-image/
 * https://golang.org/pkg/image/#pkg-examples