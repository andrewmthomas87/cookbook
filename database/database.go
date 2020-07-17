// Package database provides database implementations of all the domain repositories.
package database

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/jmoiron/sqlx"
)

type recipesRepository struct {
	db *sqlx.DB
}

func (r *recipesRepository) CreateRecipe(ctx context.Context, init *models.Recipe) (models.RecipeID, error) {
	sql := "INSERT INTO recipes (name, yields, updated) VALUES (?, ?, ?)"
	res, err := r.db.ExecContext(ctx, sql, init.Name, init.Yields, init.Updated)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return models.RecipeID(id), nil
}

func (r *recipesRepository) GetRecipes(ctx context.Context) ([]*models.Recipe, error) {
	sql := "SELECT id, name, yields, updated FROM recipes"
	var rs []*models.Recipe
	if err := r.db.Select(&rs, sql); err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *recipesRepository) CreateIngredient(ctx context.Context, init *models.Ingredient) (models.IngredientID, error) {
	sql := "INSERT INTO ingredients (recipeid, value) VALUES (?, ?)"
	res, err := r.db.ExecContext(ctx, sql, init.RecipeID, init.Value)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return models.IngredientID(id), nil
}

func (r *recipesRepository) CreateIngredients(ctx context.Context, inits []*models.Ingredient) ([]models.IngredientID, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	sql := "INSERT INTO ingredients (recipeid, value) VALUES (?, ?)"
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	var ids []models.IngredientID
	for _, init := range inits {
		res, err := stmt.Exec(init.RecipeID, init.Value)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		ids = append(ids, models.IngredientID(id))
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *recipesRepository) CreateInstruction(ctx context.Context, init *models.Instruction) (models.InstructionID, error) {
	sql := "INSERT INTO instructions (recipeid, value) VALUES (?, ?)"
	res, err := r.db.ExecContext(ctx, sql, init.RecipeID, init.Value)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return models.InstructionID(id), nil
}

func (r *recipesRepository) CreateInstructions(ctx context.Context, inits []*models.Instruction) ([]models.InstructionID, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	sql := "INSERT INTO instructions (recipeid, value) VALUES (?, ?)"
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	var ids []models.InstructionID
	for _, init := range inits {
		res, err := stmt.Exec(init.RecipeID, init.Value)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		ids = append(ids, models.InstructionID(id))
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ids, nil
}

func NewRecipesRepository(db *sqlx.DB) models.RecipesRepository {
	return &recipesRepository{db: db}
}
