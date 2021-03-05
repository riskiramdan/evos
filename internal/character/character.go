package character

import (
	"context"
	"errors"
	"time"

	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/types"
)

// Errors
var (
	ErrInvalidPower    = errors.New("Invalid Power")
	ErrCharacterExists = errors.New("Character already exists")
)

// Characters character
type Characters struct {
	ID              int        `json:"id" db:"id"`
	CharacterTypeID int        `json:"characterTypeID" db:"characterTypeID"`
	Name            string     `json:"name" db:"name"`
	Power           int        `json:"power" db:"power"`
	Value           int        `json:"value"`
	CreatedAt       time.Time  `json:"createdAt" db:"createdAt"`
	CreatedBy       string     `json:"createdBy" db:"createdBy"`
	UpdatedAt       *time.Time `json:"updatedAt" db:"updatedAt"`
	UpdatedBy       string     `json:"updatedBy" db:"updatedBy"`
}

//FindAllCharacterParams params for find all
type FindAllCharacterParams struct {
	ID    int    `json:"id"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Name  string `json:"name"`
}

// TransactionParams params for transaction
type TransactionParams struct {
	CharacterTypeID int    `json:"characterTypeID,omitempty"`
	Name            string `json:"name"`
	Power           *int   `json:"power,omitempty"`
}

// Storage represents the character storage interface
type Storage interface {
	FindAll(ctx context.Context, params *FindAllCharacterParams) ([]*Characters, *types.Error)
	FindByID(ctx context.Context, characterID int) (*Characters, *types.Error)
	Insert(ctx context.Context, character *Characters) (*Characters, *types.Error)
	Update(ctx context.Context, character *Characters) (*Characters, *types.Error)
	Delete(ctx context.Context, characterID int) *types.Error
}

// ServiceInterface represents the character service interface
type ServiceInterface interface {
	ListCharacters(ctx context.Context, params *FindAllCharacterParams) ([]*Characters, int, *types.Error)
	GetCharacter(ctx context.Context, characterID int) (*Characters, *types.Error)
	CreateCharacter(ctx context.Context, params *TransactionParams) (*Characters, *types.Error)
	UpdateCharacter(ctx context.Context, characterID int, params *TransactionParams) (*Characters, *types.Error)
}

// Service is the domain logic implementation of character Service interface
type Service struct {
	characterStorage Storage
}

func calculateValue(power, characterTypeID int) int {
	value := 0
	if characterTypeID == 0 {
		return value
	}
	switch characterTypeID {
	case 1:
		{
			value = calculatePercentage(power, 150)
		}
	case 2:
		{
			value = 2 + calculatePercentage(power, 110)
		}
	case 3:
		{
			value = calculatePercentage(power, 300)
			if power < 20 {
				value = calculatePercentage(power, 200)
			}
		}
	default:
		{
			value = 0
		}
	}
	return value
}

func calculatePercentage(value, amount int) int {
	return value * amount / 100
}

// ListCharacters is listing characters
func (s *Service) ListCharacters(ctx context.Context, params *FindAllCharacterParams) ([]*Characters, int, *types.Error) {
	characters, err := s.characterStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".characterservice->Listcharacters()" + err.Path
		return nil, 0, err
	}
	params.Page = 0
	params.Limit = 0
	allcharacters, err := s.characterStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".characterservice->Listcharacters()" + err.Path
		return nil, 0, err
	}

	for _, v := range characters {
		v.Value = calculateValue(v.Power, v.CharacterTypeID)
	}

	return characters, len(allcharacters), nil
}

// GetCharacter is get character
func (s *Service) GetCharacter(ctx context.Context, characterID int) (*Characters, *types.Error) {
	character, err := s.characterStorage.FindByID(ctx, characterID)
	if err != nil {
		err.Path = ".characterservice->GetCharacter()" + err.Path
		return nil, err
	}

	return character, nil
}

// CreateCharacter create character
func (s *Service) CreateCharacter(ctx context.Context, params *TransactionParams) (*Characters, *types.Error) {
	characters, _, errType := s.ListCharacters(ctx, &FindAllCharacterParams{
		Name: params.Name,
	})
	if errType != nil {
		errType.Path = ".characterservice->CreateCharacter()" + errType.Path
		return nil, errType
	}
	if len(characters) > 0 {
		return nil, &types.Error{
			Path:    ".characterservice->CreateCharacter()",
			Message: ErrCharacterExists.Error(),
			Error:   ErrCharacterExists,
			Type:    "validation-error",
		}
	}

	now := time.Now()

	character := &Characters{
		Name:            params.Name,
		CharacterTypeID: params.CharacterTypeID,
		Power:           *params.Power,
		CreatedBy:       "Admin",
		CreatedAt:       now,
		UpdatedBy:       "Admin",
		UpdatedAt:       &now,
	}

	character, errType = s.characterStorage.Insert(ctx, character)
	if errType != nil {
		errType.Path = ".characterservice->CreateCharacter()" + errType.Path
		return nil, errType
	}

	return character, nil
}

// UpdateCharacter update a character
func (s *Service) UpdateCharacter(ctx context.Context, characterID int, params *TransactionParams) (*Characters, *types.Error) {
	character, err := s.GetCharacter(ctx, characterID)
	if err != nil {
		err.Path = ".CharacterService->UpdateCharacter()" + err.Path
		return nil, err
	}

	if params.Name != "" {
		characters, _, err := s.ListCharacters(ctx, &FindAllCharacterParams{
			Name: params.Name,
		})
		if err != nil {
			err.Path = ".CharacterService->UpdateCharacter()" + err.Path
			return nil, err
		}
		if len(characters) > 0 {
			return nil, &types.Error{
				Path:    ".CharacterService->CreateCharacter()",
				Message: data.ErrAlreadyExist.Error(),
				Error:   data.ErrAlreadyExist,
				Type:    "validation-error",
			}
		}
		character.Name = params.Name
	}
	if params.Power != nil {
		character.Power = *params.Power
	}

	now := time.Now()
	character.UpdatedAt = &now

	character, err = s.characterStorage.Update(ctx, character)
	if err != nil {
		err.Path = ".CharacterService->UpdateCharacter()" + err.Path
		return nil, err
	}

	return character, nil
}

// NewService creates a new character AppService
func NewService(
	characterStorage Storage,
) *Service {
	return &Service{
		characterStorage: characterStorage,
	}
}
