package services

import (
	"context"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

type characterService struct {
	repo repositories.CharacterRepository
}

// NewCharacterService cria uma nova instância do serviço de personagens
func NewCharacterService(repo repositories.CharacterRepository) CharacterService {
	return &characterService{
		repo: repo,
	}
}

func (s *characterService) Create(ctx context.Context, character *entities.Character) error {
	return s.repo.Create(ctx, character)
}

func (s *characterService) Update(ctx context.Context, character *entities.Character) error {
	return s.repo.Update(ctx, character)
}

func (s *characterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *characterService) GetByID(ctx context.Context, id string) (*entities.Character, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *characterService) GetByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error) {
	return s.repo.GetByUserAndGuild(ctx, userID, guildID)
}

func (s *characterService) ListByUser(ctx context.Context, userID string) ([]*entities.Character, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *characterService) ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error) {
	return s.repo.ListByGuild(ctx, guildID)
}

// List lista todos os personagens ativos
func (s *characterService) List(ctx context.Context) ([]*entities.Character, error) {
	return s.repo.List(ctx)
}

func (s *characterService) Search(ctx context.Context, query string) ([]*entities.Character, error) {
	return s.repo.Search(ctx, query)
}

func (s *characterService) CountByGuild(ctx context.Context, guildID string) (int64, error) {
	return s.repo.CountByGuild(ctx, guildID)
}

func (s *characterService) CountByUser(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountByUser(ctx, userID)
}
