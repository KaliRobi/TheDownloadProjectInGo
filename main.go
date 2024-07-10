package main

import (
	"download-project/download"
	"download-project/server"
)

func main() {

	var urlsSlice = []string{

		"https://raw.githubusercontent.com/GoogleContainerTools/distroless/main/java/README.md",
		"https://raw.githubusercontent.com/golang/go/master/README.md",
		"http://thisisanabsolutellyinvalidurl.org.e",
		"https://pkg.go.dev/github.com/posener/goreadme",
		"https://pkg.go.dev/go.jpap.org/godoc-readme-gen",
		"https://github.com/golang/example/blob/master/README.md",
		// "http://thisisanabsolutellyinvalidurl.org.e",
	}

	content, err := download.ReturnContentOrFail(urlsSlice)
	server.StartServer(content, err)
}
