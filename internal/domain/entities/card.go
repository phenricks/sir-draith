package entities

import (
	"errors"
	"time"
)

// Card-related errors
var (
	ErrInvalidID     = errors.New("invalid card ID")
	ErrInvalidName   = errors.New("invalid card name")
	ErrInvalidType   = errors.New("invalid card type")
	ErrInvalidRarity = errors.New("invalid card rarity")
	ErrNoEffects     = errors.New("card must have at least one effect")
)

// CardType represents the type of a card
type CardType string

const (
	CardTypeSpell    CardType = "spell"
	CardTypeCreature CardType = "creature"
	CardTypeArtifact CardType = "artifact"
	CardTypeEnchant  CardType = "enchant"
)

// CardRarity represents the rarity of a card
type CardRarity string

const (
	CardRarityCommon    CardRarity = "common"
	CardRarityUncommon  CardRarity = "uncommon"
	CardRarityRare      CardRarity = "rare"
	CardRarityEpic      CardRarity = "epic"
	CardRarityLegendary CardRarity = "legendary"
)

// Effect represents a card effect
type Effect struct {
	Type       string      `bson:"type"`
	Value      int         `bson:"value"`
	Duration   int         `bson:"duration"`
	Target     string      `bson:"target"`
	Conditions []string    `bson:"conditions,omitempty"`
	Properties interface{} `bson:"properties,omitempty"`
}

// Requirements represents the requirements to play a card
type Requirements struct {
	Level      int            `bson:"level"`
	Class      string         `bson:"class,omitempty"`
	Attributes map[string]int `bson:"attributes,omitempty"`
	Skills     []string       `bson:"skills,omitempty"`
	Resources  map[string]int `bson:"resources,omitempty"`
}

// Card represents a card in the game
type Card struct {
	ID          string     `bson:"_id"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	Type        CardType   `bson:"type"`
	Rarity      CardRarity `bson:"rarity"`
	Cost        int        `bson:"cost"`
	Attack      int        `bson:"attack,omitempty"`
	Defense     int        `bson:"defense,omitempty"`
	Effects     []string   `bson:"effects,omitempty"`
	Keywords    []string   `bson:"keywords,omitempty"`
	ImageURL    string     `bson:"image_url,omitempty"`
	CreatedAt   int64      `bson:"created_at"`
	UpdatedAt   int64      `bson:"updated_at"`
}

// NewCard creates a new card instance
func NewCard(id, name string, cardType CardType, rarity CardRarity, description string, cost int, attack, defense int, effects []string, keywords []string, imageURL string) *Card {
	now := time.Now()
	return &Card{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        cardType,
		Rarity:      rarity,
		Cost:        cost,
		Attack:      attack,
		Defense:     defense,
		Effects:     effects,
		Keywords:    keywords,
		ImageURL:    imageURL,
		CreatedAt:   now.Unix(),
		UpdatedAt:   now.Unix(),
	}
}

// Validate validates the card data
func (c *Card) Validate() error {
	if c.ID == "" {
		return ErrInvalidID
	}
	if c.Name == "" {
		return ErrInvalidName
	}
	if c.Type == "" {
		return ErrInvalidType
	}
	if c.Rarity == "" {
		return ErrInvalidRarity
	}
	if c.Cost < 0 {
		return errors.New("cost cannot be negative")
	}
	if c.Type == CardTypeCreature {
		if c.Attack < 0 {
			return errors.New("attack cannot be negative")
		}
		if c.Defense < 0 {
			return errors.New("defense cannot be negative")
		}
	}
	return nil
}
