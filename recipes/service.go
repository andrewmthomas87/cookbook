// Package recipes provides recipe APIs.
package recipes

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
)

// Service provides recipes operations.
type Service interface {
	GetRecipes(ctx context.Context) ([]*models.Recipe, error)
}

type service struct {
	rr models.RecipesRepository
}

// NewService creates a recipes service with the necessary dependencies.
func NewService(rr models.RecipesRepository) Service {
	return &service{rr: rr}
}

func (s *service) GetRecipes(ctx context.Context) ([]*models.Recipe, error) {
	return s.rr.GetRecipes(ctx)
}
