package main

import (
	"flag"
	"tools/pkg/download/http"
	"tools/pkg/upload"
)

func main() {
	directoryURL := flag.String("url", "", "The URL of the directory to download")
	uploadJars := flag.Bool("upload-jars", false, "Implement the functionality to upload JAR files")
	localPath := flag.String("local-path", "~/.m2/repository", "Local file path")
	remoteMaven := flag.String("maven-url", "", "Maven repository URL")
	repoID := flag.String("repo-id", "", "Maven repository ID")

	flag.Parse()

	if *directoryURL != "" {
		err := http.DownloadDirectory(*directoryURL, "")
		if err != nil {
			panic(err)
		}
	}

	if *uploadJars {
		err := upload.UploadIt(upload.MvnUploadJars{}, *localPath, *remoteMaven, *repoID)
		if err != nil {
			panic(err)
		}
	}

}
