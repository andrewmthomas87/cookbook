package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4"
)

func main() {
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	config.Config.Database = "cookbook"
	config.Config.User = "postgres"
	config.Config.Password = "Asdf195789"
	conn, err := pgx.ConnectConfig(context.Background(), config);
	if err != nil {
		log.Fatal(err)
	}

	cs := parseFolders(os.Args[1])
	for category, rs := range cs {
		fmt.Println(category)
		for _, r := range rs {
			var id int
			if err := conn.QueryRow(context.Background(), "INSERT INTO recipes (category, name, yields, updated) VALUES ($1, $2, $3, $4) RETURNING id", category, r.name, r.yields, r.date).Scan(&id); err != nil {
				log.Fatal(err)
			}

			for _, ingredient := range r.ingredients {
				_, err := conn.Exec(context.Background(), "INSERT INTO ingredients (recipeid, value) VALUES ($1, $2)", id, ingredient)
				if err != nil {
					log.Fatal(err)
				}
			}
			for _, instruction := range r.instructions {
				_, err := conn.Exec(context.Background(), "INSERT INTO instructions (recipeid, value) VALUES ($1, $2)", id, instruction)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

type recipe struct {
	name         string
	yields       string
	ingredients  []string
	instructions []string
	date         string
}

func parseFolders(dirname string) map[string][]recipe {
	directories, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	cs := make(map[string][]recipe)
	for _, d := range directories {
		if d.IsDir() {
			files, err := ioutil.ReadDir(filepath.Join(dirname, d.Name()))
			if err != nil {
				log.Fatal(err)
			}

			var rs []recipe
			for _, f := range files {
				f := parseFile(filepath.Join(dirname, d.Name(), f.Name()))
				rs = append(rs, f)
			}
			cs[d.Name()] = rs
		}
	}
	return cs
}

func parseFile(filename string) recipe {
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
