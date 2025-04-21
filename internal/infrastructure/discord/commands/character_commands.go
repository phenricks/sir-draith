package commands

import (
	"context"
	"fmt"
	"strings"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/services"

	"github.com/bwmarrin/discordgo"
)

// CharacterCommands contém os comandos relacionados a personagens
type CharacterCommands struct {
	characterService *services.CharacterService
}

// NewCharacterCommands cria uma nova instância dos comandos de personagem
func NewCharacterCommands(characterService *services.CharacterService) *CharacterCommands {
	return &CharacterCommands{
		characterService: characterService,
	}
}

// Register registra todos os comandos de personagem
func (cc *CharacterCommands) Register(registry *CommandRegistry) {
	// Comando para criar um novo personagem
	registry.RegisterCommand(&Command{
		Name:        "criar",
		Aliases:     []string{"create", "new"},
		Description: "Cria um novo personagem",
		Usage:       "criar <nome>",
		Category:    "Personagem",
		Handler:     cc.handleCreate,
	})

	// Comando para ver informações do personagem
	registry.RegisterCommand(&Command{
		Name:        "personagem",
		Aliases:     []string{"char", "info"},
		Description: "Mostra informações do seu personagem",
		Usage:       "personagem",
		Category:    "Personagem",
		Handler:     cc.handleInfo,
	})

	// Comando para listar personagens
	registry.RegisterCommand(&Command{
		Name:        "personagens",
		Aliases:     []string{"chars", "list"},
		Description: "Lista todos os personagens do servidor",
		Usage:       "personagens",
		Category:    "Personagem",
		Handler:     cc.handleList,
	})

	// Comando para gerenciar inventário
	registry.RegisterCommand(&Command{
		Name:        "inventario",
		Aliases:     []string{"inv", "i"},
		Description: "Gerencia o inventário do seu personagem",
		Usage:       "inventario [equipar/desequipar/usar] [item]",
		Category:    "Personagem",
		Handler:     cc.handleInventory,
	})
}

