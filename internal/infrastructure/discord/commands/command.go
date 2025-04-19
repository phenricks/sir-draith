package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sirdraith/internal/domain/repository"
)

// Command representa um comando do bot
type Command struct {
	Name        string   // Nome do comando
	Aliases     []string // Aliases alternativos para o comando
	Description string   // Descrição do comando
	Usage       string   // Como usar o comando
	Category    string   // Categoria do comando (ex: admin, geral, etc)
	Handler     CommandHandlerFunc
}

// CommandHandlerFunc é a função que executa a lógica do comando
type CommandHandlerFunc func(ctx *CommandContext) error

// CommandContext contém o contexto de execução do comando
type CommandContext struct {
	Session  *discordgo.Session
	Message  *discordgo.MessageCreate
	Args     []string
	Registry *CommandRegistry
}

// CommandRegistry mantém o registro de todos os comandos disponíveis
type CommandRegistry struct {
	commands         map[string]*Command
	defaultPrefix    string
	configRepository repository.ConfigRepository
}

// NewCommandRegistry cria um novo registro de comandos
func NewCommandRegistry(defaultPrefix string, configRepo repository.ConfigRepository) *CommandRegistry {
	return &CommandRegistry{
		commands:         make(map[string]*Command),
		defaultPrefix:    defaultPrefix,
		configRepository: configRepo,
	}
}

// RegisterCommand registra um novo comando
func (r *CommandRegistry) RegisterCommand(cmd *Command) {
	r.commands[cmd.Name] = cmd
	
	// Registra aliases
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

// GetCommand retorna um comando pelo nome ou alias
func (r *CommandRegistry) GetCommand(name string) *Command {
	return r.commands[name]
}

// GetPrefix retorna o prefixo atual dos comandos para um servidor específico
func (r *CommandRegistry) GetPrefix() string {
	return r.defaultPrefix
}

// GetGuildPrefix retorna o prefixo específico de um servidor
func (r *CommandRegistry) GetGuildPrefix(guildID string) (string, error) {
	if r.configRepository == nil {
		return r.defaultPrefix, nil
	}

	config, err := r.configRepository.GetGuildConfig(guildID)
	if err != nil {
		return r.defaultPrefix, fmt.Errorf("erro ao obter prefixo do servidor: %w", err)
	}

	return config.Prefix, nil
}

// SetGuildPrefix define o prefixo específico de um servidor
func (r *CommandRegistry) SetGuildPrefix(guildID string, newPrefix string) error {
	if r.configRepository == nil {
		return fmt.Errorf("repositório de configurações não inicializado")
	}

	return r.configRepository.UpdateGuildPrefix(guildID, newPrefix)
}

// GetCommands retorna todos os comandos registrados
func (r *CommandRegistry) GetCommands() map[string]*Command {
	return r.commands
}

// HandleMessage processa uma mensagem e executa o comando apropriado
func (r *CommandRegistry) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := m.Content
	
	// Obtém o prefixo específico do servidor
	prefix := r.defaultPrefix
	if m.GuildID != "" {
		var err error
		prefix, err = r.GetGuildPrefix(m.GuildID)
		if err != nil {
			return err
		}
	}

	// Verifica se a mensagem começa com o prefixo
	if len(content) <= len(prefix) || content[:len(prefix)] != prefix {
		return nil
	}

	// Remove o prefixo e divide a mensagem em argumentos
	args := ParseArgs(content[len(prefix):])
	if len(args) == 0 {
		return nil
	}

	// Obtém o comando e os argumentos
	cmdName := args[0]
	cmd := r.GetCommand(cmdName)
	if cmd == nil {
		return nil
	}

	// Cria o contexto do comando
	ctx := &CommandContext{
		Session:  s,
		Message:  m,
		Args:     args[1:],
		Registry: r,
	}

	// Executa o comando
	return cmd.Handler(ctx)
}

// ParseArgs divide uma string em argumentos, respeitando aspas
func ParseArgs(content string) []string {
	var args []string
	var current string
	inQuotes := false

	for _, char := range content {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if !inQuotes {
				if current != "" {
					args = append(args, current)
					current = ""
				}
			} else {
				current += string(char)
			}
		default:
			current += string(char)
		}
	}

	if current != "" {
		args = append(args, current)
	}

	return args
} 