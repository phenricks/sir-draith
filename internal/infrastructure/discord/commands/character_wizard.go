package commands

import (
	"context"
	"fmt"
	"strings"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"
	"sirdraith/internal/domain/services"

	"github.com/bwmarrin/discordgo"
)

// WizardStep representa os passos do wizard de criação de personagem
type WizardStep int

const (
	StepInitial WizardStep = iota
	StepClass
	StepBackground
	StepSkills
	StepConfirmation
)

// CharacterWizard represents the character creation wizard
type CharacterWizard struct {
	session          *discordgo.Session
	character        *entities.Character
	characterService *services.CharacterService
	messageID        string // ID da mensagem do wizard
	channelID        string // ID do canal onde o wizard está ativo
	currentStep      WizardStep
}

// NewCharacterWizard creates a new character creation wizard
func NewCharacterWizard(session *discordgo.Session, characterService *services.CharacterService) *CharacterWizard {
	return &CharacterWizard{
		session:          session,
		characterService: characterService,
	}
}

// handleAttributeSelection processa a seleção de atributos
func (w *CharacterWizard) handleAttributeSelection(i *discordgo.InteractionCreate) error {
	data := i.MessageComponentData().CustomID

	// Extrair atributo e ação
	var attribute string
	var isIncrease bool

	switch data {
	case "attr_str_up":
		attribute = "strength"
		isIncrease = true
	case "attr_str_down":
		attribute = "strength"
		isIncrease = false
	case "attr_dex_up":
		attribute = "dexterity"
		isIncrease = true
	case "attr_dex_down":
		attribute = "dexterity"
		isIncrease = false
	case "attr_con_up":
		attribute = "constitution"
		isIncrease = true
	case "attr_con_down":
		attribute = "constitution"
		isIncrease = false
	case "attr_int_up":
		attribute = "intelligence"
		isIncrease = true
	case "attr_int_down":
		attribute = "intelligence"
		isIncrease = false
	case "attr_wis_up":
		attribute = "wisdom"
		isIncrease = true
	case "attr_wis_down":
		attribute = "wisdom"
		isIncrease = false
	case "attr_cha_up":
		attribute = "charisma"
		isIncrease = true
	case "attr_cha_down":
		attribute = "charisma"
		isIncrease = false
	case "attr_confirm":
		return w.handleAttributeConfirmation(i)
	default:
		return fmt.Errorf("ação de atributo inválida")
	}

	// Atualizar valor do atributo
	currentValue := w.character.Attributes.GetValue(attribute)
	remainingPoints := w.getRemainingPoints()

	if isIncrease {
		// Verificar se pode aumentar
		if currentValue >= 15 {
			return w.respondError(i, "Valor máximo atingido para este atributo")
		}
		cost := w.getAttributeCost(currentValue + 1)
		if remainingPoints < cost {
			return w.respondError(i, "Pontos insuficientes")
		}
		w.character.Attributes.SetValue(attribute, currentValue+1)
	} else {
		// Verificar se pode diminuir
		if currentValue <= 8 {
			return w.respondError(i, "Valor mínimo atingido para este atributo")
		}
		w.character.Attributes.SetValue(attribute, currentValue-1)
	}

	// Atualizar interface
	return w.updateAttributeInterface(i)
}

// getAttributeCost retorna o custo para aumentar um atributo para determinado valor
func (w *CharacterWizard) getAttributeCost(value int) int {
	if value <= 13 {
		return 1
	}
	if value == 14 {
		return 2
	}
	return 3 // value == 15
}

// getRemainingPoints calcula pontos restantes para distribuir
func (w *CharacterWizard) getRemainingPoints() int {
	totalPoints := 27
	usedPoints := 0

	// Calcular pontos gastos
	attributes := []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"}
	for _, attr := range attributes {
		value := w.character.Attributes.GetValue(attr)
		for v := 9; v <= value; v++ {
			usedPoints += w.getAttributeCost(v)
		}
	}

	return totalPoints - usedPoints
}

// respondError responde com uma mensagem de erro
func (w *CharacterWizard) respondError(i *discordgo.InteractionCreate, message string) error {
	_, err := w.session.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "❌ " + message,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	return err
}