// handleCreate lida com o comando de criar personagem
func (cc *CharacterCommands) handleCreate(ctx *CommandContext) error {
	// Verifica se já tem um personagem
	character, err := cc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err == nil && character != nil {
		return ctx.Reply("Você já possui um personagem neste servidor!")
	}

	// Verifica se forneceu um nome
	if len(ctx.Args) == 0 {
		return ctx.Reply("Por favor, forneça um nome para seu personagem!")
	}

	// Inicia o wizard de criação de personagem
	name := strings.Join(ctx.Args, " ")
	wizard := NewCharacterWizard(ctx.Session, cc.characterService)
	wizard.character = &entities.Character{
		UserID:  ctx.Message.Author.ID,
		GuildID: ctx.Message.GuildID,
		Name:    name,
	}
	wizard.channelID = ctx.Message.ChannelID

	// Registra o wizard no registro de comandos
	ctx.Registry.wizards[ctx.Message.Author.ID] = wizard

	// Inicia o processo de criação
	embed := &discordgo.MessageEmbed{
		Title:       "⚔️ Criação de Personagem",
		Description: fmt.Sprintf("Vamos criar seu personagem **%s**!\nEscolha sua classe:", name),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "⚔️ Guerreiro",
				Value: "Mestre das armas e da guerra, especialista em combate corpo a corpo",
			},
			{
				Name:  "🏹 Arqueiro",
				Value: "Ágil e preciso, domina o combate à distância",
			},
			{
				Name:  "🔮 Mago",
				Value: "Estudioso das artes arcanas, manipula a magia para diversos fins",
			},
			{
				Name:  "🛡️ Paladino",
				Value: "Guerreiro sagrado que combina combate com poderes divinos",
			},
			{
				Name:  "🌿 Druida",
				Value: "Guardião da natureza com poderes elementais e metamorfose",
			},
			{
				Name:  "✝️ Clérigo",
				Value: "Servo divino com poderes de cura e proteção",
			},
			{
				Name:  "🎭 Bardo",
				Value: "Artista versátil que combina música com magia",
			},
			{
				Name:  "👻 Bruxo",
				Value: "Místico que fez pacto com entidades poderosas",
			},
			{
				Name:  "🌟 Feiticeiro",
				Value: "Usuário nato de magia com poderes inatos",
			},
			{
				Name:  "🗡️ Ladino",
				Value: "Especialista em furtividade e ataques precisos",
			},
			{
				Name:  "🥋 Monge",
				Value: "Mestre das artes marciais e da energia ki",
			},
			{
				Name:  "💢 Bárbaro",
				Value: "Guerreiro selvagem movido pela fúria",
			},
		},
	}

	// Cria os botões de seleção de classe
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Guerreiro",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_warrior",
					Emoji: discordgo.ComponentEmoji{
						Name: "⚔️",
					},
				},
				discordgo.Button{
					Label:    "Arqueiro",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_ranger",
					Emoji: discordgo.ComponentEmoji{
						Name: "🏹",
					},
				},
				discordgo.Button{
					Label:    "Mago",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_mage",
					Emoji: discordgo.ComponentEmoji{
						Name: "🔮",
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Paladino",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_paladin",
					Emoji: discordgo.ComponentEmoji{
						Name: "🛡️",
					},
				},
				discordgo.Button{
					Label:    "Druida",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_druid",
					Emoji: discordgo.ComponentEmoji{
						Name: "🌿",
					},
				},
				discordgo.Button{
					Label:    "Clérigo",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_cleric",
					Emoji: discordgo.ComponentEmoji{
						Name: "✝️",
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Bardo",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_bard",
					Emoji: discordgo.ComponentEmoji{
						Name: "🎭",
					},
				},
				discordgo.Button{
					Label:    "Bruxo",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_warlock",
					Emoji: discordgo.ComponentEmoji{
						Name: "👻",
					},
				},
				discordgo.Button{
					Label:    "Feiticeiro",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_sorcerer",
					Emoji: discordgo.ComponentEmoji{
						Name: "🌟",
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Ladino",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_rogue",
					Emoji: discordgo.ComponentEmoji{
						Name: "🗡️",
					},
				},
				discordgo.Button{
					Label:    "Monge",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_monk",
					Emoji: discordgo.ComponentEmoji{
						Name: "🥋",
					},
				},
				discordgo.Button{
					Label:    "Bárbaro",
					Style:    discordgo.PrimaryButton,
					CustomID: "class_barbarian",
					Emoji: discordgo.ComponentEmoji{
						Name: "💢",
					},
				},
			},
		},
	}

	// Envia a mensagem inicial do wizard
	msg, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		return fmt.Errorf("erro ao enviar mensagem inicial: %w", err)
	}

	// Armazena o ID da mensagem no wizard
	wizard.messageID = msg.ID

	return nil
}

