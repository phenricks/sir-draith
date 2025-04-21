package commands

import (
	"fmt"
	"log"
	"sirdraith/internal/domain/repository"
	"sirdraith/internal/domain/services"
	"sirdraith/internal/infrastructure/mongodb/repositories"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
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
	Command  string
}

// Reply envia uma resposta simples para o canal
func (ctx *CommandContext) Reply(message string) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, message)
	return err
}

// CommandRegistry mantém o registro de todos os comandos disponíveis
type CommandRegistry struct {
	commands         map[string]*Command
	defaultPrefix    string
	configRepository repository.ConfigRepository
	characterService *services.CharacterService
	wizards          map[string]*CharacterWizard // Mapa de wizards ativos por userID
	session          *discordgo.Session
	db               *mongo.Database // Add MongoDB database field
}

// NewCommandRegistry cria um novo registro de comandos
func NewCommandRegistry(session *discordgo.Session, defaultPrefix string, configRepo repository.ConfigRepository, characterService *services.CharacterService, db *mongo.Database) *CommandRegistry {
	registry := &CommandRegistry{
		commands:         make(map[string]*Command),
		defaultPrefix:    defaultPrefix,
		configRepository: configRepo,
		characterService: characterService,
		wizards:          make(map[string]*CharacterWizard),
		session:          session,
		db:               db,
	}

	// Registrar comandos
	registry.RegisterCommands()

	return registry
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
		Command:  cmdName,
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

// HandleInteraction processa uma interação de componente
func (r *CommandRegistry) HandleInteraction(i *discordgo.InteractionCreate) error {
	// Verifica se é uma interação de componente
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	// Registra a interação para debug
	log.Printf("Componente acionado: %s por %s", i.MessageComponentData().CustomID, i.Member.User.Username)

	// Busca o wizard ativo para o usuário
	wizard, exists := r.wizards[i.Member.User.ID]
	if !exists {
		return fmt.Errorf("nenhum wizard ativo para este usuário")
	}

	// Processa a interação no wizard
	return wizard.HandleInteraction(i)
}

// RegisterCommands registra todos os comandos disponíveis
func (r *CommandRegistry) RegisterCommands() {
	// Registrar comandos de personagem
	characterCommands := NewCharacterCommands(r.characterService)
	characterCommands.Register(r)

	// Registrar comandos de perícia
	skillCommands := NewSkillCommands(r.characterService)
	skillCommands.Register(r)

	// Registrar comandos de deck
	deckService := services.NewDeckService(
		repositories.NewMongoDeckRepository(r.db),
		repositories.NewMongoCardRepository(r.db),
	)
	deckCommands := NewDeckCommands(deckService)
	deckCommands.Register(r)
}

// GetWizard retorna o wizard ativo para um usuário
func (r *CommandRegistry) GetWizard(userID string) *CharacterWizard {
	return r.wizards[userID]
}

// SetWizard define o wizard ativo para um usuário
func (r *CommandRegistry) SetWizard(userID string, wizard *CharacterWizard) {
	r.wizards[userID] = wizard
}

// RemoveWizard remove o wizard ativo de um usuário
func (r *CommandRegistry) RemoveWizard(userID string) {
	delete(r.wizards, userID)
}
