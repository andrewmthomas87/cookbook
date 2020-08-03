package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	is := parseImages(os.Args[1])
	for name, image := range is {
		if _, err := conn.Exec(context.Background(), "UPDATE recipes SET image=$1 WHERE name=$2", image, name); err != nil {
			log.Fatal(err)
		}
	}
}

func parseImages(dirname string) map[string]string {
	directories, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	is := make(map[string]string)
	for _, d := range directories {
		if d.IsDir() {
			files, err := ioutil.ReadDir(filepath.Join(dirname, d.Name()))
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".png") || strings.HasSuffix(f.Name(), ".jpg") || strings.HasSuffix(f.Name(), ".jpeg") {
					is[d.Name()] = f.Name()
					break
				}
			}
		}
	}
	return is
}
