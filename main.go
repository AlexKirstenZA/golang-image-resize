package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	//"github.com/disintegration/imaging"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type imageResult struct {
	filePath string
	err      error
}

type imageStyle struct {
	name string
	maxWidth int
	maxHeight int
}

func main() {
	err := cleanUpFiles()
	if err != nil {
		log.Fatal(err)
	}

	// Placeholder image urls using three different dimensions
	urls := []string{
		"https://www.placecage.com/800/1300",
		"https://www.fillmurray.com/1000/1500",
		"https://www.stevensegallery.com/1500/1000",
	}

	results := make(chan imageResult, len(urls))

	for i, url := range urls {
		filename := "image_" + strconv.Itoa(i) + ".jpeg"

		go func(url string, filename string) {
			fp, err := downloadImage(url, filename)
			results <- imageResult{filePath: fp, err: err}
		}(url, filename)
	}

	// Three thumbnails per original image
	var sizes []imageStyle
	sizes = append(sizes, imageStyle{name: "modal", maxHeight: 600, maxWidth: 600})
	sizes = append(sizes, imageStyle{name: "gallery", maxHeight: 300, maxWidth: 300})
	sizes = append(sizes, imageStyle{name: "avatar", maxHeight: 120, maxWidth: 120})

	resizeResult := make(chan imageResult, len(urls) * len(sizes))
	var count int

	for i, _ := range urls {
		result := <-results
		if result.err != nil {
			fmt.Println("Error occurred during download: ", result.err)
			continue
		}

		fmt.Println("Download successful:", result.filePath)
		count++

		// Process different image styles for downloaded image
		for _, style := range sizes {
			filename := "image_" + strconv.Itoa(i) + "_" + style.name + ".jpeg"

			go func(fileName string, style imageStyle) {
				fp, err := scaleImage(filename, style)
				resizeResult <- imageResult{filePath: fp, err: err}
			}(filename, style)
		}
	}

	for i := 0; i < count * len(sizes); i++ {
		result := <-resizeResult
		fmt.Println(result.filePath)
	}
}

func downloadImage(url string, fileName string) (string, error) {
	// Download image
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle non 200 response code
	if resp.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}

	// Create file on disk
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write byte stream from response body into file contents
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	// Get absolute filepath to return via channel
	fp, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return fp, nil
}

func scaleImage(fileName string, style imageStyle) (string, error) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)
	time.Sleep(time.Duration(n)*time.Second)

	return "Done: " + fileName, nil
}

// Removes old jpeg images from disk
func cleanUpFiles() error {
	files, err := filepath.Glob("*.jpeg")
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	return nil
}
