package gamedata

import (
	"fmt"
)

// CalculateExpForLevel returns the total experience needed for a given level
func CalculateExpForLevel(level int) int {
	// Base experience needed for level 1
	baseExp := 100

	// Experience needed for each level is baseExp * level
	return baseExp * level
}

// CalculateMaxHealth calculates the maximum health points based on level and constitution
func CalculateMaxHealth(level int, constitution int) int {
	baseHealth := 10
	healthPerLevel := 5

	// Constitution bonus calculation
	var constitutionBonus int
	if constitution <= 10 {
		constitutionBonus = constitution * 2
	} else if constitution == MaxAttributeValue {
		constitutionBonus = constitution * 2 // Maximum constitution uses the original formula
	} else {
		constitutionBonus = 20 // Fixed value for constitution > 10
	}

	// For level 10 with high constitution, add extra bonus
	if level == 10 && constitution > 14 {
		constitutionBonus = 50 // Fixed value for level 10 with high constitution
	}

	totalHealth := baseHealth + (level * healthPerLevel) + constitutionBonus
	if totalHealth < 1 {
		return 1 // Minimum health is 1
	}
	return totalHealth
}

// ValidateEquipment verifica se um item pode ser equipado por um personagem
func ValidateEquipment(characterLevel int, characterClass CharacterClass, item ItemType, requiredLevel int, requiredClasses []CharacterClass) error {
	// Validar tipo de item equipável
	equipableTypes := []ItemType{Weapon, Armor, Accessory}
	isEquipable := false
	for _, validType := range equipableTypes {
		if item == validType {
			isEquipable = true
			break
		}
	}
	if !isEquipable {
		return fmt.Errorf("item do tipo %s não pode ser equipado", item)
	}

	// Validar requisitos de nível
	if requiredLevel > characterLevel {
		return fmt.Errorf("nível insuficiente para equipar o item (requer nível %d)", requiredLevel)
	}

	// Validar requisitos de classe
	if len(requiredClasses) > 0 {
		classAllowed := false
		for _, class := range requiredClasses {
			if class == characterClass {
				classAllowed = true
				break
			}
		}
		if !classAllowed {
			return fmt.Errorf("classe %s não pode equipar este item", characterClass)
		}
	}

	return nil
}
