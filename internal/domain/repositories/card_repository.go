package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// CardRepository defines the interface for card data persistence
type CardRepository interface {
	// Create stores a new card in the repository
	Create(ctx context.Context, card *entities.Card) error

	// Update modifies an existing card in the repository
	Update(ctx context.Context, card *entities.Card) error

	// Delete removes a card from the repository
	Delete(ctx context.Context, id string) error

	// FindByID retrieves a card by its ID
	FindByID(ctx context.Context, id string) (*entities.Card, error)

	// FindAll retrieves all cards from the repository
	FindAll(ctx context.Context) ([]*entities.Card, error)

	// FindByType retrieves all cards of a specific type
	FindByType(ctx context.Context, cardType entities.CardType) ([]*entities.Card, error)

	// FindByRarity retrieves all cards of a specific rarity
	FindByRarity(ctx context.Context, rarity entities.CardRarity) ([]*entities.Card, error)
}
