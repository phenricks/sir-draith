package repositories

import (
	"context"

	"sirdraith/internal/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create stores a new user in the database
	Create(ctx context.Context, user *entities.User) error

	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id string) (*entities.User, error)

	// FindByDiscordID retrieves a user by their Discord ID
	FindByDiscordID(ctx context.Context, discordID string) (*entities.User, error)

	// Update updates an existing user in the database
	Update(ctx context.Context, user *entities.User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id string) error
} 