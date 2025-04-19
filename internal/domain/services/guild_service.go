package services

import (
	"context"
	"time"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// GuildService handles business logic for guild operations
type GuildService struct {
	guildRepo repositories.GuildRepository
}

// NewGuildService creates a new instance of GuildService
func NewGuildService(guildRepo repositories.GuildRepository) *GuildService {
	return &GuildService{
		guildRepo: guildRepo,
	}
}

// CreateGuild creates a new guild from Discord information
func (s *GuildService) CreateGuild(ctx context.Context, discordID, name string) (*entities.Guild, error) {
	guild := &entities.Guild{
		DiscordID:       discordID,
		Name:            name,
		Prefix:          "!", // Default prefix
		ModeratorRoles:  []string{},
		EnabledFeatures: []string{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.guildRepo.Create(ctx, guild); err != nil {
		return nil, err
	}

	return guild, nil
}

// GetGuildByDiscordID retrieves a guild by its Discord ID
func (s *GuildService) GetGuildByDiscordID(ctx context.Context, discordID string) (*entities.Guild, error) {
	return s.guildRepo.FindByDiscordID(ctx, discordID)
}

// UpdateGuild updates guild information
func (s *GuildService) UpdateGuild(ctx context.Context, guild *entities.Guild) error {
	guild.UpdatedAt = time.Now()
	return s.guildRepo.Update(ctx, guild)
}

// UpdateGuildPrefix updates the guild's command prefix
func (s *GuildService) UpdateGuildPrefix(ctx context.Context, guildID, newPrefix string) error {
	guild, err := s.guildRepo.FindByDiscordID(ctx, guildID)
	if err != nil {
		return err
	}

	guild.Prefix = newPrefix
	guild.UpdatedAt = time.Now()
	return s.guildRepo.Update(ctx, guild)
}

// AddModeratorRole adds a role ID to the guild's moderator roles
func (s *GuildService) AddModeratorRole(ctx context.Context, guildID, roleID string) error {
	guild, err := s.guildRepo.FindByDiscordID(ctx, guildID)
	if err != nil {
		return err
	}

	// Check if role already exists
	for _, r := range guild.ModeratorRoles {
		if r == roleID {
			return nil
		}
	}

	guild.ModeratorRoles = append(guild.ModeratorRoles, roleID)
	guild.UpdatedAt = time.Now()
	return s.guildRepo.Update(ctx, guild)
} 