// updateAttributeInterface atualiza a interface de seleção de atributos
func (w *CharacterWizard) updateAttributeInterface(i *discordgo.InteractionCreate) error {
	remainingPoints := w.getRemainingPoints()

	// Criar embed atualizado
	embed := &discordgo.MessageEmbed{
		Title:       "📊 Distribuição de Atributos",
		Description: fmt.Sprintf("Pontos restantes: %d", remainingPoints),
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Adicionar campos para cada atributo
	attributes := []struct {
		name    string
		emoji   string
		getFunc func() int
	}{
		{"Força", "💪", w.character.Attributes.GetStrength},
		{"Destreza", "🏃", w.character.Attributes.GetDexterity},
		{"Constituição", "❤️", w.character.Attributes.GetConstitution},
		{"Inteligência", "🧠", w.character.Attributes.GetIntelligence},
		{"Sabedoria", "🦉", w.character.Attributes.GetWisdom},
		{"Carisma", "👑", w.character.Attributes.GetCharisma},
	}

	for _, attr := range attributes {
		value := attr.getFunc()
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s %s", attr.emoji, attr.name),
			Value:  fmt.Sprintf("%d", value),
			Inline: true,
		})
	}

	// Criar botões de ajuste
	components := make([]discordgo.MessageComponent, 0)

	// Botões para cada atributo
	attributeButtons := []struct {
		name string
		id   string
	}{
		{"Força", "str"},
		{"Destreza", "dex"},
		{"Constituição", "con"},
		{"Inteligência", "int"},
		{"Sabedoria", "wis"},
		{"Carisma", "cha"},
	}

	for _, attr := range attributeButtons {
		row := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "+" + attr.name,
					Style:    discordgo.SuccessButton,
					CustomID: "attr_" + attr.id + "_up",
				},
				discordgo.Button{
					Label:    "-" + attr.name,
					Style:    discordgo.DangerButton,
					CustomID: "attr_" + attr.id + "_down",
				},
			},
		}
		components = append(components, row)
	}

	// Botão de confirmação
	if remainingPoints == 0 {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Confirmar Atributos",
					Style:    discordgo.SuccessButton,
					CustomID: "attr_confirm",
					Emoji: discordgo.ComponentEmoji{
						Name: "✅",
					},
				},
			},
		})
	}

	// Atualizar mensagem
	return w.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// handleAttributeConfirmation processa a confirmação dos atributos
func (w *CharacterWizard) handleAttributeConfirmation(i *discordgo.InteractionCreate) error {
	if w.getRemainingPoints() > 0 {
		return w.respondError(i, "Você ainda tem pontos para distribuir")
	}

	// Criar embed de origem
	embed := &discordgo.MessageEmbed{
		Title:       "🌎 Origem do Personagem",
		Description: "Escolha a origem do seu personagem:",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "🏰 Nobre",
				Value: "Nascido em uma família nobre, com acesso a educação e recursos",
			},
			{
				Name:  "🏘️ Plebeu",
				Value: "Origem humilde, mas com forte determinação",
			},
			{
				Name:  "🌲 Selvagem",
				Value: "Criado na natureza, longe da civilização",
			},
		},
	}

	// Criar botões de seleção
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Nobre",
					Style:    discordgo.PrimaryButton,
					CustomID: "background_noble",
					Emoji: discordgo.ComponentEmoji{
						Name: "🏰",
					},
				},
				discordgo.Button{
					Label:    "Plebeu",
					Style:    discordgo.PrimaryButton,
					CustomID: "background_commoner",
					Emoji: discordgo.ComponentEmoji{
						Name: "🏘️",
					},
				},
				discordgo.Button{
					Label:    "Selvagem",
					Style:    discordgo.PrimaryButton,
					CustomID: "background_wild",
					Emoji: discordgo.ComponentEmoji{
						Name: "🌲",
					},
				},
			},
		},
	}

	// Atualizar mensagem
	err := w.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

	if err == nil {
		w.currentStep = StepBackground
	}
	return err
}

// handleBackgroundSelection processa a seleção de origem do personagem
func (w *CharacterWizard) handleBackgroundSelection(i *discordgo.InteractionCreate) error {
	data := i.MessageComponentData().CustomID

	// Extrair origem selecionada
	var background string
	switch data {
	case "background_noble":
		background = "Nobre"
	case "background_commoner":
		background = "Plebeu"
	case "background_wild":
		background = "Selvagem"
	default:
		return fmt.Errorf("origem inválida")
	}

	// Atualizar origem do personagem
	w.character.Background = background
	w.currentStep = StepSkills

	// Inicializar perícias do personagem
	w.character.InitializeSkills()

	// Mostrar tela de seleção de perícias
	return w.showSkillSelection(i)
}

