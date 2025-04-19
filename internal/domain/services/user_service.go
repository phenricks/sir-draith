package services

import (
	"context"
	"time"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user from Discord information
func (s *UserService) CreateUser(ctx context.Context, discordID, username, discriminator, avatarURL string) (*entities.User, error) {
	user := &entities.User{
		DiscordID:     discordID,
		Username:      username,
		Discriminator: discriminator,
		AvatarURL:     avatarURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByDiscordID retrieves a user by their Discord ID
func (s *UserService) GetUserByDiscordID(ctx context.Context, discordID string) (*entities.User, error) {
	return s.userRepo.FindByDiscordID(ctx, discordID)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
} 