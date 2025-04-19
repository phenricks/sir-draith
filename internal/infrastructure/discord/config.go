package discord

// Config contém as configurações necessárias para o cliente Discord
type Config struct {
	Token     string
	Prefix    string
	GuildID   string // ID do servidor principal (opcional)
}

// NewConfig cria uma nova configuração com valores padrão
func NewConfig(token string) *Config {
	return &Config{
		Token:  token,
		Prefix: "!", // Prefixo padrão para comandos
	}
} 