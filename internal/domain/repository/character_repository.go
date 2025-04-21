package repository

import (
	"context"
	"errors"

	"sirdraith/internal/domain/entities"
)

var (
	ErrCharacterNotFound = errors.New("character not found")
)

// CharacterRepository define a interface para operações com personagens
type CharacterRepository interface {
	Create(ctx context.Context, character *entities.Character) error
	Update(ctx context.Context, character *entities.Character) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entities.Character, error)
	GetByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error)
	ListByUser(ctx context.Context, userID string) ([]*entities.Character, error)
	ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error)
	List(ctx context.Context) ([]*entities.Character, error)
	Search(ctx context.Context, query string) ([]*entities.Character, error)
	CountByGuild(ctx context.Context, guildID string) (int64, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
}
