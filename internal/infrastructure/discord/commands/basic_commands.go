package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var startTime = time.Now()

// BasicCommands retorna os comandos básicos do bot
func BasicCommands() []*Command {
	return []*Command{
		PingCommand(),
		HelpCommand(),
		InfoCommand(),
		UptimeCommand(),
		InviteCommand(),
		StatsCommand(),
		PrefixCommand(),
	}
}

// PingCommand cria o comando de ping
func PingCommand() *Command {
	return &Command{
		Name:        "ping",
		Aliases:     []string{"p"},
		Description: "Verifica a latência do bot",
		Usage:       "ping",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			start := time.Now()
			msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "🏓 Pong!")
			if err != nil {
				return err
			}

			latency := time.Since(start).Milliseconds()
			_, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID,
				fmt.Sprintf("🏓 Pong!\nLatência: %dms\nLatência da API: %dms",
					latency, ctx.Session.HeartbeatLatency().Milliseconds()))
			return err
		},
	}
}

// HelpCommand cria o comando de ajuda
func HelpCommand() *Command {
	return &Command{
		Name:        "help",
		Aliases:     []string{"h", "ajuda"},
		Description: "Mostra a lista de comandos disponíveis",
		Usage:       "help [comando]",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			// Se um comando específico foi solicitado
			if len(ctx.Args) > 0 {
				cmdName := strings.ToLower(ctx.Args[0])
				cmd := ctx.Registry.GetCommand(cmdName)
				if cmd == nil {
					_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID,
						fmt.Sprintf("❌ Comando `%s` não encontrado.", cmdName))
					return err
				}

				embed := &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("📖 Ajuda: %s", cmd.Name),
					Color:       0x00ff00,
					Description: cmd.Description,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Uso",
							Value:  fmt.Sprintf("`%s%s`", ctx.Registry.GetPrefix(), cmd.Usage),
							Inline: false,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Sir Draith - Bot Medieval",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}

				if len(cmd.Aliases) > 0 {
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:   "Aliases",
						Value:  fmt.Sprintf("`%s`", strings.Join(cmd.Aliases, "`, `")),
						Inline: false,
					})
				}

				_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
				return err
			}

			// Agrupa comandos por categoria
			categories := make(map[string][]*Command)
			for _, cmd := range ctx.Registry.GetCommands() {
				if categories[cmd.Category] == nil {
					categories[cmd.Category] = []*Command{}
				}
				// Evita duplicatas de aliases
				found := false
				for _, c := range categories[cmd.Category] {
					if c.Name == cmd.Name {
						found = true
						break
					}
				}
				if !found {
					categories[cmd.Category] = append(categories[cmd.Category], cmd)
				}
			}

			embed := &discordgo.MessageEmbed{
				Title:       "📜 Comandos Disponíveis",
				Color:       0x00ff00,
				Description: fmt.Sprintf("Use `%shelp <comando>` para mais informações sobre um comando específico.", ctx.Registry.GetPrefix()),
				Fields:      []*discordgo.MessageEmbedField{},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Adiciona campos para cada categoria
			for category, cmds := range categories {
				var cmdList []string
				for _, cmd := range cmds {
					cmdList = append(cmdList, fmt.Sprintf("`%s`", cmd.Name))
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   category,
					Value:  strings.Join(cmdList, ", "),
					Inline: false,
				})
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// InfoCommand cria o comando de informações
func InfoCommand() *Command {
	return &Command{
		Name:        "info",
		Aliases:     []string{"i", "sobre"},
		Description: "Mostra informações sobre o bot",
		Usage:       "info",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			embed := &discordgo.MessageEmbed{
				Title: "ℹ️ Sobre Sir Draith",
				Color: 0x00ff00,
				Description: strings.Join([]string{
					"Sir Draith é um bot medieval para RPG de mesa no Discord.",
					"",
					"**Características:**",
					"• Sistema de personagens",
					"• Sistema de batalha com cartas",
					"• Sistema de narrativa",
					"• Gerenciamento de campanhas",
				}, "\n"),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Versão",
						Value:  "1.0.0",
						Inline: true,
					},
					{
						Name:   "Desenvolvedor",
						Value:  "SamSepi0l",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// UptimeCommand cria o comando para mostrar o tempo de atividade do bot
func UptimeCommand() *Command {
	return &Command{
		Name:        "uptime",
		Aliases:     []string{"tempo", "online"},
		Description: "Mostra há quanto tempo o bot está online",
		Usage:       "uptime",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			uptime := time.Since(startTime)
			days := int(uptime.Hours() / 24)
			hours := int(uptime.Hours()) % 24
			minutes := int(uptime.Minutes()) % 60
			seconds := int(uptime.Seconds()) % 60

			embed := &discordgo.MessageEmbed{
				Title:       "⏱️ Tempo Online",
				Color:       0x00ff00,
				Description: fmt.Sprintf("Estou online há: **%d dias, %d horas, %d minutos e %d segundos**", days, hours, minutes, seconds),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// InviteCommand cria o comando para obter o link de convite do bot
func InviteCommand() *Command {
	return &Command{
		Name:        "invite",
		Aliases:     []string{"convite", "convidar"},
		Description: "Mostra o link para adicionar o bot ao seu servidor",
		Usage:       "invite",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			app, err := ctx.Session.Application("@me")
			if err != nil {
				return fmt.Errorf("erro ao obter informações do bot: %w", err)
			}

			// Permissões necessárias para o bot funcionar
			permissions := discordgo.PermissionManageRoles |
				discordgo.PermissionManageMessages |
				discordgo.PermissionReadMessages |
				discordgo.PermissionSendMessages |
				discordgo.PermissionEmbedLinks |
				discordgo.PermissionAttachFiles |
				discordgo.PermissionReadMessageHistory |
				discordgo.PermissionMentionEveryone |
				discordgo.PermissionUseExternalEmojis |
				discordgo.PermissionKickMembers |
				discordgo.PermissionBanMembers

			inviteLink := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
				app.ID, permissions)

			embed := &discordgo.MessageEmbed{
				Title:       "🔗 Link de Convite",
				Color:       0x00ff00,
				Description: fmt.Sprintf("Clique [aqui](%s) para me adicionar ao seu servidor!", inviteLink),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// StatsCommand cria o comando para mostrar estatísticas do bot
func StatsCommand() *Command {
	return &Command{
		Name:        "stats",
		Aliases:     []string{"estatisticas", "status"},
		Description: "Mostra estatísticas do bot",
		Usage:       "stats",
		Category:    "Geral",
		Handler: func(ctx *CommandContext) error {
			// Conta servidores, usuários e canais
			var (
				guildCount   = len(ctx.Session.State.Guilds)
				userCount    = 0
				channelCount = 0
			)

			for _, guild := range ctx.Session.State.Guilds {
				userCount += guild.MemberCount
				channelCount += len(guild.Channels)
			}

			embed := &discordgo.MessageEmbed{
				Title: "📊 Estatísticas do Bot",
				Color: 0x00ff00,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🏰 Servidores",
						Value:  fmt.Sprintf("%d", guildCount),
						Inline: true,
					},
					{
						Name:   "👥 Usuários",
						Value:  fmt.Sprintf("%d", userCount),
						Inline: true,
					},
					{
						Name:   "💬 Canais",
						Value:  fmt.Sprintf("%d", channelCount),
						Inline: true,
					},
					{
						Name:   "⚡ Latência",
						Value:  fmt.Sprintf("%dms", ctx.Session.HeartbeatLatency().Milliseconds()),
						Inline: true,
					},
					{
						Name:   "⏱️ Uptime",
						Value:  fmt.Sprintf("%s", time.Since(startTime).Round(time.Second)),
						Inline: true,
					},
					{
						Name:   "🤖 Versão",
						Value:  "1.0.0",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// PrefixCommand cria o comando para mostrar ou alterar o prefixo dos comandos
func PrefixCommand() *Command {
	return &Command{
		Name:        "prefix",
		Aliases:     []string{"prefixo"},
		Description: "Mostra ou altera o prefixo dos comandos",
		Usage:       "prefix [novo prefixo]",
		Category:    "Admin",
		Handler: func(ctx *CommandContext) error {
			// Se não houver argumentos, mostra o prefixo atual
			if len(ctx.Args) == 0 {
				prefix, err := ctx.Registry.GetGuildPrefix(ctx.Message.GuildID)
				if err != nil {
					return fmt.Errorf("erro ao obter prefixo: %w", err)
				}

				embed := &discordgo.MessageEmbed{
					Title:       "⚙️ Prefixo Atual",
					Color:       0x00ff00,
					Description: fmt.Sprintf("O prefixo atual é: `%s`", prefix),
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Use %sprefix <novo prefixo> para alterar", prefix),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}

				_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
				return err
			}

			// Verifica se o usuário tem permissão para alterar o prefixo
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionManageServer == 0 {
				return fmt.Errorf("você não tem permissão para alterar o prefixo")
			}

			// Valida o novo prefixo
			newPrefix := ctx.Args[0]
			if len(newPrefix) > 3 {
				return fmt.Errorf("o prefixo não pode ter mais de 3 caracteres")
			}

			// Atualiza o prefixo no banco de dados
			err = ctx.Registry.SetGuildPrefix(ctx.Message.GuildID, newPrefix)
			if err != nil {
				return fmt.Errorf("erro ao atualizar prefixo: %w", err)
			}

			embed := &discordgo.MessageEmbed{
				Title:       "⚙️ Prefixo Alterado",
				Color:       0x00ff00,
				Description: fmt.Sprintf("O prefixo foi alterado para: `%s`", newPrefix),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Sir Draith - Bot Medieval",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
} 