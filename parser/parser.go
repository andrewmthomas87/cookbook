package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	foldersCmd := flag.NewFlagSet("folders", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatal("expected 'file' or 'folders' subcommands")
	}

	switch os.Args[1] {
	case "file":
		fileCmd.Parse(os.Args[2:])

		fmt.Println(parseFile(fileCmd.Arg(0)))
	case "folders":
		foldersCmd.Parse(os.Args[2:])

		parseFolders(foldersCmd.Arg(0))
	default:
		log.Fatal("expected 'file' or 'folders' subcommands")
	}
}

type recipe struct {
	name         string
	yields       string
	ingredients  []string
	instructions []string
	date         string
}

func parseFolders(dirname string) {
	directories, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range directories {
		if d.IsDir() {
			fmt.Println("Entering '" + d.Name() + "'")

			files, err := ioutil.ReadDir(filepath.Join(dirname, d.Name()))
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				parseFile(filepath.Join(dirname, d.Name(), f.Name()))
			}
		}
	}
}

func parseFile(filename string) recipe {
	fmt.Println("Parsing '" + filename + "'")

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	name, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal("expected name: ", err)
	}

	name = bytes.TrimSpace(name)
	if len(name) == 0 {
		log.Fatal("expected name")
	}

	if empty, err := r.ReadBytes('\n'); err != nil || len(bytes.TrimSpace(empty)) > 0 {
		log.Fatal("expected empty line: ", err)
	}

	yields, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal("expected yields: ", err)
	}

	yields = bytes.TrimSpace(yields)
	if len(yields) == 0 {
		log.Fatal("expected yields")
	}

	if empty, err := r.ReadBytes('\n'); err != nil || len(bytes.TrimSpace(empty)) > 0 {
		log.Fatal("expected empty line: ", err)
	}

	var ingredients []string
	for {
		ingredient, err := r.ReadBytes('\n')
		if err != nil {
			log.Fatal("expected ingredient: ", err)
		}

		ingredient = bytes.TrimSpace(ingredient)
		if len(ingredient) == 0 {
			break
		}

		ingredients = append(ingredients, string(ingredient))
	}

	var instructions []string
	for {
		instruction, err := r.ReadBytes('\n')
		if err != nil {
			log.Fatal("expected instruction: ", err)
		}

		instruction = bytes.TrimSpace(instruction)
		if len(instruction) == 0 {
			break
		}

		instructions = append(instructions, string(instruction))
	}

	date, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal("expected date: ", err)
	}

	date = bytes.TrimSpace(date)
	if len(date) == 0 {
		log.Fatal("expected date")
	}

	for {
		if empty, err := r.ReadBytes('\n'); err == io.EOF {
			break
		} else if err != nil || len(bytes.TrimSpace(empty)) > 0 {
			log.Fatal("unexpected trailing data")
		}
	}

	return recipe{
		name:         string(name),
		yields:       string(yields),
		ingredients:  ingredients,
		instructions: instructions,
		date:         string(date),
	}
}
