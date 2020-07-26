// Package recipes provides recipe APIs.
package recipes

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
)

// Service provides recipes operations.
type Service interface {
	GetRecipes(ctx context.Context) ([]*models.Recipe, error)
	GetRecipe(ctx context.Context, id models.RecipeID) (*models.Recipe, []*models.Ingredient, []*models.Instruction, error)
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

func (s *service) GetRecipe(ctx context.Context, id models.RecipeID) (*models.Recipe, []*models.Ingredient, []*models.Instruction, error) {
	r, err := s.rr.GetRecipe(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	ingredients, err := s.rr.GetIngredients(ctx, r.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	instructions, err := s.rr.GetInstructions(ctx, r.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	return r, ingredients, instructions, nil
}
