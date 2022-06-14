package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type downloadResult struct {
	filePath string
	err error
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

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))

	results := make(chan downloadResult, len(urls))

	for i, url := range urls {
		filename := "image_" + strconv.Itoa(i) + ".jpeg"
		go downloadImage(wg, url, filename, results)
	}

	for _ = range urls {
		result := <-results
		if result.err != nil {
			fmt.Println("Error occured during download: ", result.err)
			continue
		}

		fmt.Println("Download successful:", result.filePath)
	}

	wg.Wait()
}

func downloadImage(wg *sync.WaitGroup, url string, fileName string, ch chan downloadResult) {
	defer wg.Done()

	// Download image
	resp, err := http.Get(url)
	if err != nil {
		ch <- downloadResult{filePath: "", err: err}
		return
	}
	defer resp.Body.Close()

	// Handle non 200 response code
	if resp.StatusCode != 200 {
		ch <- downloadResult{filePath: "", err: errors.New("Received non 200 response code")}
		return
	}

	// Create file on disk
	file, err := os.Create(fileName)
	if err != nil {
		ch <- downloadResult{filePath: "", err: err}
		return
	}
	defer file.Close()

	// Write byte stream from response body into file contents
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		ch <- downloadResult{filePath: "", err: err}
		return
	}

	// Get absolute filepath to return via channel
	fp, err := filepath.Abs(fileName)
	if err != nil {
		ch <- downloadResult{filePath: "", err: err}
		return
	}

	ch <- downloadResult{filePath: fp, err: err}
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
