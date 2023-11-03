package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	defaultDownloadDirectory = "downloadByURL"
)

func DownloadDirectory(directoryURL string, downloadDir string) error {
	if downloadDir == "" {
		downloadDir = defaultDownloadDirectory
	}
	resp, err := http.Get(directoryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.HasSuffix(href, "/") {
			// This is a directory, we could recursively call DownloadDirectory here
			dir := downloadDir + "/" + href
			DownloadDirectory(directoryURL+href, dir)
		} else if href != "" {
			// This is a file, download it
			fmt.Println("Downloading file:", directoryURL)
			err := DownloadFile(href, downloadDir)
			if err != nil {
				fmt.Println("Error downloading file:", err)
			}
		}
	})

	return nil
}

func DownloadFile(fileURL, dirPath string) error {
	u, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	fmt.Println("current directory: ", dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0755)
	}
	filePath := path.Join(dirPath, path.Base(fileURL))
	fmt.Println("file path: ", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
