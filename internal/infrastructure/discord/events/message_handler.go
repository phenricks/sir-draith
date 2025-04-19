package events

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// MessageHandler handles all message-related events in Discord
type MessageHandler struct {
	// Configurações do handler
	config struct {
		maxMentions int // Número máximo de menções permitidas em uma mensagem
		logChannel  string // ID do canal de logs (opcional)
	}
}

// NewMessageHandler creates a new instance of MessageHandler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		config: struct {
			maxMentions int
			logChannel  string
		}{
			maxMentions: 5, // Valor padrão
			logChannel:  "", // Será configurado pelo servidor
		},
	}
}

// Handle processes all message-related events
func (h *MessageHandler) Handle(s *discordgo.Session, i interface{}) error {
	switch event := i.(type) {
	case *discordgo.MessageCreate:
		return h.handleMessageCreate(s, event)
	case *discordgo.MessageDelete:
		return h.handleMessageDelete(s, event)
	case *discordgo.MessageUpdate:
		return h.handleMessageUpdate(s, event)
	default:
		return fmt.Errorf("evento não suportado para MessageHandler")
	}
}

// handleMessageCreate processes new messages
func (h *MessageHandler) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Ignora mensagens do próprio bot
	if m.Author.ID == s.State.User.ID {
		return nil
	}

	// Log da mensagem
	log.Printf("Nova mensagem de %s#%s no canal %s: %s",
		m.Author.Username, m.Author.Discriminator, m.ChannelID, m.Content)

	// Verifica menções excessivas (anti-spam)
	if len(m.Mentions) > h.config.maxMentions {
		// Tenta deletar a mensagem
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Printf("Erro ao deletar mensagem com menções excessivas: %v", err)
		}

		// Avisa o usuário
		warning := fmt.Sprintf("🛡️ %s, por favor evite mencionar muitos usuários de uma vez.", m.Author.Mention())
		_, err = s.ChannelMessageSend(m.ChannelID, warning)
		if err != nil {
			log.Printf("Erro ao enviar aviso sobre menções: %v", err)
		}

		return nil
	}

	// Registra a mensagem no canal de logs, se configurado
	if h.config.logChannel != "" {
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator),
				IconURL: m.Author.AvatarURL(""),
			},
			Description: m.Content,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Canal",
					Value:  fmt.Sprintf("<#%s>", m.ChannelID),
					Inline: true,
				},
			},
			Color:     0x3498db, // Azul
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID da Mensagem: %s", m.ID),
			},
		}

		// Adiciona anexos, se houver
		if len(m.Attachments) > 0 {
			attachments := make([]string, 0, len(m.Attachments))
			for _, a := range m.Attachments {
				attachments = append(attachments, a.URL)
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Anexos",
				Value:  strings.Join(attachments, "\n"),
				Inline: false,
			})
		}

		_, err := s.ChannelMessageSendEmbed(h.config.logChannel, embed)
		if err != nil {
			log.Printf("Erro ao registrar mensagem no canal de logs: %v", err)
		}
	}

	return nil
}

// handleMessageDelete processes deleted messages
func (h *MessageHandler) handleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) error {
	// Log da deleção
	log.Printf("Mensagem %s deletada no canal %s", m.ID, m.ChannelID)

	// Registra a deleção no canal de logs, se configurado
	if h.config.logChannel != "" {
		embed := &discordgo.MessageEmbed{
			Title:       "📝 Mensagem Deletada",
			Description: m.Content,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Canal",
					Value:  fmt.Sprintf("<#%s>", m.ChannelID),
					Inline: true,
				},
			},
			Color:     0xe74c3c, // Vermelho
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID da Mensagem: %s", m.ID),
			},
		}

		_, err := s.ChannelMessageSendEmbed(h.config.logChannel, embed)
		if err != nil {
			log.Printf("Erro ao registrar deleção no canal de logs: %v", err)
		}
	}

	return nil
}

// handleMessageUpdate processes edited messages
func (h *MessageHandler) handleMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) error {
	// Ignora atualizações vazias ou do próprio bot
	if m.Author != nil && m.Author.ID == s.State.User.ID {
		return nil
	}

	// Log da edição
	log.Printf("Mensagem %s editada no canal %s", m.ID, m.ChannelID)

	// Registra a edição no canal de logs, se configurado
	if h.config.logChannel != "" {
		embed := &discordgo.MessageEmbed{
			Title:       "📝 Mensagem Editada",
			Description: m.Content,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Canal",
					Value:  fmt.Sprintf("<#%s>", m.ChannelID),
					Inline: true,
				},
			},
			Color:     0xf1c40f, // Amarelo
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID da Mensagem: %s", m.ID),
			},
		}

		_, err := s.ChannelMessageSendEmbed(h.config.logChannel, embed)
		if err != nil {
			log.Printf("Erro ao registrar edição no canal de logs: %v", err)
		}
	}

	return nil
}

// SetLogChannel configura o canal de logs
func (h *MessageHandler) SetLogChannel(channelID string) {
	h.config.logChannel = channelID
}

// SetMaxMentions configura o número máximo de menções permitidas
func (h *MessageHandler) SetMaxMentions(max int) {
	h.config.maxMentions = max
} 