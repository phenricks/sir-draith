package events

import (
	"fmt"
	"log"
	"time"

	"sirdraith/internal/domain/repository"

	"github.com/bwmarrin/discordgo"
)

// ReadyHandler handles the Ready event
type ReadyHandler struct {
	session *discordgo.Session
}

func NewReadyHandler(s *discordgo.Session) *ReadyHandler {
	return &ReadyHandler{session: s}
}

func (h *ReadyHandler) Handle(s *discordgo.Session, i interface{}) error {
	log.Printf("Bot está pronto! Conectado como %s", s.State.User.Username)
	return h.session.UpdateGameStatus(0, "Protegendo o Reino")
}

// GuildCreateHandler handles the GuildCreate event
type GuildCreateHandler struct {
	session *discordgo.Session
}

func NewGuildCreateHandler(s *discordgo.Session) *GuildCreateHandler {
	return &GuildCreateHandler{session: s}
}

func (h *GuildCreateHandler) Handle(s *discordgo.Session, i interface{}) error {
	guild, ok := i.(*discordgo.GuildCreate)
	if !ok {
		return fmt.Errorf("evento inválido para GuildCreate handler")
	}

	log.Printf("Bot adicionado ao servidor: %s", guild.Name)

	// Envia mensagem de boas-vindas no canal do sistema
	if guild.SystemChannelID != "" {
		embed := &discordgo.MessageEmbed{
			Title: "🗡️ Sir Draith - O Protetor do Reino",
			Description: fmt.Sprintf("Saudações, nobres membros do Reino **%s**!\n\n"+
				"Eu sou **Sir Draith**, um cavaleiro místico treinado nas artes da guerra e da magia, "+
				"designado para servir e proteger vosso reino digital. "+
				"Estou aqui para auxiliar em vossas aventuras e manter a ordem em vosso servidor.\n\n"+
				"**🛡️ Meus Serviços:**\n"+
				"• Sistema de RPG com criação de personagens\n"+
				"• Batalhas épicas com sistema de cartas\n"+
				"• Gerenciamento de campanhas e narrativas\n"+
				"• Proteção e moderação do servidor\n\n"+
				"**📜 Comandos Principais:**\n"+
				"• `!ajuda` - Lista todos os comandos disponíveis\n"+
				"• `!criar` - Inicia a criação de um personagem\n"+
				"• `!info` - Mostra informações detalhadas sobre mim\n"+
				"• `!config` - Configura as opções do servidor (Admin)\n\n"+
				"Que vossa jornada neste reino seja repleta de aventuras e glórias! 🏰", guild.Name),
			Color: 0x2b2d31, // Cor temática medieval (cinza escuro)
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: s.State.User.AvatarURL(""),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Sir Draith - Bot Medieval • Use !ajuda para mais informações",
			},
			Timestamp: time.Now().Format(time.RFC3339),
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

	log.Printf("Bot removido do servidor: %s", guild.ID)
	return nil
}

// GuildMemberAddHandler handles the GuildMemberAdd event
type GuildMemberAddHandler struct {
	guildConfigRepo repository.ConfigRepository
}

func NewGuildMemberAddHandler(repo repository.ConfigRepository) *GuildMemberAddHandler {
	return &GuildMemberAddHandler{
		guildConfigRepo: repo,
	}
}

func (h *GuildMemberAddHandler) Handle(s *discordgo.Session, i interface{}) error {
	member, ok := i.(*discordgo.GuildMemberAdd)
	if !ok {
		return fmt.Errorf("evento inválido para GuildMemberAdd handler")
	}

	log.Printf("Novo membro entrou no servidor %s: %s#%s",
		member.GuildID, member.User.Username, member.User.Discriminator)

	// Busca a configuração do servidor
	config, err := h.guildConfigRepo.GetGuildConfig(member.GuildID)
	if err != nil {
		log.Printf("Erro ao buscar configuração do servidor: %v", err)
		return err
	}

	// Determina o canal de boas-vindas
	channelID := config.WelcomeChannel
	if channelID == "" {
		// Se não houver canal configurado, tenta encontrar um apropriado
		guild, err := s.Guild(member.GuildID)
		if err != nil {
			return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
		}

		// Tenta usar o canal de sistema primeiro
		channelID = guild.SystemChannelID
		if channelID == "" {
			// Procura um canal com "welcome", "bem-vindo" ou similar no nome
			channels, err := s.GuildChannels(member.GuildID)
			if err != nil {
				return fmt.Errorf("erro ao buscar canais do servidor: %w", err)
			}

			for _, ch := range channels {
				if ch.Type == discordgo.ChannelTypeGuildText &&
					(stringContainsAny(ch.Name, []string{"welcome", "bem-vindo", "boas-vindas"})) {
					channelID = ch.ID
					break
				}
			}
		}
	}

	if channelID != "" {
		guild, err := s.Guild(member.GuildID)
		if err != nil {
			return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
		}

		embed := &discordgo.MessageEmbed{
			Title: "🎉 Bem-vindo(a) ao Reino!",
			Description: fmt.Sprintf("Saudações, nobre %s! Que vossa jornada neste reino seja repleta de aventuras e glórias!\n\n"+
				"🛡️ **Dicas para começar:**\n"+
				"• Use `!help` para ver todos os comandos disponíveis\n"+
				"• Leia as regras do servidor para uma convivência harmoniosa\n"+
				"• Apresente-se aos outros membros do reino",
				member.User.Mention()),
			Color: 0x2ecc71, // Verde esmeralda
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: member.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎭 Membro",
					Value:  fmt.Sprintf("%s#%s", member.User.Username, member.User.Discriminator),
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
			return err
		}
	}

	return nil
}

// GuildMemberRemoveHandler handles the GuildMemberRemove event
type GuildMemberRemoveHandler struct {
	guildConfigRepo repository.ConfigRepository
}

func NewGuildMemberRemoveHandler(repo repository.ConfigRepository) *GuildMemberRemoveHandler {
	return &GuildMemberRemoveHandler{
		guildConfigRepo: repo,
	}
}

func (h *GuildMemberRemoveHandler) Handle(s *discordgo.Session, i interface{}) error {
	member, ok := i.(*discordgo.GuildMemberRemove)
	if !ok {
		return fmt.Errorf("evento inválido para GuildMemberRemove handler")
	}

	log.Printf("Membro saiu do servidor %s: %s#%s",
		member.GuildID, member.User.Username, member.User.Discriminator)

	// Busca a configuração do servidor
	config, err := h.guildConfigRepo.GetGuildConfig(member.GuildID)
	if err != nil {
		log.Printf("Erro ao buscar configuração do servidor: %v", err)
		return err
	}

	// Determina o canal de despedida
	channelID := config.GoodbyeChannel
	if channelID == "" {
		// Se não houver canal configurado, usa o canal do sistema
		guild, err := s.Guild(member.GuildID)
		if err != nil {
			return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
		}
		channelID = guild.SystemChannelID
	}

	if channelID != "" {
		guild, err := s.Guild(member.GuildID)
		if err != nil {
			return fmt.Errorf("erro ao buscar informações do servidor: %w", err)
		}

		embed := &discordgo.MessageEmbed{
			Title: "👋 Até logo, nobre aventureiro!",
			Description: fmt.Sprintf("O nobre %s partiu de nosso reino. Que sua jornada seja próspera por onde quer que vá!",
				member.User.Username),
			Color: 0xe74c3c, // Vermelho
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: member.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🎭 Membro",
					Value:  fmt.Sprintf("%s#%s", member.User.Username, member.User.Discriminator),
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

		_, err = s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Printf("Erro ao enviar mensagem de despedida: %v", err)
			return err
		}
	}

	return nil
}
