package validation

import (
	"fmt"
	"math"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"
)

// CalculateExpForLevel calcula a experiência necessária para um determinado nível
func CalculateExpForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	if level > gamedata.MaxLevel {
		return -1
	}

	// Fórmula: Exp = Base * (1.5 ^ (Level - 1))
	baseExp := 100.0
	multiplier := math.Pow(1.5, float64(level-1))
	return int(baseExp * multiplier)
}

// CalculateMaxHealth calcula o HP máximo do personagem baseado no nível e constituição
func CalculateMaxHealth(level int, constitution int) int {
	if level < 1 || level > gamedata.MaxLevel {
		return 0
	}

	// Base HP por nível
	baseHP := 10

	// Bônus de constituição (começa em -1 para constituição 10)
	constBonus := (constitution - 10) / 2

	// HP total = (Base + Bônus Const) * Nível
	totalHP := (baseHP + constBonus) * level

	// Garantir mínimo de 1 HP
	if totalHP < 1 {
		return 1
	}

	return totalHP
}

// CalculateAttackBonus calcula o bônus de ataque baseado na classe e atributos
func CalculateAttackBonus(class gamedata.CharacterClass, level int, attributes *gamedata.Attributes) int {
	if attributes == nil {
		return 0
	}

	// Bônus base por nível (proficiency bonus)
	proficiencyBonus := (level-1)/4 + 2

	// Bônus de atributo baseado na classe
	var attrBonus int
	switch class {
	case gamedata.Warrior:
		attrBonus = (attributes.Strength - 10) / 2
	case gamedata.Ranger:
		attrBonus = (attributes.Dexterity - 10) / 2
	case gamedata.Mage:
		attrBonus = (attributes.Intelligence - 10) / 2
	default:
		return 0
	}

	return proficiencyBonus + attrBonus
}

// CalculateDefense calcula a defesa base do personagem
func CalculateDefense(class gamedata.CharacterClass, attributes *gamedata.Attributes) int {
	if attributes == nil {
		return 10
	}

	// Defesa base é 10
	baseDefense := 10

	// Bônus de atributo baseado na classe
	var attrBonus int
	switch class {
	case gamedata.Warrior:
		// Guerreiros usam força ou destreza (o maior) para defesa
		strBonus := (attributes.Strength - 10) / 2
		dexBonus := (attributes.Dexterity - 10) / 2
		if strBonus > dexBonus {
			attrBonus = strBonus
		} else {
			attrBonus = dexBonus
		}
	case gamedata.Ranger:
		// Arqueiros usam destreza
		attrBonus = (attributes.Dexterity - 10) / 2
	case gamedata.Mage:
		// Magos usam inteligência
		attrBonus = (attributes.Intelligence - 10) / 2
	}

	return baseDefense + attrBonus
}

// ValidateCharacterProgression valida a progressão do personagem
func ValidateCharacterProgression(character *entities.Character) error {
	if character == nil {
		return fmt.Errorf("personagem não pode ser nulo")
	}

	// Validar nível
	if character.Level < gamedata.StartingLevel || character.Level > gamedata.MaxLevel {
		return fmt.Errorf("nível inválido: deve estar entre %d e %d", gamedata.StartingLevel, gamedata.MaxLevel)
	}

	// Validar experiência
	minExp := CalculateExpForLevel(character.Level)
	maxExp := CalculateExpForLevel(character.Level + 1)
	if maxExp != -1 && character.Experience >= maxExp {
		return fmt.Errorf("experiência suficiente para subir de nível")
	}
	if character.Experience < minExp {
		return fmt.Errorf("experiência insuficiente para o nível atual")
	}

	// Validar HP máximo
	expectedHP := CalculateMaxHealth(character.Level, character.Attributes.Constitution)
	if character.Combat.MaxHealth != expectedHP {
		return fmt.Errorf("HP máximo incorreto: esperado %d, atual %d", expectedHP, character.Combat.MaxHealth)
	}

	// Validar HP atual
	if character.Combat.Health > character.Combat.MaxHealth {
		return fmt.Errorf("HP atual não pode ser maior que o máximo")
	}
	if character.Combat.Health < 0 {
		return fmt.Errorf("HP atual não pode ser negativo")
	}

	return nil
}
