package repository

import "sirdraith/internal/domain/model"

// ConfigRepository define as operações para gerenciar configurações de servidores
type ConfigRepository interface {
	// GetGuildConfig retorna a configuração de um servidor específico
	GetGuildConfig(guildID string) (*model.GuildConfig, error)

	// UpdateGuildPrefix atualiza o prefixo de comandos de um servidor
	UpdateGuildPrefix(guildID string, newPrefix string) error

	// EnsureGuildConfig garante que existe uma configuração para o servidor
	EnsureGuildConfig(guildID string) (*model.GuildConfig, error)

	// UpdateGuildConfig atualiza a configuração completa de um servidor
	UpdateGuildConfig(guildID string, config *model.GuildConfig) error
}
