package events

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ReadyHandler handles the Ready event
type ReadyHandler struct{}

func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{}
}

func (h *ReadyHandler) Handle(s *discordgo.Session, i interface{}) error {
	ready, ok := i.(*discordgo.Ready)
	if !ok {
		return fmt.Errorf("evento inválido para Ready handler")
	}

	log.Printf("Bot está pronto! Conectado como: %s#%s\n", ready.User.Username, ready.User.Discriminator)
	log.Printf("Presente em %d servidores\n", len(ready.Guilds))

	// Set bot's status
	err := s.UpdateGameStatus(0, "Protegendo o Reino | !help")
	if err != nil {
		log.Printf("Erro ao atualizar status: %v", err)
	}

	return nil
}

// GuildCreateHandler handles the GuildCreate event
type GuildCreateHandler struct{}

func NewGuildCreateHandler() *GuildCreateHandler {
	return &GuildCreateHandler{}
}

func (h *GuildCreateHandler) Handle(s *discordgo.Session, i interface{}) error {
	guild, ok := i.(*discordgo.GuildCreate)
	if !ok {
		return fmt.Errorf("evento inválido para GuildCreate handler")
	}

	log.Printf("Bot adicionado ao servidor: %s (ID: %s)\n", guild.Name, guild.ID)

	// Send welcome message to the system channel if available
	if guild.SystemChannelID != "" {
		embed := &discordgo.MessageEmbed{
			Title: "🏰 Saudações, nobre servidor!",
			Description: "Eu sou Sir Draith, vosso fiel servo e guardião deste reino digital. " +
				"Estou aqui para auxiliar em vossa jornada com meus poderes místicos e conhecimentos ancestrais.\n\n" +
				"Use `!help` para descobrir todos os meus comandos e habilidades.",
			Color:     0x00ff00,
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Sir Draith - Bot Medieval",
			},
		}

		_, err := s.ChannelMessageSendEmbed(guild.SystemChannelID, embed)
		if err != nil {
			log.Printf("Erro ao enviar mensagem de boas-vindas: %v", err)
		}
	}

	return nil
}

// GuildDeleteHandler handles the GuildDelete event
type GuildDeleteHandler struct{}

func NewGuildDeleteHandler() *GuildDeleteHandler {
	return &GuildDeleteHandler{}
}

func (h *GuildDeleteHandler) Handle(s *discordgo.Session, i interface{}) error {
	guild, ok := i.(*discordgo.GuildDelete)
	if !ok {
		return fmt.Errorf("evento inválido para GuildDelete handler")
	}

	log.Printf("Bot removido do servidor: %s (ID: %s)\n", guild.Name, guild.ID)
	return nil
}

// InteractionCreateHandler handles the InteractionCreate event
type InteractionCreateHandler struct{}

func NewInteractionCreateHandler() *InteractionCreateHandler {
	return &InteractionCreateHandler{}
}

func (h *InteractionCreateHandler) Handle(s *discordgo.Session, i interface{}) error {
	interaction, ok := i.(*discordgo.InteractionCreate)
	if !ok {
		return fmt.Errorf("evento inválido para InteractionCreate handler")
	}

	// Handle different types of interactions
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		return h.handleApplicationCommand(s, interaction)
	case discordgo.InteractionMessageComponent:
		return h.handleMessageComponent(s, interaction)
	}

	return nil
}

func (h *InteractionCreateHandler) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	// Resposta padrão para comandos não implementados
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚧 Este comando ainda não está disponível. Use os comandos com prefixo por enquanto!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	// Responde à interação
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Erro ao responder ao comando %s: %v", data.Name, err)
		return err
	}

	return nil
}

func (h *InteractionCreateHandler) handleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.MessageComponentData()

	// Log do componente acionado
	log.Printf("Componente acionado: %s por %s#%s",
		data.CustomID, i.Member.User.Username, i.Member.User.Discriminator)

	// Resposta padrão para componentes não implementados
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚧 Esta interação ainda não está disponível!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	// Responde à interação
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Printf("Erro ao responder ao componente %s: %v", data.CustomID, err)
		return err
	}

	return nil
} 