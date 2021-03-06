package models

import (
	"context"
)

// RecipeID uniquely identifies a particular restaurant.
type RecipeID int64

// Recipe models a recipe.
type Recipe struct {
	ID           RecipeID
	Category     string
	Name         string
	Yields       string
	Updated      string
	Image        string
	Ingredients  []*Ingredient
	Instructions []*Instruction
}

// IngredientID uniquely identifies a particular ingredient.
type IngredientID int64

// Ingredient models an ingredient.
type Ingredient struct {
	ID       IngredientID
	RecipeID RecipeID
	Value    string
}

// InstructionID uniquely identifies a particular instruction.
type InstructionID int64

// Instruction models an instruction.
type Instruction struct {
	ID       InstructionID
	RecipeID RecipeID
	Value    string
}

// RecipesRepository provides access to a recipes store.
type RecipesRepository interface {
	CreateRecipe(ctx context.Context, init *Recipe) (RecipeID, error)
	GetRecipes(ctx context.Context) ([]*Recipe, error)
	GetRecipe(ctx context.Context, id RecipeID) (*Recipe, error)

	CreateIngredient(ctx context.Context, init *Ingredient) (IngredientID, error)
	CreateIngredients(ctx context.Context, inits []*Ingredient) ([]IngredientID, error)
	GetIngredients(ctx context.Context, id RecipeID) ([]*Ingredient, error)

	CreateInstruction(ctx context.Context, init *Instruction) (InstructionID, error)
	CreateInstructions(ctx context.Context, inits []*Instruction) ([]InstructionID, error)
	GetInstructions(ctx context.Context, id RecipeID) ([]*Instruction, error)
}
