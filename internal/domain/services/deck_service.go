package services

import (
	"context"
	"fmt"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// DeckService encapsula a lógica de negócio relacionada a decks
type DeckService struct {
	deckRepo repositories.DeckRepository
	cardRepo repositories.CardRepository
}

// NewDeckService cria uma nova instância do serviço de decks
func NewDeckService(deckRepo repositories.DeckRepository, cardRepo repositories.CardRepository) *DeckService {
	return &DeckService{
		deckRepo: deckRepo,
		cardRepo: cardRepo,
	}
}

// CreateDeck cria um novo deck
func (s *DeckService) CreateDeck(ctx context.Context, userID, guildID, name, description, class string) (*entities.Deck, error) {
	deck := entities.NewDeck(userID, guildID, name, description, class)
	if err := s.deckRepo.Create(ctx, deck); err != nil {
		return nil, fmt.Errorf("erro ao criar deck: %w", err)
	}
	return deck, nil
}

// AddCardToDeck adiciona uma carta ao deck
func (s *DeckService) AddCardToDeck(ctx context.Context, deckID string, cardID string) error {
	deck, err := s.deckRepo.FindByID(ctx, deckID)
	if err != nil {
		return fmt.Errorf("erro ao buscar deck: %w", err)
	}
	if deck == nil {
		return fmt.Errorf("deck não encontrado")
	}

	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return fmt.Errorf("erro ao buscar carta: %w", err)
	}
	if card == nil {
		return fmt.Errorf("carta não encontrada")
	}

	config := entities.DefaultDeckConfig()
	if err := deck.AddCard(cardID, config); err != nil {
		return fmt.Errorf("erro ao adicionar carta: %w", err)
	}

	return s.deckRepo.Update(ctx, deck)
}

// RemoveCardFromDeck remove uma carta do deck
func (s *DeckService) RemoveCardFromDeck(ctx context.Context, deckID string, cardID string) error {
	deck, err := s.deckRepo.FindByID(ctx, deckID)
	if err != nil {
		return fmt.Errorf("erro ao buscar deck: %w", err)
	}
	if deck == nil {
		return fmt.Errorf("deck não encontrado")
	}

	if err := deck.RemoveCard(cardID); err != nil {
		return fmt.Errorf("erro ao remover carta: %w", err)
	}

	return s.deckRepo.Update(ctx, deck)
}

// GetDeck retorna um deck pelo ID
func (s *DeckService) GetDeck(ctx context.Context, deckID string) (*entities.Deck, error) {
	return s.deckRepo.FindByID(ctx, deckID)
}

// ListDecksByUser lista todos os decks de um usuário
func (s *DeckService) ListDecksByUser(ctx context.Context, userID string) ([]*entities.Deck, error) {
	return s.deckRepo.FindByUser(ctx, userID)
}

// ListDecksByGuild lista todos os decks de um servidor
func (s *DeckService) ListDecksByGuild(ctx context.Context, guildID string) ([]*entities.Deck, error) {
	return s.deckRepo.FindByGuild(ctx, guildID)
}

// ValidateDeck valida um deck de acordo com as regras
func (s *DeckService) ValidateDeck(ctx context.Context, deckID string) error {
	deck, err := s.deckRepo.FindByID(ctx, deckID)
	if err != nil {
		return fmt.Errorf("erro ao buscar deck: %w", err)
	}
	if deck == nil {
		return fmt.Errorf("deck não encontrado")
	}

	config := entities.DefaultDeckConfig()
	return deck.Validate(config)
}

// DeleteDeck remove um deck
func (s *DeckService) DeleteDeck(ctx context.Context, deckID string) error {
	return s.deckRepo.Delete(ctx, deckID)
}
