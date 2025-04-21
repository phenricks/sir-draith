package model

import "time"

// GuildConfig representa as configurações específicas de um servidor
type GuildConfig struct {
	ID             string    `bson:"_id"`             // ID do servidor
	Prefix         string    `bson:"prefix"`          // Prefixo de comandos personalizado
	WelcomeChannel string    `bson:"welcome_channel"` // Canal para mensagens de boas-vindas
	GoodbyeChannel string    `bson:"goodbye_channel"` // Canal para mensagens de despedida
	CreatedAt      time.Time `bson:"created_at"`      // Data de criação
	UpdatedAt      time.Time `bson:"updated_at"`      // Data da última atualização
}

// NewGuildConfig cria uma nova configuração de servidor com valores padrão
func NewGuildConfig(guildID string) *GuildConfig {
	now := time.Now()
	return &GuildConfig{
		ID:             guildID,
		Prefix:         "!", // Prefixo padrão
		WelcomeChannel: "",  // Canal de boas-vindas vazio por padrão
		GoodbyeChannel: "",  // Canal de despedida vazio por padrão
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
