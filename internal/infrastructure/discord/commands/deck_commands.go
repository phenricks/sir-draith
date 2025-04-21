package commands

import (
	"context"
	"fmt"
	"strings"

	"sirdraith/internal/domain/services"

	"github.com/bwmarrin/discordgo"
)

// DeckCommands encapsula os comandos relacionados a decks
type DeckCommands struct {
	deckService *services.DeckService
}

// NewDeckCommands cria uma nova instância de DeckCommands
func NewDeckCommands(deckService *services.DeckService) *DeckCommands {
	return &DeckCommands{
		deckService: deckService,
	}
}

// Register registra os comandos de deck
func (dc *DeckCommands) Register(registry *CommandRegistry) {
	registry.RegisterCommand(&Command{
		Name:        "criar-deck",
		Description: "Cria um novo deck",
		Usage:       "criar-deck <nome> <classe> [descrição]",
		Category:    "Decks",
		Handler:     dc.handleCreate,
	})

	registry.RegisterCommand(&Command{
		Name:        "editar-deck",
		Description: "Edita um deck existente",
		Usage:       "editar-deck <id> <adicionar/remover> <carta>",
		Category:    "Decks",
		Handler:     dc.handleEdit,
	})

	registry.RegisterCommand(&Command{
		Name:        "deletar-deck",
		Description: "Remove um deck",
		Usage:       "deletar-deck <id>",
		Category:    "Decks",
		Handler:     dc.handleDelete,
	})

	registry.RegisterCommand(&Command{
		Name:        "listar-decks",
		Description: "Lista seus decks",
		Usage:       "listar-decks",
		Category:    "Decks",
		Handler:     dc.handleList,
	})

	registry.RegisterCommand(&Command{
		Name:        "ver-deck",
		Description: "Mostra detalhes de um deck",
		Usage:       "ver-deck <id>",
		Category:    "Decks",
		Handler:     dc.handleView,
	})
}

// handleCreate processa o comando de criar deck
func (dc *DeckCommands) handleCreate(ctx *CommandContext) error {
	if len(ctx.Args) < 2 {
		return ctx.Reply("Por favor, forneça o nome e a classe do deck!")
	}

	name := ctx.Args[0]
	class := ctx.Args[1]
	description := ""
	if len(ctx.Args) > 2 {
		description = strings.Join(ctx.Args[2:], " ")
	}

	deck, err := dc.deckService.CreateDeck(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID, name, description, class)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("Erro ao criar deck: %s", err))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎴 Deck Criado",
		Description: fmt.Sprintf("Deck **%s** criado com sucesso!", deck.Name),
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  deck.ID.Hex(),
				Inline: true,
			},
			{
				Name:   "Classe",
				Value:  deck.Class,
				Inline: true,
			},
			{
				Name:   "Descrição",
				Value:  deck.Description,
				Inline: false,
			},
		},
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleEdit processa o comando de editar deck
func (dc *DeckCommands) handleEdit(ctx *CommandContext) error {
	if len(ctx.Args) < 3 {
		return ctx.Reply("Por favor, forneça o ID do deck, a ação (adicionar/remover) e o ID da carta!")
	}

	deckID := ctx.Args[0]
	action := strings.ToLower(ctx.Args[1])
	cardID := ctx.Args[2]

	var err error
	switch action {
	case "adicionar":
		err = dc.deckService.AddCardToDeck(context.Background(), deckID, cardID)
	case "remover":
		err = dc.deckService.RemoveCardFromDeck(context.Background(), deckID, cardID)
	default:
		return ctx.Reply("Ação inválida! Use 'adicionar' ou 'remover'.")
	}

	if err != nil {
		return ctx.Reply(fmt.Sprintf("Erro ao editar deck: %s", err))
	}

	return ctx.Reply("✅ Deck atualizado com sucesso!")
}

// handleDelete processa o comando de deletar deck
func (dc *DeckCommands) handleDelete(ctx *CommandContext) error {
	if len(ctx.Args) < 1 {
		return ctx.Reply("Por favor, forneça o ID do deck!")
	}

	deckID := ctx.Args[0]
	if err := dc.deckService.DeleteDeck(context.Background(), deckID); err != nil {
		return ctx.Reply(fmt.Sprintf("Erro ao deletar deck: %s", err))
	}

	return ctx.Reply("✅ Deck deletado com sucesso!")
}

// handleList processa o comando de listar decks
func (dc *DeckCommands) handleList(ctx *CommandContext) error {
	decks, err := dc.deckService.ListDecksByUser(context.Background(), ctx.Message.Author.ID)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("Erro ao listar decks: %s", err))
	}

	if len(decks) == 0 {
		return ctx.Reply("Você não possui nenhum deck!")
	}

	embed := &discordgo.MessageEmbed{
		Title: "📚 Seus Decks",
		Color: 0x0099ff,
	}

	for _, deck := range decks {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: deck.Name,
			Value: fmt.Sprintf(
				"ID: %s\nClasse: %s\nCartas: %d",
				deck.ID.Hex(),
				deck.Class,
				len(deck.Cards),
			),
			Inline: true,
		})
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleView processa o comando de ver deck
func (dc *DeckCommands) handleView(ctx *CommandContext) error {
	if len(ctx.Args) < 1 {
		return ctx.Reply("Por favor, forneça o ID do deck!")
	}

	deckID := ctx.Args[0]
	deck, err := dc.deckService.GetDeck(context.Background(), deckID)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("Erro ao buscar deck: %s", err))
	}
	if deck == nil {
		return ctx.Reply("Deck não encontrado!")
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎴 %s", deck.Name),
		Description: deck.Description,
		Color:       0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  deck.ID.Hex(),
				Inline: true,
			},
			{
				Name:   "Classe",
				Value:  deck.Class,
				Inline: true,
			},
			{
				Name:   "Total de Cartas",
				Value:  fmt.Sprintf("%d", len(deck.Cards)),
				Inline: true,
			},
		},
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}
