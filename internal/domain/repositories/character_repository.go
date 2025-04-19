package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// CharacterRepository define as operações de persistência para personagens
type CharacterRepository interface {
	// Create cria um novo personagem
	Create(ctx context.Context, character *entities.Character) error

	// GetByID busca um personagem pelo ID
	GetByID(ctx context.Context, id string) (*entities.Character, error)

	// GetByUserID busca todos os personagens de um usuário
	GetByUserID(ctx context.Context, userID string) ([]*entities.Character, error)

	// GetByGuildID busca todos os personagens de uma guilda
	GetByGuildID(ctx context.Context, guildID string) ([]*entities.Character, error)

	// GetByUserAndGuild busca um personagem específico de um usuário em uma guilda
	GetByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error)

	// Update atualiza um personagem existente
	Update(ctx context.Context, character *entities.Character) error

	// Delete remove um personagem
	Delete(ctx context.Context, id string) error

	// List lista todos os personagens ativos
	List(ctx context.Context) ([]*entities.Character, error)

	// ListByUser lista todos os personagens de um usuário
	ListByUser(ctx context.Context, userID string) ([]*entities.Character, error)

	// ListByGuild lista todos os personagens de uma guilda
	ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error)

	// Search busca personagens por nome ou título
	Search(ctx context.Context, query string) ([]*entities.Character, error)

	// CountByGuild conta o número de personagens em uma guilda
	CountByGuild(ctx context.Context, guildID string) (int64, error)

	// CountByUser conta o número de personagens de um usuário
	CountByUser(ctx context.Context, userID string) (int64, error)
}
