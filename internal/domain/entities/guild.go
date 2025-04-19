package entities

import "time"

// Guild represents a Discord server in our system
type Guild struct {
	ID              string    `bson:"_id"`
	DiscordID       string    `bson:"discord_id"`
	Name            string    `bson:"name"`
	Prefix          string    `bson:"prefix"`
	WelcomeChannel  string    `bson:"welcome_channel"`
	ModeratorRoles  []string  `bson:"moderator_roles"`
	EnabledFeatures []string  `bson:"enabled_features"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
} 