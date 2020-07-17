package recipes

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/go-kit/kit/endpoint"
)

type getRecipesResponse struct {
	rs  []Recipe `json:"rs,omitempty"`
	Err error    `json:"err,omitempty"`
}

func (r getRecipesResponse) error() error { return r.Err }

func makeGetRecipesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := s.GetRecipes(ctx)
		if err != nil {
			return nil, err
		}
		return getRecipesResponse{
			rs:  convertRecipes(rs),
			Err: err,
		}, nil
	}
}

// Recipe is a read model for models.Recipe.
type Recipe struct {
	ID      models.RecipeID `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Yields  string          `json:"yields,omitempty"`
	Updated string          `json:"updated,omitempty"`
}

func convertRecipes(recipes []*models.Recipe) []Recipe {
	var rs []Recipe
	for _, r := range recipes {
		rs = append(rs, Recipe{
			ID:      r.ID,
			Name:    r.Name,
			Yields:  r.Yields,
			Updated: r.Updated,
		})
	}
	return rs
}
