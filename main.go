package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* use clipboard content if github url */

	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
		if !isGithubUrl(url) {
			fmt.Println("provided url is not a github url")
			os.Exit(1)
		}
	} else {
		cb, err := clipboard.ReadAll()
		if err != nil {
			panic("cannot read from clipboard")
		}

		if isGithubUrl(cb) {
			url = cb
		} else {
			fmt.Println("expected argument github link not found")
			os.Exit(1)
		}
	}

	filename := path.Base(url)
	if fileExists(filename) {
		fmt.Printf("\nfile %s already exists", filename)
		fmt.Print("\n\nPress the Enter to override")
		fmt.Scanln() // wait for Enter Key
		delPrevLine()
	}

	fileContent := fetch(url)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(fileContent)
	fmt.Printf("\nCreated file %s \n", filename)
}

func fetch(url string) string {
	if strings.HasPrefix(url, "https://github.com") {
		url = strings.Replace(url, "https://github.com", "https://raw.githubusercontent.com", 1)
		url = strings.Replace(url, "/blob", "", 1)
	}

	// log.Debug().Msgf("fetching url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func isGithubUrl(url string) bool {
	return strings.HasPrefix(url, "https://raw.githubusercontent.com") || strings.HasPrefix(url, "https://github.com")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func delPrevLine() {
	fmt.Print("\033[A\033[A\033[A\033[0J\033[A")
}
