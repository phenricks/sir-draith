package model

import "time"

// GuildConfig representa as configurações específicas de um servidor
type GuildConfig struct {
	ID        string    `bson:"_id"`        // ID do servidor
	Prefix    string    `bson:"prefix"`     // Prefixo de comandos personalizado
	CreatedAt time.Time `bson:"created_at"` // Data de criação
	UpdatedAt time.Time `bson:"updated_at"` // Data da última atualização
}

// NewGuildConfig cria uma nova configuração de servidor com valores padrão
func NewGuildConfig(guildID string) *GuildConfig {
	now := time.Now()
	return &GuildConfig{
		ID:        guildID,
		Prefix:    "!", // Prefixo padrão
		CreatedAt: now,
		UpdatedAt: now,
	}
}
