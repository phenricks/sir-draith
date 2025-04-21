package gamedata

import "fmt"

// ClassRequirements define os requisitos mínimos de atributos para cada classe
var ClassRequirements = map[CharacterClass]Attributes{
	Warrior: {
		Strength:     13,
		Constitution: 12,
	},
	Mage: {
		Intelligence: 13,
		Wisdom:       12,
	},
	Ranger: {
		Dexterity:    13,
		Constitution: 12,
	},
	Cleric: {
		Wisdom:       13,
		Constitution: 12,
	},
	Paladin: {
		Strength: 12,
		Charisma: 13,
	},
	Druid: {
		Wisdom:       13,
		Intelligence: 12,
	},
	Barbarian: {
		Strength:     13,
		Constitution: 13,
	},
	Monk: {
		Dexterity: 13,
		Wisdom:    12,
	},
	Bard: {
		Charisma:  13,
		Dexterity: 12,
	},
	Warlock: {
		Charisma:     13,
		Intelligence: 12,
	},
	Sorcerer: {
		Charisma:     13,
		Constitution: 12,
	},
	Rogue: {
		Dexterity:    13,
		Intelligence: 12,
	},
}

// ClassEquipmentRestrictions define as restrições de equipamento para cada classe
var ClassEquipmentRestrictions = map[CharacterClass][]ItemType{
	Warrior:   {Weapon, Armor},
	Mage:      {Weapon, Accessory},
	Ranger:    {Weapon, Armor},
	Cleric:    {Weapon, Armor},
	Paladin:   {Weapon, Armor},
	Druid:     {Weapon, Accessory},
	Barbarian: {Weapon, Armor},
	Monk:      {Weapon, Accessory},
	Bard:      {Weapon, Accessory},
	Warlock:   {Weapon, Accessory},
	Sorcerer:  {Weapon, Accessory},
	Rogue:     {Weapon, Armor},
}

// ValidateClassRequirements verifica se os atributos atendem aos requisitos da classe
func ValidateClassRequirements(class CharacterClass, attributes *Attributes) error {
	requirements, exists := ClassRequirements[class]
	if !exists {
		return fmt.Errorf("classe %s não encontrada", class)
	}

	// Verificar força
	if requirements.Strength > 0 && attributes.Strength < requirements.Strength {
		return fmt.Errorf("força insuficiente para a classe %s (mínimo: %d)", class, requirements.Strength)
	}

	// Verificar destreza
	if requirements.Dexterity > 0 && attributes.Dexterity < requirements.Dexterity {
		return fmt.Errorf("destreza insuficiente para a classe %s (mínimo: %d)", class, requirements.Dexterity)
	}

	// Verificar constituição
	if requirements.Constitution > 0 && attributes.Constitution < requirements.Constitution {
		return fmt.Errorf("constituição insuficiente para a classe %s (mínimo: %d)", class, requirements.Constitution)
	}

	// Verificar inteligência
	if requirements.Intelligence > 0 && attributes.Intelligence < requirements.Intelligence {
		return fmt.Errorf("inteligência insuficiente para a classe %s (mínimo: %d)", class, requirements.Intelligence)
	}

	// Verificar sabedoria
	if requirements.Wisdom > 0 && attributes.Wisdom < requirements.Wisdom {
		return fmt.Errorf("sabedoria insuficiente para a classe %s (mínimo: %d)", class, requirements.Wisdom)
	}

	// Verificar carisma
	if requirements.Charisma > 0 && attributes.Charisma < requirements.Charisma {
		return fmt.Errorf("carisma insuficiente para a classe %s (mínimo: %d)", class, requirements.Charisma)
	}

	return nil
}

// ValidateClassEquipment verifica se um tipo de equipamento pode ser usado pela classe
func ValidateClassEquipment(class CharacterClass, itemType ItemType) error {
	restrictions, exists := ClassEquipmentRestrictions[class]
	if !exists {
		return fmt.Errorf("classe %s não encontrada", class)
	}

	for _, allowedType := range restrictions {
		if allowedType == itemType {
			return nil
		}
	}

	return fmt.Errorf("a classe %s não pode usar equipamentos do tipo %s", class, itemType)
}

// GetClassPrimaryAttribute retorna o atributo principal da classe
func GetClassPrimaryAttribute(class CharacterClass) string {
	switch class {
	case Warrior, Barbarian:
		return "strength"
	case Ranger, Monk, Rogue:
		return "dexterity"
	case Mage:
		return "intelligence"
	case Cleric, Druid:
		return "wisdom"
	case Paladin, Bard, Warlock, Sorcerer:
		return "charisma"
	default:
		return ""
	}
}

// CalculateClassBonus calcula o bônus de classe baseado no atributo principal
func CalculateClassBonus(class CharacterClass, attributes *Attributes) int {
	primaryAttr := GetClassPrimaryAttribute(class)
	if primaryAttr == "" {
		return 0
	}

	attrValue := attributes.GetValue(primaryAttr)
	bonus := (attrValue - 10) / 2 // Bônus padrão baseado no modificador do atributo

	// Bônus específicos por classe
	switch class {
	case Warrior:
		bonus += 2 // Bônus adicional de combate
	case Mage:
		bonus += 1 // Bônus adicional de magia
	case Ranger:
		if attributes.Dexterity >= 15 {
			bonus += 1 // Bônus de precisão
		}
	case Barbarian:
		if attributes.Constitution >= 15 {
			bonus += 1 // Bônus de resistência
		}
	}

	return bonus
}
