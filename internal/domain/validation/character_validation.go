package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"
)

const (
	// Limites de personagem
	MinNameLength     = 3
	MaxNameLength     = 32
	MaxLevel          = 20
	MinAttributeValue = 3
	MaxAttributeValue = 20
	MaxInventorySize  = 50
	MaxEquipmentSize  = 10

	// Limites de item
	MinItemNameLength = 3
	MaxItemNameLength = 64
	MaxItemQuantity   = 99
	MinItemValue      = 0
	MaxItemValue      = 100000

	// Experiência e Gold
	BaseExpPerLevel = 100 // Experiência base para o primeiro nível
	ExpMultiplier   = 1.5 // Multiplicador de experiência por nível
	StartingGold    = 100 // Ouro inicial
	MaxGold         = 1000000
)

var (
	// Expressões regulares para validação
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-']+$`)

	// Raridades de item e seus multiplicadores
	rarityMultipliers = map[string]float64{
		"comum":    1.0,
		"incomum":  1.5,
		"raro":     2.0,
		"épico":    3.0,
		"lendário": 5.0,
	}

	// Slots de equipamento válidos
	validEquipmentSlots = map[string]bool{
		"Mão Principal":  true,
		"Mão Secundária": true,
		"Cabeça":         true,
		"Corpo":          true,
		"Mãos":           true,
		"Pés":            true,
		"Anel":           true,
		"Amuleto":        true,
	}

	// Tipos de item válidos
	validItemTypes = map[string]bool{
		"Arma":       true,
		"Armadura":   true,
		"Consumível": true,
		"Recurso":    true,
		"Tesouro":    true,
	}
)

// ValidateNewCharacter valida os dados de um novo personagem
func ValidateNewCharacter(character *entities.Character) error {
	if character == nil {
		return fmt.Errorf("personagem não pode ser nulo")
	}

	// Validar nome
	if err := validateName(character.Name); err != nil {
		return err
	}

	// Validar classe
	if err := validateClass(character.Class); err != nil {
		return err
	}

	// Validar atributos
	if err := ValidateAttributes(character.Attributes); err != nil {
		return err
	}

	// Validar nível inicial
	if character.Level != gamedata.StartingLevel {
		return fmt.Errorf("nível inicial deve ser %d", gamedata.StartingLevel)
	}

	// Validar experiência inicial
	if character.Experience != 0 {
		return fmt.Errorf("experiência inicial deve ser 0")
	}

	// Validar ouro inicial
	if character.Gold != gamedata.StartingGold {
		return fmt.Errorf("ouro inicial deve ser %d", gamedata.StartingGold)
	}

	// Validar datas
	now := time.Now()
	if character.CreatedAt.After(now) {
		return fmt.Errorf("data de criação inválida")
	}

	return nil
}

// validateName valida o nome do personagem
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < gamedata.MinNameLength {
		return fmt.Errorf("nome muito curto (mínimo %d caracteres)", gamedata.MinNameLength)
	}
	if len(name) > gamedata.MaxNameLength {
		return fmt.Errorf("nome muito longo (máximo %d caracteres)", gamedata.MaxNameLength)
	}
	return nil
}

// validateClass valida a classe do personagem
func validateClass(class string) error {
	class = strings.TrimSpace(strings.ToLower(class))
	validClasses := []string{"guerreiro", "mago", "arqueiro"}
	for _, validClass := range validClasses {
		if class == validClass {
			return nil
		}
	}
	return fmt.Errorf("classe inválida: %s", class)
}

// ValidateAttributes valida os atributos do personagem
func ValidateAttributes(attrs *entities.Attributes) error {
	if attrs == nil {
		return fmt.Errorf("atributos não podem ser nulos")
	}

	// Validar valores individuais
	attributes := []int{
		attrs.Strength,
		attrs.Dexterity,
		attrs.Constitution,
		attrs.Intelligence,
		attrs.Wisdom,
		attrs.Charisma,
	}

	for _, value := range attributes {
		if value < gamedata.MinAttributeValue {
			return fmt.Errorf("atributo não pode ser menor que %d", gamedata.MinAttributeValue)
		}
		if value > gamedata.MaxAttributeValue {
			return fmt.Errorf("atributo não pode ser maior que %d", gamedata.MaxAttributeValue)
		}
	}

	// Validar soma total
	total := 0
	for _, value := range attributes {
		total += value
	}
	if total > gamedata.AttributePoints {
		return fmt.Errorf("total de pontos de atributo excede o limite de %d", gamedata.AttributePoints)
	}

	return nil
}

