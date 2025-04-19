package entities

import "time"

// User represents a Discord user in our system
type User struct {
	ID            string    `bson:"_id"`
	DiscordID     string    `bson:"discord_id"`
	Username      string    `bson:"username"`
	Discriminator string    `bson:"discriminator"`
	AvatarURL     string    `bson:"avatar_url"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
} 