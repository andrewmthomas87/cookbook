package recipes

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/go-kit/kit/endpoint"
)

type getRecipesResponse struct {
	Rs  []Recipe `json:"rs,omitempty"`
	Err error    `json:"err,omitempty"`
}

func (r getRecipesResponse) error() error { return r.Err }

func makeGetRecipesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := s.GetRecipes(ctx)
		return getRecipesResponse{
			Rs:  convertRecipes(rs),
			Err: err,
		}, nil
	}
}

type getRecipeRequest struct {
	Id models.RecipeID `json:"id,omitempty"`
}

type getRecipeResponse struct {
	R            Recipe        `json:"r,omitempty"`
	Ingredients  []Ingredient  `json:"ingredients,omitempty"`
	Instructions []Instruction `json:"instructions,omitempty"`
	Err          error         `json:"err,omitempty"`
}

func (r getRecipeResponse) error() error { return r.Err }

func makeGetRecipeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRecipeRequest)
		r, ingredients, instructions, err := s.GetRecipe(ctx, req.Id)
		return getRecipeResponse{
			R:            convertRecipe(r),
			Ingredients:  convertIngredients(ingredients),
			Instructions: convertInstructions(instructions),
			Err:          err,
		}, nil
	}
}

// Recipe is a read model for models.Recipe.
type Recipe struct {
	ID       models.RecipeID `json:"id,omitempty"`
	Category string          `json:"category,omitempty"`
	Name     string          `json:"name,omitempty"`
	Yields   string          `json:"yields,omitempty"`
	Updated  string          `json:"updated,omitempty"`
}

// Ingredient is a read model for models.Ingredient.
type Ingredient struct {
	ID       models.IngredientID `json:"id,omitempty"`
	RecipeID models.RecipeID     `json:"recipeId,omitempty"`
	Value    string              `json:"value,omitempty"`
}

// Instruction is a read model for models.Instruction.
type Instruction struct {
	ID       models.InstructionID `json:"id,omitempty"`
	RecipeID models.RecipeID      `json:"recipeId,omitempty"`
	Value    string               `json:"value,omitempty"`
}

func convertRecipes(recipes []*models.Recipe) []Recipe {
	var rs []Recipe
	for _, r := range recipes {
		rs = append(rs, Recipe{
			ID:       r.ID,
			Category: r.Category,
			Name:     r.Name,
			Yields:   r.Yields,
			Updated:  r.Updated,
		})
	}
	return rs
}

func convertRecipe(recipe *models.Recipe) Recipe {
	if recipe == nil {
		return Recipe{}
	}
	return Recipe{
		ID:       recipe.ID,
		Category: recipe.Category,
		Name:     recipe.Name,
		Yields:   recipe.Yields,
		Updated:  recipe.Updated,
	}
}

func convertIngredients(ingredients []*models.Ingredient) []Ingredient {
	var is []Ingredient
	for _, i := range ingredients {
		is = append(is, Ingredient{
			ID:       i.ID,
			RecipeID: i.RecipeID,
			Value:    i.Value,
		})
	}
	return is
}

func convertInstructions(instructions []*models.Instruction) []Instruction {
	var is []Instruction
	for _, i := range instructions {
		is = append(is, Instruction{
			ID:       i.ID,
			RecipeID: i.RecipeID,
			Value:    i.Value,
		})
	}
	return is
}
