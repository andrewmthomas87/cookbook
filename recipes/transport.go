package recipes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the recipes service.
func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}

	getRecipesHandler := kithttp.NewServer(
		makeGetRecipesEndpoint(s),
		nopDecodeRequest,
		encodeResponse,
		opts...,
	)
	getRecipeHandler := kithttp.NewServer(
		makeGetRecipeEndpoint(s),
		decodeGetRecipeRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/recipes/recipes", getRecipesHandler).Methods("GET")
	r.Handle("/recipes/recipe/{id}", getRecipeHandler).Methods("GET")

	return r
}

func nopDecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetRecipeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		return nil, err
	}
	return getRecipeRequest{Id: models.RecipeID(id)}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