// showSkillSelection mostra a interface de seleção de perícias
func (w *CharacterWizard) showSkillSelection(i *discordgo.InteractionCreate) error {
	// Obter perícias disponíveis para a classe
	availableSkills := gamedata.GetSkillsForClass(w.character.Class)

	// Criar embed com informações
	embed := &discordgo.MessageEmbed{
		Title:       "🎯 Seleção de Perícias",
		Description: "Escolha as perícias em que seu personagem será proficiente:",
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Agrupar perícias por atributo base
	skillsByAttr := make(map[string][]gamedata.Skill)
	for _, skill := range availableSkills {
		baseAttr := gamedata.SkillBaseAttribute[skill]
		skillsByAttr[baseAttr] = append(skillsByAttr[baseAttr], skill)
	}

	// Adicionar campos para cada grupo de perícias
	attrEmojis := map[string]string{
		"strength":     "💪",
		"dexterity":    "🏃",
		"constitution": "❤️",
		"intelligence": "🧠",
		"wisdom":       "🦉",
		"charisma":     "👑",
	}

	for attr, skills := range skillsByAttr {
		skillList := ""
		for _, skill := range skills {
			// Verificar se a perícia já está selecionada
			isProficient := false
			for _, prof := range w.character.Skills {
				if prof.Skill == skill && prof.IsProficient {
					isProficient = true
					break
				}
			}
			status := "⚪" // Não selecionada
			if isProficient {
				status = "🟢" // Selecionada
			}
			skillList += fmt.Sprintf("%s %s\n", status, skill)
		}

		if skillList != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s Perícias de %s", attrEmojis[attr], strings.Title(attr)),
				Value:  skillList,
				Inline: true,
			})
		}
	}

	// Criar botões para cada perícia disponível
	components := make([]discordgo.MessageComponent, 0)
	currentRow := make([]discordgo.MessageComponent, 0)

	for _, skill := range availableSkills {
		// Verificar se a perícia já está selecionada
		isProficient := false
		for _, prof := range w.character.Skills {
			if prof.Skill == skill && prof.IsProficient {
				isProficient = true
				break
			}
		}

		button := discordgo.Button{
			Label:    string(skill),
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("skill_%s", skill),
			Emoji: discordgo.ComponentEmoji{
				Name: "⚪",
			},
		}

		if isProficient {
			button.Style = discordgo.SuccessButton
			button.Emoji.Name = "🟢"
		}

		currentRow = append(currentRow, button)

		// Criar nova linha a cada 3 botões
		if len(currentRow) == 3 {
			components = append(components, discordgo.ActionsRow{
				Components: currentRow,
			})
			currentRow = make([]discordgo.MessageComponent, 0)
		}
	}

	// Adicionar última linha se houver botões restantes
	if len(currentRow) > 0 {
		components = append(components, discordgo.ActionsRow{
			Components: currentRow,
		})
	}

	// Adicionar botão de confirmação
	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Confirmar Perícias",
				Style:    discordgo.SuccessButton,
				CustomID: "skills_confirm",
				Emoji: discordgo.ComponentEmoji{
					Name: "✅",
				},
			},
		},
	})

	// Atualizar mensagem
	return w.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// handleSkillSelection processa a seleção de perícias
func (w *CharacterWizard) handleSkillSelection(i *discordgo.InteractionCreate) error {
	data := i.MessageComponentData().CustomID

	// Se for confirmação, prosseguir para criação do personagem
	if data == "skills_confirm" {
		return w.handleCharacterConfirmation(i)
	}

	// Extrair nome da perícia do ID do botão
	skillName := strings.TrimPrefix(data, "skill_")
	skill := gamedata.Skill(skillName)

	// Verificar se é uma perícia válida para a classe
	if !gamedata.ValidateSkillProficiency(w.character.Class, skill) {
		return w.respondError(i, "Perícia inválida para sua classe")
	}

	// Alternar proficiência na perícia
	found := false
	for i, prof := range w.character.Skills {
		if prof.Skill == skill {
			w.character.Skills[i].IsProficient = !w.character.Skills[i].IsProficient
			found = true
			break
		}
	}

	if !found {
		w.character.Skills = append(w.character.Skills, gamedata.SkillProficiency{
			Skill:        skill,
			IsProficient: true,
		})
	}

	// Atualizar interface
	return w.showSkillSelection(i)
}

