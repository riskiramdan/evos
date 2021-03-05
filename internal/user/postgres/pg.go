package postgres

import (
	"context"
	"fmt"

	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/types"
	"github.com/riskiramdan/evos/internal/user"
)

// Storage implements the user storage service interface
type Storage struct {
	Storage data.GenericStorage
}

// FindAll find all users
func (s *Storage) FindAll(ctx context.Context, params *user.FindAllUsersParams) ([]*user.Users, *types.Error) {

	users := []*user.Users{}
	where := `"deletedAt" IS NULL`

	if params.ID != 0 {
		where += ` AND "id" = :id`
	}
	if params.Phone != "" {
		where += ` AND "phone" = :phone`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if params.Token != "" {
		where += ` AND "token" = :token`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "createdAt" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "createdAt" DESC`, where)
	}

	err := s.Storage.Where(ctx, &users, where, map[string]interface{}{
		"id":     params.ID,
		"limit":  params.Limit,
		"phone":  params.Phone,
		"name":   "%" + params.Name + "%",
		"offset": ((params.Page - 1) * params.Limit),
		"token":  params.Token,
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserStorage->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return users, nil
}

// FindByID find user by its id
func (s *Storage) FindByID(ctx context.Context, userID int) (*user.Users, *types.Error) {
	users, err := s.FindAll(ctx, &user.FindAllUsersParams{
		ID: userID,
	})
	if err != nil {
		err.Path = ".UserStorage->FindByID()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].ID != userID {
		return nil, &types.Error{
			Path:    ".UserStorage->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// FindByPhone find user by its phone
func (s *Storage) FindByPhone(ctx context.Context, phone string) (*user.Users, *types.Error) {
	users, err := s.FindAll(ctx, &user.FindAllUsersParams{
		Phone: phone,
	})
	if err != nil {
		err.Path = ".UserStorage->FindByPhone()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].Phone != phone {
		return nil, &types.Error{
			Path:    ".UserStorage->FindByPhone()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// FindByToken find user by its token
func (s *Storage) FindByToken(ctx context.Context, token string) (*user.Users, *types.Error) {
	users, err := s.FindAll(ctx, &user.FindAllUsersParams{
		Token: token,
	})
	if err != nil {
		err.Path = ".UserStorage->FindByToken()" + err.Path
		return nil, err
	}

	if len(users) < 1 || (users[0].Token != nil && *users[0].Token != token) {
		return nil, &types.Error{
			Path:    ".UserStorage->FindByToken()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// Insert insert user
func (s *Storage) Insert(ctx context.Context, user *user.Users) (*user.Users, *types.Error) {
	err := s.Storage.Insert(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserStorage->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return user, nil
}

// Update update user
func (s *Storage) Update(ctx context.Context, user *user.Users) (*user.Users, *types.Error) {
	err := s.Storage.Update(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserStorage->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return user, nil
}

// Delete delete a user
func (s *Storage) Delete(ctx context.Context, userID int) *types.Error {
	err := s.Storage.Delete(ctx, userID)
	if err != nil {
		return &types.Error{
			Path:    ".UserStorage->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return nil
}

// NewPostgresStorage creates new user repository service
func NewPostgresStorage(
	storage data.GenericStorage,
) *Storage {
	return &Storage{
		Storage: storage,
	}
}
