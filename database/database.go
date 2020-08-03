// Package database provides database implementations of all the domain repositories.
package database

import (
	"context"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type recipesRepository struct {
	p *pgxpool.Pool
}

func (r *recipesRepository) CreateRecipe(ctx context.Context, init *models.Recipe) (models.RecipeID, error) {
	sql := "INSERT INTO recipes (category, name, yields, updated) VALUES (?, ?, ?, ?) RETURNING id"
	var id int
	if err := r.p.QueryRow(ctx, sql, init.Category, init.Name, init.Yields, init.Updated).Scan(&id); err != nil {
		return 0, err
	}
	return models.RecipeID(id), nil
}

func (r *recipesRepository) GetRecipes(ctx context.Context) ([]*models.Recipe, error) {
	sql := "SELECT id, category, name, yields, updated, image FROM recipes"
	rows, err := r.p.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rs []*models.Recipe
	for rows.Next() {
		var r models.Recipe
		if err := rows.Scan(&r.ID, &r.Category, &r.Name, &r.Yields, &r.Updated, &r.Image); err != nil {
			return nil, err
		}
		rs = append(rs, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *recipesRepository) GetRecipe(ctx context.Context, id models.RecipeID) (*models.Recipe, error) {
	sql := "SELECT id, category, name, yields, updated, image FROM recipes WHERE id=$1"
	var rec models.Recipe
	if err := r.p.QueryRow(ctx, sql, id).Scan(&rec.ID, &rec.Category, &rec.Name, &rec.Yields, &rec.Updated, &rec.Image); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *recipesRepository) CreateIngredient(ctx context.Context, init *models.Ingredient) (models.IngredientID, error) {
	sql := "INSERT INTO ingredients (recipeid, value) VALUES (?, ?)"
	var id int
	if err := r.p.QueryRow(ctx, sql, init.RecipeID, init.Value).Scan(&id); err != nil {
		return 0, err
	}
	return models.IngredientID(id), nil
}

func (r *recipesRepository) CreateIngredients(ctx context.Context, inits []*models.Ingredient) ([]models.IngredientID, error) {
	tx, err := r.p.Begin(ctx)
	if err != nil {
		return nil, err
	}
	sql := "INSERT INTO ingredients (recipeid, value) VALUES (?, ?) RETURNING id"
	var ids []models.IngredientID
	for _, init := range inits {
		var id int
		if err := tx.QueryRow(ctx, sql, init.RecipeID, init.Value).Scan(&id); err != nil {
			_ = tx.Rollback(ctx)
			return nil, err
		}
		ids = append(ids, models.IngredientID(id))
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *recipesRepository) GetIngredients(ctx context.Context, id models.RecipeID) ([]*models.Ingredient, error) {
	sql := "SELECT id, recipeid, value FROM ingredients WHERE recipeid=$1"
	rows, err := r.p.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var is []*models.Ingredient
	for rows.Next() {
		var i models.Ingredient
		if err := rows.Scan(&i.ID, &i.RecipeID, &i.Value); err != nil {
			return nil, err
		}
		is = append(is, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return is, nil
}

func (r *recipesRepository) CreateInstruction(ctx context.Context, init *models.Instruction) (models.InstructionID, error) {
	sql := "INSERT INTO instructions (recipeid, value) VALUES (?, ?) RETURNING id"
	var id int
	if err := r.p.QueryRow(ctx, sql, init.RecipeID, init.Value).Scan(&id); err != nil {
		return 0, err
	}
	return models.InstructionID(id), nil
}

func (r *recipesRepository) CreateInstructions(ctx context.Context, inits []*models.Instruction) ([]models.InstructionID, error) {
	tx, err := r.p.Begin(ctx)
	if err != nil {
		return nil, err
	}
	sql := "INSERT INTO instructions (recipeid, value) VALUES (?, ?) RETURNING id"
	var ids []models.InstructionID
	for _, init := range inits {
		var id int
		if err := tx.QueryRow(ctx, sql, init.RecipeID, init.Value).Scan(&id); err != nil {
			_ = tx.Rollback(ctx)
			return nil, err
		}
		ids = append(ids, models.InstructionID(id))
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *recipesRepository) GetInstructions(ctx context.Context, id models.RecipeID) ([]*models.Instruction, error) {
	sql := "SELECT id, recipeid, value FROM instructions WHERE recipeid=$1"
	rows, err := r.p.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var is []*models.Instruction
	for rows.Next() {
		var i models.Instruction
		if err := rows.Scan(&i.ID, &i.RecipeID, &i.Value); err != nil {
			return nil, err
		}
		is = append(is, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return is, nil
}

func NewRecipesRepository(p *pgxpool.Pool) models.RecipesRepository {
	return &recipesRepository{p: p}
}