// handleClassSelection processa a seleção de classe do personagem
func (w *CharacterWizard) handleClassSelection(data string) error {
	var class gamedata.CharacterClass
	switch data {
	case "class_warrior":
		class = gamedata.Warrior
	case "class_ranger":
		class = gamedata.Ranger
	case "class_mage":
		class = gamedata.Mage
	case "class_paladin":
		class = gamedata.Paladin
	case "class_druid":
		class = gamedata.Druid
	case "class_cleric":
		class = gamedata.Cleric
	case "class_bard":
		class = gamedata.Bard
	case "class_warlock":
		class = gamedata.Warlock
	case "class_sorcerer":
		class = gamedata.Sorcerer
	case "class_rogue":
		class = gamedata.Rogue
	case "class_monk":
		class = gamedata.Monk
	case "class_barbarian":
		class = gamedata.Barbarian
	default:
		return fmt.Errorf("classe inválida: %s", data)
	}

	// Atualiza a classe e inicializa atributos base
	w.character.Class = class
	w.character.Attributes = gamedata.GetBaseAttributesForClass(class)
	w.currentStep = StepBackground
	return nil
}

// handleCharacterConfirmation processa a confirmação final do personagem
func (w *CharacterWizard) handleCharacterConfirmation(i *discordgo.InteractionCreate) error {
	// Criar o personagem no banco
	character, err := w.characterService.CreateCharacter(
		context.Background(),
		w.character.UserID,
		w.character.GuildID,
		w.character.Name,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar personagem: %w", err)
	}

	// Atualizar os dados do personagem
	character.Class = w.character.Class
	character.Background = w.character.Background
	character.Attributes = w.character.Attributes
	character.Skills = w.character.Skills

	// Atualizar o personagem no banco
	if err := w.characterService.UpdateCharacter(context.Background(), character); err != nil {
		return fmt.Errorf("erro ao atualizar personagem: %w", err)
	}

	// Atualiza o personagem no wizard com o ID gerado
	w.character = character

	// Criar embed de sucesso
	embed := &discordgo.MessageEmbed{
		Title:       "✨ Personagem Criado com Sucesso!",
		Description: fmt.Sprintf("**%s** está pronto para a aventura!", character.Name),
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Origem",
				Value:  character.Background,
				Inline: true,
			},
			{
				Name:   "Classe",
				Value:  string(character.Class),
				Inline: true,
			},
			{
				Name:   "Nível",
				Value:  fmt.Sprintf("%d", character.Level),
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
				Name: "Perícias",
				Value: func() string {
					proficientSkills := make([]string, 0)
					for _, skill := range character.Skills {
						if skill.IsProficient {
							proficientSkills = append(proficientSkills, string(skill.Skill))
						}
					}
					if len(proficientSkills) == 0 {
						return "Nenhuma perícia selecionada"
					}
					return strings.Join(proficientSkills, "\n")
				}(),
				Inline: false,
			},
		},
	}

	// Atualiza a mensagem usando ChannelMessageEditComplex
	_, err = w.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    w.channelID,
		ID:         w.messageID,
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar mensagem: %w", err)
	}

	// Remove o wizard do registro após concluir
	return nil
}

// startCharacterCreation inicia o processo de criação de personagem
func (w *CharacterWizard) startCharacterCreation(i *discordgo.InteractionCreate) error {
	// Criar embed inicial
	embed := &discordgo.MessageEmbed{
		Title:       "📊 Distribuição de Atributos",
		Description: "Distribua os pontos entre os atributos do seu personagem:",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "💪 Força",
				Value:  "8",
				Inline: true,
			},
			{
				Name:   "🏃 Destreza",
				Value:  "8",
				Inline: true,
			},
			{
				Name:   "❤️ Constituição",
				Value:  "8",
				Inline: true,
			},
			{
				Name:   "🧠 Inteligência",
				Value:  "8",
				Inline: true,
			},
			{
				Name:   "🦉 Sabedoria",
				Value:  "8",
				Inline: true,
			},
			{
				Name:   "👑 Carisma",
				Value:  "8",
				Inline: true,
			},
		},
	}

	// Criar botões de ajuste
	components := make([]discordgo.MessageComponent, 0)

	// Botões para cada atributo
	attributeButtons := []struct {
		name string
		id   string
	}{
		{"Força", "str"},
		{"Destreza", "dex"},
		{"Constituição", "con"},
		{"Inteligência", "int"},
		{"Sabedoria", "wis"},
		{"Carisma", "cha"},
	}

	for _, attr := range attributeButtons {
		row := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "+" + attr.name,
					Style:    discordgo.SuccessButton,
					CustomID: "attr_" + attr.id + "_up",
				},
				discordgo.Button{
					Label:    "-" + attr.name,
					Style:    discordgo.DangerButton,
					CustomID: "attr_" + attr.id + "_down",
				},
			},
		}
		components = append(components, row)
	}

	// Atualizar mensagem
	return w.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// HandleInteraction processa interações com o wizard
