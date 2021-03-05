package postgres

import (
	"context"
	"fmt"

	"github.com/riskiramdan/evos/internal/character"
	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/types"
)

// Storage implements the character storage service interface
type Storage struct {
	Storage data.GenericStorage
}

// FindAll find all characters
func (s *Storage) FindAll(ctx context.Context, params *character.FindAllCharacterParams) ([]*character.Characters, *types.Error) {

	characters := []*character.Characters{}
	where := `"deletedAt" IS NULL`

	if params.ID != 0 {
		where += ` AND "id" = :id`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "createdAt" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "createdAt" DESC`, where)
	}

	err := s.Storage.Where(ctx, &characters, where, map[string]interface{}{
		"id":     params.ID,
		"limit":  params.Limit,
		"name":   "%" + params.Name + "%",
		"offset": ((params.Page - 1) * params.Limit),
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".CharacterStorage->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return characters, nil
}

// FindByID find character by its id
func (s *Storage) FindByID(ctx context.Context, characterID int) (*character.Characters, *types.Error) {
	characters, err := s.FindAll(ctx, &character.FindAllCharacterParams{
		ID: characterID,
	})
	if err != nil {
		err.Path = ".CharacterStorage->FindByID()" + err.Path
		return nil, err
	}

	if len(characters) < 1 || characters[0].ID != characterID {
		return nil, &types.Error{
			Path:    ".CharacterStorage->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return characters[0], nil
}

// Insert insert character
func (s *Storage) Insert(ctx context.Context, character *character.Characters) (*character.Characters, *types.Error) {
	err := s.Storage.Insert(ctx, character)
	if err != nil {
		return nil, &types.Error{
			Path:    ".CharacterStorage->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return character, nil
}

// Update update character
func (s *Storage) Update(ctx context.Context, character *character.Characters) (*character.Characters, *types.Error) {
	err := s.Storage.Update(ctx, character)
	if err != nil {
		return nil, &types.Error{
			Path:    ".CharacterStorage->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return character, nil
}

// Delete delete a character
func (s *Storage) Delete(ctx context.Context, characterID int) *types.Error {
	err := s.Storage.Delete(ctx, characterID)
	if err != nil {
		return &types.Error{
			Path:    ".CharacterStorage->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return nil
}

// NewPostgresStorage creates new character repository service
func NewPostgresStorage(
	storage data.GenericStorage,
) *Storage {
	return &Storage{
		Storage: storage,
	}
}
