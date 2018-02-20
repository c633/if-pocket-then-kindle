package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/motemen/go-pocket/api"
)

func mkdir(item api.Item) string {
	url := item.GivenURL
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	dir, err := ioutil.TempDir("", "if-pocket-then-kindle")
	if err != nil {
		log.Fatal(err)
	}
	tmpfn := filepath.Join(dir, fileName)
	return tmpfn
}

func download(fp string, item api.Item) {
	url := item.GivenURL
	log.Println("Downloading", url, "to", fp)

	output, err := os.Create(fp)
	if err != nil {
		fmt.Println("Error while creating", fp, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
	}
}

func convert(fp string, item api.Item) string {
	log.Printf("Running k2pdfopt...")
	cmd := exec.Command(k2pdfopt, "-ui- -dev kp3 -x", fp)
	if err := cmd.Run(); err != nil {
		log.Println("Error while converting", "-", err)
	}
	converted := strings.TrimSuffix(fp, filepath.Ext(fp))
	return converted + "_k2opt" + filepath.Ext(fp)
}