func (w *CharacterWizard) HandleInteraction(i *discordgo.InteractionCreate) error {
	// Verifica se a interação é para a mensagem correta
	if i.Message.ID != w.messageID || i.ChannelID != w.channelID {
		return fmt.Errorf("interação inválida para este wizard")
	}

	data := i.MessageComponentData().CustomID

	// Processa a seleção de classe
	if strings.HasPrefix(data, "class_") {
		if err := w.handleClassSelection(data); err != nil {
			return w.respondError(i, fmt.Sprintf("Erro ao selecionar classe: %s", err))
		}

		// Após selecionar a classe, mostra a tela de origem
		embed := &discordgo.MessageEmbed{
			Title:       "🌎 Origem do Personagem",
			Description: "Escolha a origem do seu personagem:",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "🏰 Nobre",
					Value: "Nascido em uma família nobre, com acesso a educação e recursos",
				},
				{
					Name:  "🏘️ Plebeu",
					Value: "Origem humilde, mas com forte determinação",
				},
				{
					Name:  "🌲 Selvagem",
					Value: "Criado na natureza, longe da civilização",
				},
			},
		}

		// Criar botões de seleção
		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Nobre",
						Style:    discordgo.PrimaryButton,
						CustomID: "background_noble",
						Emoji: discordgo.ComponentEmoji{
							Name: "🏰",
						},
					},
					discordgo.Button{
						Label:    "Plebeu",
						Style:    discordgo.PrimaryButton,
						CustomID: "background_commoner",
						Emoji: discordgo.ComponentEmoji{
							Name: "🏘️",
						},
					},
					discordgo.Button{
						Label:    "Selvagem",
						Style:    discordgo.PrimaryButton,
						CustomID: "background_wild",
						Emoji: discordgo.ComponentEmoji{
							Name: "🌲",
						},
					},
				},
			},
		}

		// Responde à interação com a nova mensagem
		return w.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
			},
		})
	}

	// Processa a seleção de origem
	if strings.HasPrefix(data, "background_") {
		return w.handleBackgroundSelection(i)
	}

	// Processa a seleção de perícias
	if strings.HasPrefix(data, "skill_") || data == "skills_confirm" {
		return w.handleSkillSelection(i)
	}

	// Processa ajustes de atributos
	if strings.HasPrefix(data, "attr_") {
		return w.handleAttributeSelection(i)
	}

	// Processa confirmação final
	if data == "confirm" || data == "restart" {
		return w.handleCharacterConfirmation(i)
	}

	return fmt.Errorf("ação inválida: %s", data)
}

func createClassButtons() []discordgo.MessageComponent {
	// Primeira linha: Classes marciais básicas
	row1 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Guerreiro",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_warrior",
			},
			discordgo.Button{
				Label:    "Bárbaro",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_barbarian",
			},
			discordgo.Button{
				Label:    "Monge",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_monk",
			},
			discordgo.Button{
				Label:    "Paladino",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_paladin",
			},
		},
	}

	// Segunda linha: Classes de furtividade e destreza
	row2 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Arqueiro",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_ranger",
			},
			discordgo.Button{
				Label:    "Ladino",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_rogue",
			},
			discordgo.Button{
				Label:    "Bardo",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_bard",
			},
		},
	}

	// Terceira linha: Classes mágicas
	row3 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Mago",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_mage",
			},
			discordgo.Button{
				Label:    "Bruxo",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_warlock",
			},
			discordgo.Button{
				Label:    "Feiticeiro",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_sorcerer",
			},
		},
	}

	// Quarta linha: Classes divinas e naturais
	row4 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Clérigo",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_cleric",
			},
			discordgo.Button{
				Label:    "Druida",
				Style:    discordgo.PrimaryButton,
				CustomID: "class_druid",
			},
		},
	}

	return []discordgo.MessageComponent{row1, row2, row3, row4}
}

// ... existing code ...
