package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"

	"sirdraith/internal/domain/repository"
	"sirdraith/internal/domain/services"
	"sirdraith/internal/infrastructure/discord/commands"
	"sirdraith/internal/infrastructure/discord/events"
	"sirdraith/internal/infrastructure/mongodb"
)

// Client representa o cliente Discord do bot
type Client struct {
	session          *discordgo.Session
	configRepo       repository.ConfigRepository
	commandRegistry  *commands.CommandRegistry
	eventManager     *events.EventManager
	characterService *services.CharacterService
}

// NewClient cria uma nova instância do cliente Discord
func NewClient(token string, db *mongo.Database) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token do Discord não pode estar vazio")
	}

	if db == nil {
		return nil, fmt.Errorf("banco de dados não pode estar vazio")
	}

	// Cria a sessão do Discord
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar sessão do Discord: %w", err)
	}

	// Inicializa repositórios
	configRepo := mongodb.NewConfigRepository(db)
	characterRepo := mongodb.NewCharacterRepository(db)

	// Inicializa serviços
	characterService := services.NewCharacterService(characterRepo)

	// Inicializa o registro de comandos e gerenciador de eventos
	registry := commands.NewCommandRegistry(session, "!", configRepo, characterService, db)
	eventManager := events.NewEventManager(session, configRepo, registry)

	client := &Client{
		session:          session,
		configRepo:       configRepo,
		commandRegistry:  registry,
		eventManager:     eventManager,
		characterService: characterService,
	}

	// Registra os comandos
	client.registerCommands()

	return client, nil
}

// registerCommands registra todos os comandos disponíveis
func (c *Client) registerCommands() {
	// Registra comandos básicos
	for _, cmd := range commands.BasicCommands() {
		c.commandRegistry.RegisterCommand(cmd)
	}

	// Registra comandos administrativos
	for _, cmd := range commands.AdminCommands() {
		c.commandRegistry.RegisterCommand(cmd)
	}

	// Registra comandos utilitários
	for _, cmd := range commands.UtilityCommands() {
		c.commandRegistry.RegisterCommand(cmd)
	}

	// Registra comandos de personagem
	characterCommands := commands.NewCharacterCommands(c.characterService)
	characterCommands.Register(c.commandRegistry)
}

// Connect estabelece a conexão com o Discord
func (c *Client) Connect() error {
	// Registra os handlers padrão
	c.eventManager.RegisterDefaultHandlers()

	// Registra os handlers de eventos
	c.session.AddHandler(events.WrapHandler(events.EventReady, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventGuildCreate, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventGuildDelete, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventGuildMemberAdd, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventGuildMemberRem, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventMessageCreate, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventMessageDelete, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventMessageUpdate, c.eventManager))
	c.session.AddHandler(events.WrapHandler(events.EventInteractionCreate, c.eventManager))

	// Registra o handler de mensagens para comandos
	c.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if err := c.commandRegistry.HandleMessage(s, m); err != nil {
			log.Printf("Erro ao processar comando: %v\n", err)
		}
	})

	// Registra o handler de interações
	c.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if err := c.commandRegistry.HandleInteraction(i); err != nil {
			log.Printf("Erro ao processar interação: %v\n", err)
		}
	})

	// Define intents necessários
	c.session.Identify.Intents = discordgo.IntentsAll

	// Abre a conexão com o Discord
	if err := c.session.Open(); err != nil {
		return fmt.Errorf("erro ao abrir conexão com Discord: %w", err)
	}

	log.Println("Bot conectado ao Discord com sucesso!")
	return nil
}

// Disconnect desconecta o bot do Discord
func (c *Client) Disconnect() error {
	if err := c.session.Close(); err != nil {
		return fmt.Errorf("erro ao desconectar do Discord: %w", err)
	}
	return nil
}
