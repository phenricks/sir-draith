package entities

import "time"

// CommandType represents the type of command
type CommandType string

const (
	// TextCommand represents a simple text response command
	TextCommand CommandType = "text"
	// CustomCommand represents a command with custom logic
	CustomCommand CommandType = "custom"
)

// Command represents a custom bot command
type Command struct {
	ID          string      `bson:"_id"`
	GuildID     string      `bson:"guild_id"`
	Name        string      `bson:"name"`
	Description string      `bson:"description"`
	Type        CommandType `bson:"type"`
	Response    string      `bson:"response"`
	CreatedBy   string      `bson:"created_by"`
	Enabled     bool        `bson:"enabled"`
	CreatedAt   time.Time   `bson:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at"`
} 