// ValidateItem valida os dados de um item
func ValidateItem(item *entities.Item) error {
	if item == nil {
		return fmt.Errorf("item não pode ser nulo")
	}

	// Validar nome
	name := strings.TrimSpace(item.Name)
	if len(name) < gamedata.MinItemNameLength {
		return fmt.Errorf("nome do item muito curto (mínimo %d caracteres)", gamedata.MinItemNameLength)
	}
	if len(name) > gamedata.MaxItemNameLength {
		return fmt.Errorf("nome do item muito longo (máximo %d caracteres)", gamedata.MaxItemNameLength)
	}

	// Validar descrição
	description := strings.TrimSpace(item.Description)
	if len(description) < gamedata.MinDescriptionLength {
		return fmt.Errorf("descrição do item muito curta (mínimo %d caracteres)", gamedata.MinDescriptionLength)
	}
	if len(description) > gamedata.MaxDescriptionLength {
		return fmt.Errorf("descrição do item muito longa (máximo %d caracteres)", gamedata.MaxDescriptionLength)
	}

	// Validar tipo
	if err := validateItemType(item.Type); err != nil {
		return err
	}

	// Validar raridade
	if err := validateRarity(item.Rarity); err != nil {
		return err
	}

	// Validar atributos do item
	if item.Stats != nil {
		if err := validateItemStats(item.Stats); err != nil {
			return err
		}
	}

	return nil
}

// validateItemType valida o tipo do item
func validateItemType(itemType string) error {
	itemType = strings.TrimSpace(strings.ToLower(itemType))
	validTypes := []string{
		"arma", "armadura", "escudo", "anel",
		"amuleto", "poção", "pergaminho", "outro",
	}
	for _, validType := range validTypes {
		if itemType == validType {
			return nil
		}
	}
	return fmt.Errorf("tipo de item inválido: %s", itemType)
}

// validateRarity valida a raridade do item
func validateRarity(rarity string) error {
	rarity = strings.TrimSpace(strings.ToLower(rarity))
	validRarities := []string{
		"comum", "incomum", "raro",
		"muito raro", "lendário", "único",
	}
	for _, validRarity := range validRarities {
		if rarity == validRarity {
			return nil
		}
	}
	return fmt.Errorf("raridade inválida: %s", rarity)
}

// validateItemStats valida os atributos de um item
func validateItemStats(stats *entities.ItemStats) error {
	if stats == nil {
		return nil
	}

	attributes := []int{
		stats.Strength,
		stats.Dexterity,
		stats.Constitution,
		stats.Intelligence,
		stats.Wisdom,
		stats.Charisma,
	}

	for _, value := range attributes {
		if value < gamedata.MinItemAttributeValue {
			return fmt.Errorf("bônus de atributo do item não pode ser menor que %d", gamedata.MinItemAttributeValue)
		}
		if value > gamedata.MaxItemAttributeValue {
			return fmt.Errorf("bônus de atributo do item não pode ser maior que %d", gamedata.MaxItemAttributeValue)
		}
	}

	return nil
}

// ValidateEquipment valida se um item pode ser equipado por um personagem
func ValidateEquipment(character *entities.Character, item *entities.Item) error {
	if character == nil {
		return fmt.Errorf("personagem não pode ser nulo")
	}
	if item == nil {
		return fmt.Errorf("item não pode ser nulo")
	}

	// Validar tipo de item equipável
	equipableTypes := []string{"arma", "armadura", "escudo", "anel", "amuleto"}
	isEquipable := false
	for _, validType := range equipableTypes {
		if item.Type == validType {
			isEquipable = true
			break
		}
	}
	if !isEquipable {
		return fmt.Errorf("item do tipo %s não pode ser equipado", item.Type)
	}

	// Validar requisitos de nível
	if item.RequiredLevel > character.Level {
		return fmt.Errorf("nível insuficiente para equipar o item (requer nível %d)", item.RequiredLevel)
	}

	// Validar requisitos de classe
	if len(item.RequiredClasses) > 0 {
		classAllowed := false
		for _, class := range item.RequiredClasses {
			if strings.EqualFold(class, character.Class) {
				classAllowed = true
				break
			}
		}
		if !classAllowed {
			return fmt.Errorf("classe %s não pode equipar este item", character.Class)
		}
	}

	return nil
}
