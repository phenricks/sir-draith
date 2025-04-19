package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// AdminCommands retorna os comandos administrativos do bot
func AdminCommands() []*Command {
	return []*Command{
		ClearCommand(),
		KickCommand(),
		BanCommand(),
		UnbanCommand(),
	}
}

// ClearCommand cria o comando para limpar mensagens
func ClearCommand() *Command {
	return &Command{
		Name:        "clear",
		Aliases:     []string{"limpar", "clean"},
		Description: "Limpa mensagens do canal atual",
		Usage:       "clear <quantidade>",
		Category:    "Admin",
		Handler: func(ctx *CommandContext) error {
			// Verifica se o usuário tem permissão
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionManageMessages == 0 {
				return fmt.Errorf("você não tem permissão para usar este comando")
			}

			// Verifica se a quantidade foi especificada
			if len(ctx.Args) != 1 {
				return fmt.Errorf("use: %s%s", ctx.Registry.GetPrefix(), "clear <quantidade>")
			}

			// Converte a quantidade para número
			amount := 0
			_, err = fmt.Sscanf(ctx.Args[0], "%d", &amount)
			if err != nil || amount < 1 || amount > 100 {
				return fmt.Errorf("quantidade inválida. Use um número entre 1 e 100")
			}

			// Obtém as mensagens
			messages, err := ctx.Session.ChannelMessages(ctx.Message.ChannelID, amount+1, "", "", ctx.Message.ID)
			if err != nil {
				return fmt.Errorf("erro ao obter mensagens: %w", err)
			}

			// Extrai os IDs das mensagens
			messageIDs := make([]string, len(messages))
			for i, msg := range messages {
				messageIDs[i] = msg.ID
			}

			// Deleta as mensagens
			err = ctx.Session.ChannelMessagesBulkDelete(ctx.Message.ChannelID, messageIDs)
			if err != nil {
				return fmt.Errorf("erro ao deletar mensagens: %w", err)
			}

			// Envia confirmação e deleta após 5 segundos
			msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID,
				fmt.Sprintf("✅ %d mensagens foram deletadas!", len(messageIDs)-1))
			if err != nil {
				return err
			}

			time.Sleep(5 * time.Second)
			return ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, msg.ID)
		},
	}
}

// KickCommand cria o comando para expulsar usuários
func KickCommand() *Command {
	return &Command{
		Name:        "kick",
		Aliases:     []string{"expulsar"},
		Description: "Expulsa um usuário do servidor",
		Usage:       "kick <@usuário> [motivo]",
		Category:    "Admin",
		Handler: func(ctx *CommandContext) error {
			// Verifica se o usuário tem permissão
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionKickMembers == 0 {
				return fmt.Errorf("você não tem permissão para usar este comando")
			}

			// Verifica se o usuário foi mencionado
			if len(ctx.Message.Mentions) != 1 {
				return fmt.Errorf("mencione o usuário que deseja expulsar")
			}

			// Obtém o motivo
			reason := "Nenhum motivo especificado"
			if len(ctx.Args) > 1 {
				reason = strings.Join(ctx.Args[1:], " ")
			}

			// Expulsa o usuário
			err = ctx.Session.GuildMemberDelete(ctx.Message.GuildID, ctx.Message.Mentions[0].ID)
			if err != nil {
				return fmt.Errorf("erro ao expulsar usuário: %w", err)
			}

			// Envia confirmação
			embed := &discordgo.MessageEmbed{
				Title:       "👢 Usuário Expulso",
				Color:       0xff0000,
				Description: fmt.Sprintf("**Usuário:** %s\n**Motivo:** %s", ctx.Message.Mentions[0].Username, reason),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Expulso por %s", ctx.Message.Author.Username),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// BanCommand cria o comando para banir usuários
func BanCommand() *Command {
	return &Command{
		Name:        "ban",
		Aliases:     []string{"banir"},
		Description: "Bane um usuário do servidor",
		Usage:       "ban <@usuário> [motivo]",
		Category:    "Admin",
		Handler: func(ctx *CommandContext) error {
			// Verifica se o usuário tem permissão
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionBanMembers == 0 {
				return fmt.Errorf("você não tem permissão para usar este comando")
			}

			// Verifica se o usuário foi mencionado
			if len(ctx.Message.Mentions) != 1 {
				return fmt.Errorf("mencione o usuário que deseja banir")
			}

			// Obtém o motivo
			reason := "Nenhum motivo especificado"
			if len(ctx.Args) > 1 {
				reason = strings.Join(ctx.Args[1:], " ")
			}

			// Bane o usuário
			err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, ctx.Message.Mentions[0].ID, reason, 1)
			if err != nil {
				return fmt.Errorf("erro ao banir usuário: %w", err)
			}

			// Envia confirmação
			embed := &discordgo.MessageEmbed{
				Title:       "🔨 Usuário Banido",
				Color:       0xff0000,
				Description: fmt.Sprintf("**Usuário:** %s\n**Motivo:** %s", ctx.Message.Mentions[0].Username, reason),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Banido por %s", ctx.Message.Author.Username),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
}

// UnbanCommand cria o comando para desbanir usuários
func UnbanCommand() *Command {
	return &Command{
		Name:        "unban",
		Aliases:     []string{"desbanir"},
		Description: "Remove o banimento de um usuário",
		Usage:       "unban <ID do usuário>",
		Category:    "Admin",
		Handler: func(ctx *CommandContext) error {
			// Verifica se o usuário tem permissão
			perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				return fmt.Errorf("erro ao verificar permissões: %w", err)
			}

			if perms&discordgo.PermissionBanMembers == 0 {
				return fmt.Errorf("você não tem permissão para usar este comando")
			}

			// Verifica se o ID foi fornecido
			if len(ctx.Args) != 1 {
				return fmt.Errorf("forneça o ID do usuário que deseja desbanir")
			}

			// Remove o banimento
			err = ctx.Session.GuildBanDelete(ctx.Message.GuildID, ctx.Args[0])
			if err != nil {
				return fmt.Errorf("erro ao desbanir usuário: %w", err)
			}

			// Envia confirmação
			embed := &discordgo.MessageEmbed{
				Title:       "🔓 Usuário Desbanido",
				Color:       0x00ff00,
				Description: fmt.Sprintf("**ID do Usuário:** %s", ctx.Args[0]),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Desbanido por %s", ctx.Message.Author.Username),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
			return err
		},
	}
} 