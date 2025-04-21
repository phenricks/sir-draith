package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrInvalidDeckID indicates an invalid deck ID
	ErrInvalidDeckID = errors.New("invalid deck ID")
	// ErrInvalidDeckName indicates an invalid deck name
	ErrInvalidDeckName = errors.New("invalid deck name")
	// ErrInvalidDeckClass indicates an invalid deck class
	ErrInvalidDeckClass = errors.New("invalid deck class")
	// ErrDeckCardLimit indicates the deck has reached its card limit
	ErrDeckCardLimit = errors.New("deck card limit reached")
	// ErrCardNotFound indicates a card was not found in the deck
	ErrCardNotFound = errors.New("card not found in deck")
)

// DeckConfig defines configuration rules for a deck
type DeckConfig struct {
	MaxCards    int            // Maximum number of cards allowed in the deck
	MaxPerCard  int            // Maximum number of copies of a single card
	CardLimits  map[string]int // Special limits for specific cards
	ClassLimits map[string]int // Limits for cards of specific classes
}

// DefaultDeckConfig returns the default deck configuration
func DefaultDeckConfig() *DeckConfig {
	return &DeckConfig{
		MaxCards:    30,
		MaxPerCard:  2,
		CardLimits:  make(map[string]int),
		ClassLimits: make(map[string]int),
	}
}

// Deck represents a collection of cards
type Deck struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	GuildID     string             `bson:"guild_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Class       string             `bson:"class"`
	Cards       map[string]int     `bson:"cards"` // Map of card IDs to quantity
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// NewDeck creates a new deck instance
func NewDeck(userID, guildID, name, description, class string) *Deck {
	now := time.Now()
	return &Deck{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		GuildID:     guildID,
		Name:        name,
		Description: description,
		Class:       class,
		Cards:       make(map[string]int),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate checks if the deck is valid according to the configuration
func (d *Deck) Validate(config *DeckConfig) error {
	if d.ID.IsZero() {
		return ErrInvalidDeckID
	}
	if d.Name == "" {
		return ErrInvalidDeckName
	}
	if d.Class == "" {
		return ErrInvalidDeckClass
	}

	totalCards := 0
	for _, quantity := range d.Cards {
		totalCards += quantity
	}

	if totalCards > config.MaxCards {
		return ErrDeckCardLimit
	}

	return nil
}

// AddCard adds a card to the deck
func (d *Deck) AddCard(cardID string, config *DeckConfig) error {
	currentQuantity := d.Cards[cardID]
	if currentQuantity >= config.MaxPerCard {
		return ErrDeckCardLimit
	}

	totalCards := 0
	for _, quantity := range d.Cards {
		totalCards += quantity
	}
	if totalCards >= config.MaxCards {
		return ErrDeckCardLimit
	}

	d.Cards[cardID] = currentQuantity + 1
	d.UpdatedAt = time.Now()
	return nil
}

// RemoveCard removes a card from the deck
func (d *Deck) RemoveCard(cardID string) error {
	quantity, exists := d.Cards[cardID]
	if !exists || quantity <= 0 {
		return ErrCardNotFound
	}

	if quantity == 1 {
		delete(d.Cards, cardID)
	} else {
		d.Cards[cardID] = quantity - 1
	}

	d.UpdatedAt = time.Now()
	return nil
}

// GetCardCount returns the total number of cards in the deck
func (d *Deck) GetCardCount() int {
	total := 0
	for _, quantity := range d.Cards {
		total += quantity
	}
	return total
}

// GetCardQuantity returns the quantity of a specific card in the deck
func (d *Deck) GetCardQuantity(cardID string) int {
	return d.Cards[cardID]
}
