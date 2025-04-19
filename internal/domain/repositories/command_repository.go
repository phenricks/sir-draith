package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// CommandRepository defines the interface for command persistence operations
type CommandRepository interface {
	// Create stores a new command in the database
	Create(ctx context.Context, command *entities.Command) error

	// FindByID retrieves a command by its ID
	FindByID(ctx context.Context, id string) (*entities.Command, error)

	// FindByName retrieves a command by its name within a guild
	FindByName(ctx context.Context, guildID, name string) (*entities.Command, error)

	// ListByGuild retrieves all commands for a specific guild
	ListByGuild(ctx context.Context, guildID string) ([]*entities.Command, error)

	// Update updates an existing command in the database
	Update(ctx context.Context, command *entities.Command) error

	// Delete removes a command from the database
	Delete(ctx context.Context, id string) error
} 