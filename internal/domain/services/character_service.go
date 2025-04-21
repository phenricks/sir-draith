package services

import (
	"context"
	"fmt"
	"time"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"
	"sirdraith/internal/domain/repositories"
)

// CharacterService gerencia as operações de negócio relacionadas aos personagens
type CharacterService struct {
	repo repositories.CharacterRepository
}

// NewCharacterService cria uma nova instância do serviço de personagens
func NewCharacterService(repo repositories.CharacterRepository) *CharacterService {
	return &CharacterService{
		repo: repo,
	}
}

// CreateCharacter cria um novo personagem
func (s *CharacterService) CreateCharacter(ctx context.Context, userID, guildID, name string) (*entities.Character, error) {
	// Verifica se o usuário já tem um personagem neste servidor
	existing, err := s.repo.GetByUserAndGuild(ctx, userID, guildID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("usuário já possui um personagem neste servidor")
	}

	// Cria um novo personagem com valores padrão
	character := &entities.Character{
		UserID:     userID,
		GuildID:    guildID,
		Name:       name,
		Level:      1,
		Experience: 0,
		Gold:       gamedata.StartingGold,
		Inventory:  make([]entities.Item, 0),
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, character); err != nil {
		return nil, fmt.Errorf("erro ao criar personagem: %w", err)
	}

	return character, nil
}

// GetCharacter busca um personagem pelo ID
func (s *CharacterService) GetCharacter(ctx context.Context, id string) (*entities.Character, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCharacterByUserAndGuild busca um personagem pelo ID do usuário e do servidor
func (s *CharacterService) GetCharacterByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error) {
	return s.repo.GetByUserAndGuild(ctx, userID, guildID)
}

// UpdateCharacter atualiza os dados de um personagem
func (s *CharacterService) UpdateCharacter(ctx context.Context, character *entities.Character) error {
	return s.repo.Update(ctx, character)
}

// DeleteCharacter marca um personagem como inativo
func (s *CharacterService) DeleteCharacter(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// AddExperience adiciona experiência ao personagem e gerencia level up
func (s *CharacterService) AddExperience(ctx context.Context, character *entities.Character, exp int) error {
	character.AddExperience(exp)
	return s.repo.Update(ctx, character)
}

// AddGold adiciona ou remove ouro do personagem
func (s *CharacterService) AddGold(ctx context.Context, character *entities.Character, amount int) error {
	character.Gold += amount
	if character.Gold < 0 {
		character.Gold = 0
	}
	return s.repo.Update(ctx, character)
}

// AddItem adiciona um item ao inventário do personagem
func (s *CharacterService) AddItem(ctx context.Context, character *entities.Character, item entities.Item) error {
	if err := character.AddItem(item); err != nil {
		return fmt.Errorf("erro ao adicionar item: %w", err)
	}
	return s.repo.Update(ctx, character)
}

// RemoveItem remove um item do inventário do personagem
func (s *CharacterService) RemoveItem(ctx context.Context, character *entities.Character, itemName string, quantity int) error {
	if err := character.RemoveItem(itemName, quantity); err != nil {
		return fmt.Errorf("erro ao remover item: %w", err)
	}
	return s.repo.Update(ctx, character)
}

// EquipItem equipa um item do inventário
func (s *CharacterService) EquipItem(ctx context.Context, character *entities.Character, itemName string) error {
	if err := character.EquipItem(itemName); err != nil {
		return fmt.Errorf("erro ao equipar item: %w", err)
	}
	return s.repo.Update(ctx, character)
}

// UnequipItem desequipa um item
func (s *CharacterService) UnequipItem(ctx context.Context, character *entities.Character, itemName string) error {
	if err := character.UnequipItem(itemName); err != nil {
		return fmt.Errorf("erro ao desequipar item: %w", err)
	}
	return s.repo.Update(ctx, character)
}

// TakeDamage aplica dano ao personagem
func (s *CharacterService) TakeDamage(ctx context.Context, character *entities.Character, damage int) error {
	character.TakeDamage(damage)
	return s.repo.Update(ctx, character)
}

// Heal cura o personagem
func (s *CharacterService) Heal(ctx context.Context, character *entities.Character, amount int) error {
	character.Heal(amount)
	return s.repo.Update(ctx, character)
}

// ListCharactersByGuild lista todos os personagens de um servidor
func (s *CharacterService) ListCharactersByGuild(ctx context.Context, guildID string) ([]*entities.Character, error) {
	return s.repo.ListByGuild(ctx, guildID)
}

// SearchCharacters busca personagens por nome ou título
func (s *CharacterService) SearchCharacters(ctx context.Context, query string) ([]*entities.Character, error) {
	return s.repo.Search(ctx, query)
}