// handleInfo lida com o comando de ver informações do personagem
func (cc *CharacterCommands) handleInfo(ctx *CommandContext) error {
	// Busca o personagem
	character, err := cc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply("Você não possui um personagem neste servidor!")
	}

	// Cria embed com informações detalhadas
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📜 %s", character.Name),
		Description: character.Description,
		Color:       0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Nível",
				Value:  fmt.Sprintf("%d", character.Level),
				Inline: true,
			},
			{
				Name:   "Experiência",
				Value:  fmt.Sprintf("%d", character.Experience),
				Inline: true,
			},
			{
				Name:   "Ouro",
				Value:  fmt.Sprintf("%d", character.Gold),
				Inline: true,
			},
			{
				Name:   "Classe",
				Value:  string(character.Class),
				Inline: true,
			},
			{
				Name:   "Origem",
				Value:  character.Background,
				Inline: true,
			},
			{
				Name: "Atributos",
				Value: fmt.Sprintf(
					"Força: %d\nDestreza: %d\nConstituição: %d\nInteligência: %d\nSabedoria: %d\nCarisma: %d",
					character.Attributes.Strength,
					character.Attributes.Dexterity,
					character.Attributes.Constitution,
					character.Attributes.Intelligence,
					character.Attributes.Wisdom,
					character.Attributes.Charisma,
				),
				Inline: false,
			},
			{
				Name: "Combate",
				Value: fmt.Sprintf(
					"Vida: %d/%d\nArmadura: %d\nIniciativa: %d",
					character.Combat.Health,
					character.Combat.MaxHealth,
					character.Combat.Armor,
					character.Combat.Initiative,
				),
				Inline: false,
			},
		},
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleList lida com o comando de listar personagens
func (cc *CharacterCommands) handleList(ctx *CommandContext) error {
	// Busca todos os personagens do servidor
	characters, err := cc.characterService.ListCharactersByGuild(context.Background(), ctx.Message.GuildID)
	if err != nil {
		return fmt.Errorf("erro ao listar personagens: %w", err)
	}

	if len(characters) == 0 {
		return ctx.Reply("Não há personagens neste servidor ainda!")
	}

	// Cria embed com a lista de personagens
	embed := &discordgo.MessageEmbed{
		Title: "📚 Personagens do Servidor",
		Color: 0x0099ff,
	}

	for _, char := range characters {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: char.Name,
			Value: fmt.Sprintf(
				"Nível %d %s\nVida: %d/%d",
				char.Level,
				char.Class,
				char.Combat.Health,
				char.Combat.MaxHealth,
			),
			Inline: true,
		})
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleInventory lida com o comando de gerenciar inventário
func (cc *CharacterCommands) handleInventory(ctx *CommandContext) error {
	// Busca o personagem
	character, err := cc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply("Você não possui um personagem neste servidor!")
	}

	// Se não houver argumentos, mostra o inventário
	if len(ctx.Args) == 0 {
		embed := &discordgo.MessageEmbed{
			Title: "🎒 Inventário",
			Color: 0x0099ff,
		}

		// Lista itens equipados
		equipmentList := "Nenhum item equipado"
		if len(character.Equipment) > 0 {
			items := make([]string, 0)
			for _, item := range character.Equipment {
				items = append(items, fmt.Sprintf("**%s** (%s)", item.Name, item.Slot))
			}
			equipmentList = strings.Join(items, "\n")
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📦 Equipado",
			Value:  equipmentList,
			Inline: false,
		})

		// Lista itens no inventário
		inventoryList := "Inventário vazio"
		if len(character.Inventory) > 0 {
			items := make([]string, 0)
			for _, item := range character.Inventory {
				items = append(items, fmt.Sprintf("**%s** x%d", item.Name, item.Quantity))
			}
			inventoryList = strings.Join(items, "\n")
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🎒 Itens",
			Value:  inventoryList,
			Inline: false,
		})

		_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
		return err
	}

	// Processa subcomandos do inventário
	subcommand := strings.ToLower(ctx.Args[0])
	switch subcommand {
	case "equipar", "equip":
		if len(ctx.Args) < 2 {
			return ctx.Reply("Por favor, especifique qual item deseja equipar!")
		}
		itemName := strings.Join(ctx.Args[1:], " ")
		if err := cc.characterService.EquipItem(context.Background(), character, itemName); err != nil {
			return ctx.Reply(fmt.Sprintf("Erro ao equipar item: %s", err))
		}
		return ctx.Reply(fmt.Sprintf("✅ **%s** equipado com sucesso!", itemName))

	case "desequipar", "unequip":
		if len(ctx.Args) < 2 {
			return ctx.Reply("Por favor, especifique qual item deseja desequipar!")
		}
		itemName := strings.Join(ctx.Args[1:], " ")
		if err := cc.characterService.UnequipItem(context.Background(), character, itemName); err != nil {
			return ctx.Reply(fmt.Sprintf("Erro ao desequipar item: %s", err))
		}
		return ctx.Reply(fmt.Sprintf("✅ **%s** desequipado com sucesso!", itemName))

	default:
		return ctx.Reply("Subcomando inválido! Use: equipar, desequipar")
	}
}
