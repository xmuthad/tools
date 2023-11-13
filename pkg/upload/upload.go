package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MavenUploader interface {
	Upload(jarPath, pomPath, repoURL, repoID string) error
}

func UploadIt(uploader MavenUploader, localPath, repoURL, repoID string) (err error) {
	err = filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jar" {
			pomPath := strings.TrimSuffix(path, ".jar") + ".pom"
			if _, err := os.Stat(pomPath); os.IsNotExist(err) {
				fmt.Printf("POM file does not exist for %s, skipping\n", path)
				return nil
			}
			fmt.Printf("Found JAR and POM: %s\n", path)
			if err := uploader.Upload(path, pomPath, repoURL, repoID); err != nil {
				fmt.Println(err)
				// Continue on error
			}
		}
		return nil
	})
	return
}
