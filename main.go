package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {

	ARGS := os.Args[1:]

	var fileName string

	if len(ARGS) >= 2 {
		fileName = ARGS[1] + ".png"
	} else {
		fileName = "no_name.png"
	}

	videoCode, err := getVideoCode(ARGS[0])

	URL := "https://img.youtube.com/vi/" + videoCode + "/maxresdefault.jpg"

	if err != nil {
		log.Fatal(err)
	}

	err = downloadFile(URL, fileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File %s downlaod in current working directory", fileName)
}

// it clears any Youtube Url and return only the Video Id Code
func getVideoCode(url string) (string, error) {

	regex := regexp.MustCompile(`watch\?v=[a-zA-Z-0-9]*`)
	match := regex.FindStringSubmatch(url)

	if len(match) > 0 {
		result := match[0]
		code := strings.Replace(result, "watch?v=", "", 1)
		return code, nil
	}

	regex = regexp.MustCompile(".be/[a-zA-Z-0-9]*")
	match = regex.FindStringSubmatch(url)

	if len(match) > 0 {
		result := match[0]
		code := strings.Replace(result, ".be/", "", 1)
		return code, nil
	}

	err := errors.New("Error trying to get the Video Code, verify the Youtube Url")

	return "", err
}

func downloadFile(URL string, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return errors.New("Video not found")
	}

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
