package events

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GuildMemberHandler handles member-related events
type GuildMemberHandler struct{}

func NewGuildMemberHandler() *GuildMemberHandler {
	return &GuildMemberHandler{}
}

func (h *GuildMemberHandler) Handle(s *discordgo.Session, i interface{}) error {
	switch event := i.(type) {
	case *discordgo.GuildMemberAdd:
		return h.handleMemberAdd(s, event)
	case *discordgo.GuildMemberRemove:
		return h.handleMemberRemove(s, event)
	default:
		return fmt.Errorf("evento não suportado para GuildMemberHandler")
	}
}

func (h *GuildMemberHandler) handleMemberAdd(s *discordgo.Session, event *discordgo.GuildMemberAdd) error {
	log.Printf("Novo membro entrou no servidor %s: %s#%s",
		event.GuildID, event.User.Username, event.User.Discriminator)

	// Busca o canal de boas-vindas
	guild, err := s.Guild(event.GuildID)
	if err != nil {
		return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
	}

	// Tenta usar o canal de sistema, se disponível
	channelID := guild.SystemChannelID
	if channelID == "" {
		// Procura um canal com "welcome", "bem-vindo" ou similar no nome
		channels, err := s.GuildChannels(event.GuildID)
		if err != nil {
			return fmt.Errorf("erro ao buscar canais do servidor: %w", err)
		}

		for _, ch := range channels {
			if ch.Type == discordgo.ChannelTypeGuildText &&
				(containsAny(ch.Name, []string{"welcome", "bem-vindo", "boas-vindas"})) {
				channelID = ch.ID
				break
			}
		}
	}

	if channelID != "" {
		embed := &discordgo.MessageEmbed{
			Title: "🎉 Bem-vindo(a) ao Reino!",
			Description: fmt.Sprintf("Saudações, nobre %s! Que vossa jornada neste reino seja repleta de aventuras e glórias!\n\n"+
				"🛡️ **Dicas para começar:**\n"+
				"• Use `!help` para ver todos os comandos disponíveis\n"+
				"• Leia as regras do servidor para uma convivência harmoniosa\n"+
				"• Apresente-se aos outros membros do reino",
				event.User.Mention()),
			Color: 0x2ecc71, // Verde esmeralda
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: event.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎭 Membro",
					Value:  fmt.Sprintf("%s#%s", event.User.Username, event.User.Discriminator),
					Inline: true,
				},
				{
					Name:   "📅 Entrada",
					Value:  time.Now().Format("02/01/2006 15:04"),
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Membro #%d", guild.MemberCount),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, err = s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Printf("Erro ao enviar mensagem de boas-vindas: %v", err)
		}
	}

	return nil
}

func (h *GuildMemberHandler) handleMemberRemove(s *discordgo.Session, event *discordgo.GuildMemberRemove) error {
	log.Printf("Membro saiu do servidor %s: %s#%s",
		event.GuildID, event.User.Username, event.User.Discriminator)

	// Busca o canal de sistema
	guild, err := s.Guild(event.GuildID)
	if err != nil {
		return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
	}

	// Usa o mesmo canal das boas-vindas
	if guild.SystemChannelID != "" {
		embed := &discordgo.MessageEmbed{
			Title: "👋 Até logo!",
			Description: fmt.Sprintf("O nobre %s partiu de nosso reino. Que sua jornada seja próspera por onde quer que vá!",
				event.User.Username),
			Color: 0xe74c3c, // Vermelho
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: event.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎭 Membro",
					Value:  fmt.Sprintf("%s#%s", event.User.Username, event.User.Discriminator),
					Inline: true,
				},
				{
					Name:   "📅 Saída",
					Value:  time.Now().Format("02/01/2006 15:04"),
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Membros restantes: %d", guild.MemberCount-1),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		_, err = s.ChannelMessageSendEmbed(guild.SystemChannelID, embed)
		if err != nil {
			log.Printf("Erro ao enviar mensagem de despedida: %v", err)
		}
	}

	return nil
}

// containsAny verifica se uma string contém qualquer uma das substrings fornecidas
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

// contains verifica se uma string contém uma substring (case insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
