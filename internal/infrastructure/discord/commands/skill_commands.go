package commands

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"sirdraith/internal/domain/gamedata"
	"sirdraith/internal/domain/services"

	"github.com/bwmarrin/discordgo"
)

// SkillCommands encapsula os comandos relacionados a perícias
type SkillCommands struct {
	characterService *services.CharacterService
}

// NewSkillCommands cria uma nova instância de SkillCommands
func NewSkillCommands(characterService *services.CharacterService) *SkillCommands {
	return &SkillCommands{
		characterService: characterService,
	}
}

// Register registra os comandos de perícia
func (sc *SkillCommands) Register(registry *CommandRegistry) {
	registry.RegisterCommand(&Command{
		Name:        "roll",
		Description: "Realiza um teste de perícia",
		Usage:       "roll <perícia> [dificuldade]",
		Category:    "Perícias",
		Handler:     sc.handleRoll,
	})

	registry.RegisterCommand(&Command{
		Name:        "skills",
		Description: "Lista as perícias do seu personagem",
		Usage:       "skills",
		Category:    "Perícias",
		Handler:     sc.handleList,
	})

	registry.RegisterCommand(&Command{
		Name:        "skillbonus",
		Description: "Gerencia bônus temporários em perícias",
		Usage:       "skillbonus <add/remove> <perícia> <valor>",
		Category:    "Perícias",
		Handler:     sc.handleBonus,
	})
}

// handleRoll processa o comando de teste de perícia
func (sc *SkillCommands) handleRoll(ctx *CommandContext) error {
	// Buscar o personagem
	character, err := sc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply("Você não possui um personagem neste servidor!")
	}

	// Verificar argumentos
	if len(ctx.Args) < 1 {
		return ctx.Reply("Por favor, especifique qual perícia deseja testar!")
	}

	// Identificar a perícia
	skillName := strings.ToLower(ctx.Args[0])
	var skill gamedata.Skill
	found := false
	for s := range gamedata.SkillBaseAttribute {
		if strings.ToLower(string(s)) == skillName {
			skill = s
			found = true
			break
		}
	}

	if !found {
		return ctx.Reply("Perícia inválida! Use uma das perícias disponíveis.")
	}

	// Definir dificuldade (DC)
	dc := 10 // Dificuldade padrão
	if len(ctx.Args) > 1 {
		if n, err := parseInt(ctx.Args[1]); err == nil {
			dc = n
		}
	}

	// Rolar o dado
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(20) + 1
	modifier := character.GetSkillModifier(skill)
	total := roll + modifier

	// Determinar o resultado
	success := total >= dc
	resultEmoji := "❌"
	if success {
		resultEmoji = "✅"
	}

	// Criar embed com o resultado
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎲 Teste de %s", skill),
		Description: fmt.Sprintf("%s Resultado: %d", resultEmoji, total),
		Color:       getResultColor(success),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Rolagem",
				Value:  fmt.Sprintf("%d", roll),
				Inline: true,
			},
			{
				Name:   "Modificador",
				Value:  fmt.Sprintf("%+d", modifier),
				Inline: true,
			},
			{
				Name:   "Dificuldade",
				Value:  fmt.Sprintf("%d", dc),
				Inline: true,
			},
		},
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleList processa o comando de listar perícias
func (sc *SkillCommands) handleList(ctx *CommandContext) error {
	// Buscar o personagem
	character, err := sc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply("Você não possui um personagem neste servidor!")
	}

	// Agrupar perícias por atributo base
	skillsByAttr := make(map[string][]gamedata.SkillProficiency)
	for _, prof := range character.Skills {
		baseAttr := gamedata.SkillBaseAttribute[prof.Skill]
		skillsByAttr[baseAttr] = append(skillsByAttr[baseAttr], prof)
	}

	// Criar embed
	embed := &discordgo.MessageEmbed{
		Title:       "🎯 Perícias",
		Description: fmt.Sprintf("Perícias de %s", character.Name),
		Color:       0x0099ff,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Emojis para cada atributo
	attrEmojis := map[string]string{
		"strength":     "💪",
		"dexterity":    "🏃",
		"constitution": "❤️",
		"intelligence": "🧠",
		"wisdom":       "🦉",
		"charisma":     "👑",
	}

	// Adicionar campos para cada grupo de perícias
	for attr, skills := range skillsByAttr {
		skillList := ""
		for _, prof := range skills {
			// Calcular o modificador total
			modifier := character.GetSkillModifier(prof.Skill)
			profMark := "  "
			if prof.IsProficient {
				profMark = "✓ "
			}
			bonusStr := ""
			if prof.Bonus != 0 {
				bonusStr = fmt.Sprintf(" (%+d)", prof.Bonus)
			}
			skillList += fmt.Sprintf("%s%s: %+d%s\n", profMark, prof.Skill, modifier, bonusStr)
		}

		if skillList != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s %s", attrEmojis[attr], strings.Title(attr)),
				Value:  skillList,
				Inline: true,
			})
		}
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}

// handleBonus processa o comando de gerenciar bônus de perícia
func (sc *SkillCommands) handleBonus(ctx *CommandContext) error {
	// Buscar o personagem
	character, err := sc.characterService.GetCharacterByUserAndGuild(context.Background(), ctx.Message.Author.ID, ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply("Você não possui um personagem neste servidor!")
	}

	// Verificar argumentos
	if len(ctx.Args) < 3 {
		return ctx.Reply("Uso: skillbonus <add/remove> <perícia> <valor>")
	}

	action := strings.ToLower(ctx.Args[0])
	skillName := strings.ToLower(ctx.Args[1])
	bonus, err := parseInt(ctx.Args[2])
	if err != nil {
		return ctx.Reply("O valor do bônus deve ser um número!")
	}

	// Identificar a perícia
	var skill gamedata.Skill
	found := false
	for s := range gamedata.SkillBaseAttribute {
		if strings.ToLower(string(s)) == skillName {
			skill = s
			found = true
			break
		}
	}

	if !found {
		return ctx.Reply("Perícia inválida! Use uma das perícias disponíveis.")
	}

	// Aplicar ou remover o bônus
	switch action {
	case "add":
		if err := character.AddSkillBonus(skill, bonus); err != nil {
			return ctx.Reply(fmt.Sprintf("Erro ao adicionar bônus: %s", err))
		}
	case "remove":
		if err := character.AddSkillBonus(skill, -bonus); err != nil {
			return ctx.Reply(fmt.Sprintf("Erro ao remover bônus: %s", err))
		}
	default:
		return ctx.Reply("Ação inválida! Use 'add' ou 'remove'.")
	}

	// Atualizar o personagem no banco
	if err := sc.characterService.UpdateCharacter(context.Background(), character); err != nil {
		return ctx.Reply("Erro ao atualizar personagem!")
	}

	return ctx.Reply(fmt.Sprintf("✅ Bônus de perícia atualizado com sucesso! Use `/skills` para ver suas perícias."))
}

// Funções auxiliares

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func getResultColor(success bool) int {
	if success {
		return 0x00ff00 // Verde
	}
	return 0xff0000 // Vermelho
}
