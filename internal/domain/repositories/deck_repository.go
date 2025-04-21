package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// DeckRepository define a interface para persistência de decks
type DeckRepository interface {
	// Create armazena um novo deck
	Create(ctx context.Context, deck *entities.Deck) error

	// Update atualiza um deck existente
	Update(ctx context.Context, deck *entities.Deck) error

	// Delete remove um deck
	Delete(ctx context.Context, id string) error

	// FindByID busca um deck pelo ID
	FindByID(ctx context.Context, id string) (*entities.Deck, error)

	// FindByUser busca todos os decks de um usuário
	FindByUser(ctx context.Context, userID string) ([]*entities.Deck, error)

	// FindByGuild busca todos os decks de um servidor
	FindByGuild(ctx context.Context, guildID string) ([]*entities.Deck, error)

	// FindByClass busca todos os decks de uma classe
	FindByClass(ctx context.Context, class string) ([]*entities.Deck, error)
}
