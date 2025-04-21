package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// InteractionCreateHandler processa eventos de interação
type InteractionCreateHandler struct{}

// NewInteractionCreateHandler cria um novo handler de interações
func NewInteractionCreateHandler() *InteractionCreateHandler {
	return &InteractionCreateHandler{}
}

// Handle processa um evento de interação
func (h *InteractionCreateHandler) Handle(s *discordgo.Session, i interface{}) error {
	interaction, ok := i.(*discordgo.InteractionCreate)
	if !ok {
		return nil
	}

	// Registra a interação para debug
	if interaction.Member != nil && interaction.Member.User != nil {
		log.Printf("Componente acionado: %s por %s", interaction.MessageComponentData().CustomID, interaction.Member.User.Username)
	}

	return nil
}
