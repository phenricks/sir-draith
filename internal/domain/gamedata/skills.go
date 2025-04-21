package gamedata

// Skill representa uma perícia
type Skill string

const (
	// Perícias baseadas em Força
	Athletics Skill = "athletics" // Atletismo

	// Perícias baseadas em Destreza
	Acrobatics    Skill = "acrobatics"    // Acrobacia
	SleightOfHand Skill = "sleightOfHand" // Prestidigitação
	Stealth       Skill = "stealth"       // Furtividade

	// Perícias baseadas em Inteligência
	Arcana        Skill = "arcana"        // Arcana
	History       Skill = "history"       // História
	Investigation Skill = "investigation" // Investigação
	Nature        Skill = "nature"        // Natureza
	Religion      Skill = "religion"      // Religião

	// Perícias baseadas em Sabedoria
	AnimalHandling Skill = "animalHandling" // Adestramento
	Insight        Skill = "insight"        // Intuição
	Medicine       Skill = "medicine"       // Medicina
	Perception     Skill = "perception"     // Percepção
	Survival       Skill = "survival"       // Sobrevivência

	// Perícias baseadas em Carisma
	Deception    Skill = "deception"    // Enganação
	Intimidation Skill = "intimidation" // Intimidação
	Performance  Skill = "performance"  // Atuação
	Persuasion   Skill = "persuasion"   // Persuasão
)

// SkillBaseAttribute mapeia cada perícia ao seu atributo base
var SkillBaseAttribute = map[Skill]string{
	// Força
	Athletics: "strength",

	// Destreza
	Acrobatics:    "dexterity",
	SleightOfHand: "dexterity",
	Stealth:       "dexterity",

	// Inteligência
	Arcana:        "intelligence",
	History:       "intelligence",
	Investigation: "intelligence",
	Nature:        "intelligence",
	Religion:      "intelligence",

	// Sabedoria
	AnimalHandling: "wisdom",
	Insight:        "wisdom",
	Medicine:       "wisdom",
	Perception:     "wisdom",
	Survival:       "wisdom",

	// Carisma
	Deception:    "charisma",
	Intimidation: "charisma",
	Performance:  "charisma",
	Persuasion:   "charisma",
}

// ClassSkillProficiencies define as perícias em que cada classe tem proficiência por padrão
var ClassSkillProficiencies = map[CharacterClass][]Skill{
	Warrior: {
		Athletics,
		Intimidation,
		Perception,
		Survival,
	},
	Mage: {
		Arcana,
		History,
		Investigation,
		Medicine,
	},
	Ranger: {
		AnimalHandling,
		Athletics,
		Nature,
		Stealth,
		Survival,
	},
	Cleric: {
		Insight,
		Medicine,
		Persuasion,
		Religion,
	},
	Paladin: {
		Athletics,
		Intimidation,
		Medicine,
		Persuasion,
	},
	Druid: {
		AnimalHandling,
		Nature,
		Medicine,
		Survival,
	},
	Barbarian: {
		Athletics,
		Intimidation,
		Nature,
		Survival,
	},
	Monk: {
		Acrobatics,
		Athletics,
		Stealth,
		Insight,
	},
	Bard: {
		Deception,
		Performance,
		Persuasion,
		SleightOfHand,
	},
	Warlock: {
		Arcana,
		Deception,
		Intimidation,
		Persuasion,
	},
	Sorcerer: {
		Arcana,
		Deception,
		Intimidation,
		Persuasion,
	},
	Rogue: {
		Acrobatics,
		Deception,
		SleightOfHand,
		Stealth,
	},
}

// SkillProficiency representa a proficiência em uma perícia
type SkillProficiency struct {
	Skill        Skill `bson:"skill" json:"skill"`               // Nome da perícia
	IsProficient bool  `bson:"isProficient" json:"isProficient"` // Se tem proficiência
	Bonus        int   `bson:"bonus" json:"bonus"`               // Bônus adicional
}

// CalculateSkillModifier calcula o modificador total de uma perícia
func CalculateSkillModifier(skill Skill, attributes *Attributes, proficiency *SkillProficiency, level int) int {
	// Obter o atributo base da perícia
	baseAttr := SkillBaseAttribute[skill]
	if baseAttr == "" {
		return 0
	}

	// Calcular o modificador do atributo base
	attrValue := attributes.GetValue(baseAttr)
	attrMod := (attrValue - 10) / 2

	// Adicionar bônus de proficiência se aplicável
	if proficiency != nil && proficiency.IsProficient {
		// Bônus de proficiência base é 2 + (level / 4)
		profBonus := 2 + (level / 4)
		attrMod += profBonus
	}

	// Adicionar bônus extras da perícia
	if proficiency != nil {
		attrMod += proficiency.Bonus
	}

	return attrMod
}

// GetSkillsForClass retorna as perícias disponíveis para uma classe
func GetSkillsForClass(class CharacterClass) []Skill {
	skills, exists := ClassSkillProficiencies[class]
	if !exists {
		return []Skill{}
	}
	return skills
}

// ValidateSkillProficiency verifica se uma perícia pode ser usada por uma classe
func ValidateSkillProficiency(class CharacterClass, skill Skill) bool {
	skills := GetSkillsForClass(class)
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	return false
}
