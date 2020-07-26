package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrewmthomas87/cookbook/database"
	"github.com/andrewmthomas87/cookbook/recipes"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/cookbook")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		panic(err)
	}
	config.ConnConfig.Database = viper.GetString("db.database")
	config.ConnConfig.User = viper.GetString("db.user")
	config.ConnConfig.Password = viper.GetString("db.password")

	p, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	rr := database.NewRecipesRepository(p)

	var rec recipes.Service
	rec = recipes.NewService(rr)
	rec = recipes.NewLoggingService(log.With(logger, "component", "recipes"), rec)

	httpLogger := log.With(logger, "component", "http")
	mux := http.NewServeMux()

	mux.Handle("/recipes/", recipes.MakeHandler(rec, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		httpAddr := viper.GetString("server.httpAddress")
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
