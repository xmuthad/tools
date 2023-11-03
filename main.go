package main

import (
	"flag"
	"fmt"
	"tools/pkg/download/http"
)

func main() {
	directoryURL := flag.String("url", "", "The URL of the directory to download")

	flag.Parse()

	if *directoryURL == "" {
		fmt.Println("Please provide a URL using the -url flag")
		return
	}

	err := http.DownloadDirectory(*directoryURL, "")
	if err != nil {
		panic(err)
	}
}
