package main

import (
	"bufio"
	"bytes"
	"fmt"
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
					var readCmd *exec.Cmd
					if strings.HasSuffix(f.Name(), ".doc") {
						readCmd = exec.Command("antiword", filepath.Join(documentPath, d.Name(), f.Name()))
					} else if strings.HasSuffix(f.Name(), ".docx") {
						f, err := os.Open(filepath.Join(documentPath, d.Name(), f.Name()))
						if err != nil {
							log.Fatal(f)
						}

						readCmd = exec.Command("docx2txt", "-", "-")
						readCmd.Stdin = f
					} else {
						log.Fatal(fmt.Errorf("unkown file format: '%s'", f.Name()))
					}

					var b bytes.Buffer
					readCmd.Stdout = &b
					if err := readCmd.Run(); err != nil {
						log.Fatal(err)
					}

					rf, err := os.Create(filename)
					if err != nil {
						log.Fatal(err)
					}
					defer rf.Close()

					scanner := bufio.NewScanner(&b)
					for scanner.Scan() {
						rf.WriteString(strings.TrimSpace(scanner.Text()) + "\n")
					}

					codeCmd := exec.Command("code", filename)
					if err := codeCmd.Start(); err != nil {
						log.Fatal(err)
					}

					return
				}
			}
		}
	}
}
