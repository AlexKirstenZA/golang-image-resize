package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	urls := []string{
		"https://www.placecage.com/800/1300",
		"https://www.fillmurray.com/1000/1500",
		"https://www.stevensegallery.com/1500/1000",
	}

	for i, url := range urls {
		filename := "image_" + strconv.Itoa(i) + ".jpeg"
		filepath, err := downloadImage(url, filename)
		if err != nil {
			fmt.Printf("Error processing %v: %v", url, err)
		}

		fmt.Println(filepath)
	}
}

func downloadImage(url string, fileName string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
