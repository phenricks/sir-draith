package gamedata

// Attributes representa os atributos base de um personagem
type Attributes struct {
	Strength     int `bson:"strength" json:"strength"`         // Força
	Dexterity    int `bson:"dexterity" json:"dexterity"`       // Destreza
	Constitution int `bson:"constitution" json:"constitution"` // Constituição
	Intelligence int `bson:"intelligence" json:"intelligence"` // Inteligência
	Wisdom       int `bson:"wisdom" json:"wisdom"`             // Sabedoria
	Charisma     int `bson:"charisma" json:"charisma"`         // Carisma
}

// ItemStats representa os atributos de um item
type ItemStats struct {
	Attack     int `bson:"attack" json:"attack"`         // Ataque
	Defense    int `bson:"defense" json:"defense"`       // Defesa
	MagicPower int `bson:"magicPower" json:"magicPower"` // Poder Mágico
}

// CharacterClass representa as classes disponíveis para um personagem
type CharacterClass string

const (
	Warrior   CharacterClass = "warrior"
	Mage      CharacterClass = "mage"
	Rogue     CharacterClass = "rogue"
	Cleric    CharacterClass = "cleric"
	Ranger    CharacterClass = "ranger"
	Paladin   CharacterClass = "paladin"
	Druid     CharacterClass = "druid"
	Barbarian CharacterClass = "barbarian"
	Monk      CharacterClass = "monk"
	Bard      CharacterClass = "bard"
	Warlock   CharacterClass = "warlock"
	Sorcerer  CharacterClass = "sorcerer"
)

// ItemType representa os tipos de itens disponíveis
type ItemType string

const (
	Weapon     ItemType = "weapon"
	Armor      ItemType = "armor"
	Accessory  ItemType = "accessory"
	Consumable ItemType = "consumable"
	Quest      ItemType = "quest"
)

// ItemRarity representa as raridades disponíveis para itens
type ItemRarity string

const (
	Common    ItemRarity = "common"
	Uncommon  ItemRarity = "uncommon"
	Rare      ItemRarity = "rare"
	Epic      ItemRarity = "epic"
	Legendary ItemRarity = "legendary"
	Mythical  ItemRarity = "mythical"
)

// EquipmentSlot representa os slots de equipamento disponíveis
type EquipmentSlot string

const (
	Head     EquipmentSlot = "head"
	Neck     EquipmentSlot = "neck"
	Chest    EquipmentSlot = "chest"
	Legs     EquipmentSlot = "legs"
	Feet     EquipmentSlot = "feet"
	MainHand EquipmentSlot = "mainHand"
	OffHand  EquipmentSlot = "offHand"
	Ring1    EquipmentSlot = "ring1"
	Ring2    EquipmentSlot = "ring2"
	Trinket1 EquipmentSlot = "trinket1"
	Trinket2 EquipmentSlot = "trinket2"
)

// GetBaseAttributesForClass retorna os atributos base para cada classe
func GetBaseAttributesForClass(class CharacterClass) Attributes {
	switch class {
	case Warrior:
		return Attributes{
			Strength:     15,
			Dexterity:    12,
			Constitution: 14,
			Intelligence: 8,
			Wisdom:       10,
			Charisma:     10,
		}
	case Ranger:
		return Attributes{
			Strength:     12,
			Dexterity:    15,
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       14,
			Charisma:     8,
		}
	case Mage:
		return Attributes{
			Strength:     8,
			Dexterity:    12,
			Constitution: 10,
			Intelligence: 15,
			Wisdom:       14,
			Charisma:     12,
		}
	case Paladin:
		return Attributes{
			Strength:     14,
			Dexterity:    10,
			Constitution: 13,
			Intelligence: 8,
			Wisdom:       12,
			Charisma:     14,
		}
	case Druid:
		return Attributes{
			Strength:     10,
			Dexterity:    12,
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       15,
			Charisma:     12,
		}
	case Cleric:
		return Attributes{
			Strength:     12,
			Dexterity:    10,
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       15,
			Charisma:     12,
		}
	case Bard:
		return Attributes{
			Strength:     8,
			Dexterity:    12,
			Constitution: 10,
			Intelligence: 12,
			Wisdom:       10,
			Charisma:     15,
		}
	case Warlock:
		return Attributes{
			Strength:     8,
			Dexterity:    12,
			Constitution: 12,
			Intelligence: 13,
			Wisdom:       10,
			Charisma:     15,
		}
	case Sorcerer:
		return Attributes{
			Strength:     8,
			Dexterity:    12,
			Constitution: 12,
			Intelligence: 12,
			Wisdom:       10,
			Charisma:     15,
		}
	case Rogue:
		return Attributes{
			Strength:     10,
			Dexterity:    15,
			Constitution: 12,
			Intelligence: 12,
			Wisdom:       10,
			Charisma:     12,
		}
	case Monk:
		return Attributes{
			Strength:     12,
			Dexterity:    15,
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       13,
			Charisma:     8,
		}
	case Barbarian:
		return Attributes{
			Strength:     15,
			Dexterity:    12,
			Constitution: 14,
			Intelligence: 8,
			Wisdom:       10,
			Charisma:     8,
		}
	default:
		return Attributes{
			Strength:     10,
			Dexterity:    10,
			Constitution: 10,
			Intelligence: 10,
			Wisdom:       10,
			Charisma:     10,
		}
	}
}
