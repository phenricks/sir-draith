package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// GuildRepository defines the interface for guild persistence operations
type GuildRepository interface {
	// Create stores a new guild in the database
	Create(ctx context.Context, guild *entities.Guild) error

	// FindByID retrieves a guild by its ID
	FindByID(ctx context.Context, id string) (*entities.Guild, error)

	// FindByDiscordID retrieves a guild by its Discord ID
	FindByDiscordID(ctx context.Context, discordID string) (*entities.Guild, error)

	// Update updates an existing guild in the database
	Update(ctx context.Context, guild *entities.Guild) error

	// Delete removes a guild from the database
	Delete(ctx context.Context, id string) error

	// ListAll retrieves all guilds
	ListAll(ctx context.Context) ([]*entities.Guild, error)
} 