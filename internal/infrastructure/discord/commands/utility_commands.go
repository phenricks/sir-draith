package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// UtilityCommands retorna os comandos utilitários do bot
func UtilityCommands() []*Command {
	return []*Command{
		ServerCommand(),
		UserCommand(),
		RoleCommand(),
		AvatarCommand(),
		SystemChannelCommand(),
	}
}

// ServerCommand cria o comando para mostrar informações do servidor
func ServerCommand() *Command {
	return &Command{
		Name:        "server",
		Aliases:     []string{"servidor", "guild"},
		Description: "Mostra informações sobre o servidor",
		Usage:       "server",
		Category:    "Utilidades",
		Handler: func(ctx *CommandContext) error {
			guild, err := ctx.Session.Guild(ctx.Message.GuildID)
			if err != nil {
				return fmt.Errorf("erro ao obter informações do servidor: %w", err)
			}

			// Obtém o dono do servidor
			owner, err := ctx.Session.User(guild.OwnerID)
			if err != nil {
				return fmt.Errorf("erro ao obter informações do dono: %w", err)
			}

			// Formata a data de criação
			createdAt, err := discordgo.SnowflakeTimestamp(guild.ID)
			if err != nil {
				return fmt.Errorf("erro ao obter data de criação: %w", err)
			}

			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("ℹ️ Informações do Servidor: %s", guild.Name),
				Color: 0x00ff00,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: guild.IconURL(""),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "👑 Dono",
						Value:  owner.Username,
						Inline: true,
					},
					{
						Name:   "👥 Membros",
						Value:  fmt.Sprintf("%d", guild.MemberCount),
						Inline: true,
					},
					{
						Name:   "📅 Criado em",
						Value:  createdAt.Format("02/01/2006 15:04"),
						Inline: true,
					},
					{
						Name:   "🌍 Região",
						Value:  string(guild.PreferredLocale),
						Inline: true,
					},
					{
						Name:   "💬 Canais",
						Value:  fmt.Sprintf("%d", len(guild.Channels)),
						Inline: true,
					},
					{
						Name:   "🎭 Cargos",
						Value:  fmt.Sprintf("%d", len(guild.Roles)),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %s", guild.ID),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// UserCommand cria o comando para mostrar informações de um usuário
func UserCommand() *Command {
	return &Command{
		Name:        "user",
		Aliases:     []string{"usuario", "userinfo"},
		Description: "Mostra informações sobre um usuário",
		Usage:       "user [@usuário]",
		Category:    "Utilidades",
		Handler: func(ctx *CommandContext) error {
			var user *discordgo.User
			var member *discordgo.Member

			// Se não mencionar ninguém, usa o autor da mensagem
			if len(ctx.Message.Mentions) == 0 {
				user = ctx.Message.Author
				var err error
				member, err = ctx.Session.GuildMember(ctx.Message.GuildID, user.ID)
				if err != nil {
					return fmt.Errorf("erro ao obter informações do membro: %w", err)
				}
			} else {
				user = ctx.Message.Mentions[0]
				var err error
				member, err = ctx.Session.GuildMember(ctx.Message.GuildID, user.ID)
				if err != nil {
					return fmt.Errorf("erro ao obter informações do membro: %w", err)
				}
			}

			// Formata a data de entrada no servidor
			joinedAt := member.JoinedAt

			// Formata a data de criação da conta
			createdAt, err := discordgo.SnowflakeTimestamp(user.ID)
			if err != nil {
				return fmt.Errorf("erro ao obter data de criação: %w", err)
			}

			// Lista os cargos do usuário
			var roles []string
			for _, roleID := range member.Roles {
				role, err := ctx.Session.State.Role(ctx.Message.GuildID, roleID)
				if err != nil {
					continue
				}
				roles = append(roles, role.Name)
			}

			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("👤 Informações do Usuário: %s", user.Username),
				Color: 0x00ff00,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: user.AvatarURL(""),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🏷️ Tag",
						Value:  fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
						Inline: true,
					},
					{
						Name:   "📅 Conta criada",
						Value:  createdAt.Format("02/01/2006 15:04"),
						Inline: true,
					},
					{
						Name:   "📥 Entrou em",
						Value:  joinedAt.Format("02/01/2006 15:04"),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %s", user.ID),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if len(roles) > 0 {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "🎭 Cargos",
					Value:  strings.Join(roles, ", "),
					Inline: false,
				})
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// RoleCommand cria o comando para mostrar informações de um cargo
func RoleCommand() *Command {
	return &Command{
		Name:        "role",
		Aliases:     []string{"cargo"},
		Description: "Mostra informações sobre um cargo",
		Usage:       "role <nome do cargo>",
		Category:    "Utilidades",
		Handler: func(ctx *CommandContext) error {
			if len(ctx.Args) == 0 {
				return fmt.Errorf("especifique o nome do cargo")
			}

			// Procura o cargo pelo nome
			roleName := strings.Join(ctx.Args, " ")
			var role *discordgo.Role

			guild, err := ctx.Session.Guild(ctx.Message.GuildID)
			if err != nil {
				return fmt.Errorf("erro ao obter informações do servidor: %w", err)
			}

			for _, r := range guild.Roles {
				if strings.EqualFold(r.Name, roleName) {
					role = r
					break
				}
			}

			if role == nil {
				return fmt.Errorf("cargo não encontrado")
			}

			// Conta quantos membros têm o cargo
			memberCount := 0
			members, err := ctx.Session.GuildMembers(ctx.Message.GuildID, "", 1000)
			if err != nil {
				return fmt.Errorf("erro ao obter membros: %w", err)
			}

			for _, member := range members {
				for _, roleID := range member.Roles {
					if roleID == role.ID {
						memberCount++
						break
					}
				}
			}

			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("🎭 Informações do Cargo: %s", role.Name),
				Color:       role.Color,
				Description: fmt.Sprintf("**Membros:** %d", memberCount),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🔢 Posição",
						Value:  fmt.Sprintf("%d", role.Position),
						Inline: true,
					},
					{
						Name:   "🎨 Cor",
						Value:  fmt.Sprintf("#%06x", role.Color),
						Inline: true,
					},
					{
						Name:   "🔐 Mencionável",
						Value:  fmt.Sprintf("%v", role.Mentionable),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %s", role.ID),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// AvatarCommand cria o comando para mostrar o avatar de um usuário
func AvatarCommand() *Command {
	return &Command{
		Name:        "avatar",
		Aliases:     []string{"foto", "pic"},
		Description: "Mostra o avatar de um usuário",
		Usage:       "avatar [@usuário]",
		Category:    "Utilidades",
		Handler: func(ctx *CommandContext) error {
			var user *discordgo.User

			// Se não mencionar ninguém, usa o autor da mensagem
			if len(ctx.Message.Mentions) == 0 {
				user = ctx.Message.Author
			} else {
				user = ctx.Message.Mentions[0]
			}

			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("🖼️ Avatar de %s", user.Username),
				Color: 0x00ff00,
				Image: &discordgo.MessageEmbedImage{
					URL: user.AvatarURL("2048"),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Solicitado por %s", ctx.Message.Author.Username),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// SystemChannelCommand cria o comando para gerenciar o canal de sistema
func SystemChannelCommand() *Command {
	return &Command{
		Name:        "systemchannel",
		Aliases:     []string{"canal-sistema", "canal"},
		Description: "Mostra ou configura o canal de sistema do servidor",
		Usage:       "systemchannel [#canal]",
		Category:    "Utilidades",
		Handler: func(ctx *CommandContext) error {
			// Verifica se o usuário tem permissão
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionAdministrator == 0 {
				return fmt.Errorf("você precisa ter permissão de administrador para usar este comando")
			}

			guild, err := ctx.Session.Guild(ctx.Message.GuildID)
			if err != nil {
				return fmt.Errorf("erro ao obter informações do servidor: %w", err)
			}

			// Se não houver argumentos, mostra o canal atual
			if len(ctx.Args) == 0 {
				var channelName string
				if guild.SystemChannelID != "" {
					channel, err := ctx.Session.Channel(guild.SystemChannelID)
					if err != nil {
						return fmt.Errorf("erro ao obter informações do canal: %w", err)
					}
					channelName = channel.Name
				} else {
					channelName = "Nenhum canal configurado"
				}

				embed := &discordgo.MessageEmbed{
					Title:       "📢 Canal do Sistema",
					Description: "Este é o canal onde enviarei mensagens de sistema, como boas-vindas e despedidas.",
					Color:       0x00ff00,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Canal Atual",
							Value:  fmt.Sprintf("<#%s> (%s)", guild.SystemChannelID, channelName),
							Inline: false,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Use !systemchannel #canal para alterar",
					},
				}

				_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
				return err
			}

			// Extrai o ID do canal da menção (formato: <#ID>)
			channelMention := ctx.Args[0]
			if len(channelMention) < 4 || !strings.HasPrefix(channelMention, "<#") || !strings.HasSuffix(channelMention, ">") {
				return fmt.Errorf("por favor, mencione um canal válido usando #nome-do-canal")
			}

			// Remove os caracteres <#> para obter o ID
			newChannelID := channelMention[2 : len(channelMention)-1]

			// Verifica se o canal existe
			newChannel, err := ctx.Session.Channel(newChannelID)
			if err != nil {
				return fmt.Errorf("canal inválido ou não encontrado: %w", err)
			}

			// Verifica se o canal está no mesmo servidor
			if newChannel.GuildID != ctx.Message.GuildID {
				return fmt.Errorf("o canal precisa estar neste servidor")
			}

			// Atualiza o canal do sistema
			guildParams := &discordgo.GuildParams{
				SystemChannelID: newChannelID,
			}

			_, err = ctx.Session.GuildEdit(ctx.Message.GuildID, guildParams)
			if err != nil {
				return fmt.Errorf("erro ao atualizar canal do sistema: %w", err)
			}

			embed := &discordgo.MessageEmbed{
				Title:       "✅ Canal do Sistema Atualizado",
				Description: fmt.Sprintf("O canal do sistema foi alterado para <#%s>", newChannelID),
				Color:       0x00ff00,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Canal",
						Value:  fmt.Sprintf("<#%s> (%s)", newChannelID, newChannel.Name),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Todas as mensagens de sistema serão enviadas neste canal",
				},
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}
