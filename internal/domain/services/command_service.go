package services

import (
	"context"
	"strings"
	"time"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// CommandService handles business logic for command operations
type CommandService struct {
	cmdRepo repositories.CommandRepository
}

// NewCommandService creates a new instance of CommandService
func NewCommandService(cmdRepo repositories.CommandRepository) *CommandService {
	return &CommandService{
		cmdRepo: cmdRepo,
	}
}

// CreateCommand creates a new custom command
func (s *CommandService) CreateCommand(ctx context.Context, guildID, name, description, response, createdBy string) (*entities.Command, error) {
	// Normalize command name
	name = strings.ToLower(strings.TrimSpace(name))

	cmd := &entities.Command{
		GuildID:     guildID,
		Name:        name,
		Description: description,
		Type:        entities.TextCommand,
		Response:    response,
		CreatedBy:   createdBy,
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.cmdRepo.Create(ctx, cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// GetCommandByName retrieves a command by its name within a guild
func (s *CommandService) GetCommandByName(ctx context.Context, guildID, name string) (*entities.Command, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	return s.cmdRepo.FindByName(ctx, guildID, name)
}

// ListGuildCommands retrieves all commands for a specific guild
func (s *CommandService) ListGuildCommands(ctx context.Context, guildID string) ([]*entities.Command, error) {
	return s.cmdRepo.ListByGuild(ctx, guildID)
}

// UpdateCommandResponse updates the response of a text command
func (s *CommandService) UpdateCommandResponse(ctx context.Context, cmdID, newResponse string) error {
	cmd, err := s.cmdRepo.FindByID(ctx, cmdID)
	if err != nil {
		return err
	}

	if cmd.Type != entities.TextCommand {
		return nil // Only update text commands
	}

	cmd.Response = newResponse
	cmd.UpdatedAt = time.Now()
	return s.cmdRepo.Update(ctx, cmd)
}

// ToggleCommand enables or disables a command
func (s *CommandService) ToggleCommand(ctx context.Context, cmdID string, enabled bool) error {
	cmd, err := s.cmdRepo.FindByID(ctx, cmdID)
	if err != nil {
		return err
	}

	cmd.Enabled = enabled
	cmd.UpdatedAt = time.Now()
	return s.cmdRepo.Update(ctx, cmd)
} 