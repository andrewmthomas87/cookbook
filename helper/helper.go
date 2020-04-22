package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	documentPath, dataPath := os.Args[1], os.Args[2]

	directories, err := ioutil.ReadDir(documentPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range directories {
		if d.IsDir() {
			files, err := ioutil.ReadDir(filepath.Join(documentPath, d.Name()))
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				index := strings.LastIndex(f.Name(), ".")
				filename := filepath.Join(dataPath, d.Name(), f.Name()[:index]+".txt")
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					codeCmd := exec.Command("code", filename)
					if err := codeCmd.Start(); err != nil {
						log.Fatal(err)
					}

					openCmd := exec.Command("open", filepath.Join(documentPath, d.Name(), f.Name()))
					if err := openCmd.Start(); err != nil {
						log.Fatal(err)
					}

					return
				}
			}
		}
	}
}
