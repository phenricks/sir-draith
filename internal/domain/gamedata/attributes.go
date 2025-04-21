package gamedata

// GetValue retorna o valor de um atributo pelo nome
func (a *Attributes) GetValue(attribute string) int {
	switch attribute {
	case "strength":
		return a.Strength
	case "dexterity":
		return a.Dexterity
	case "constitution":
		return a.Constitution
	case "intelligence":
		return a.Intelligence
	case "wisdom":
		return a.Wisdom
	case "charisma":
		return a.Charisma
	default:
		return 0
	}
}

// SetValue define o valor de um atributo pelo nome
func (a *Attributes) SetValue(attribute string, value int) {
	switch attribute {
	case "strength":
		a.Strength = value
	case "dexterity":
		a.Dexterity = value
	case "constitution":
		a.Constitution = value
	case "intelligence":
		a.Intelligence = value
	case "wisdom":
		a.Wisdom = value
	case "charisma":
		a.Charisma = value
	}
}

// GetStrength retorna o valor de Força
func (a *Attributes) GetStrength() int {
	return a.Strength
}

// GetDexterity retorna o valor de Destreza
func (a *Attributes) GetDexterity() int {
	return a.Dexterity
}

// GetConstitution retorna o valor de Constituição
func (a *Attributes) GetConstitution() int {
	return a.Constitution
}

// GetIntelligence retorna o valor de Inteligência
func (a *Attributes) GetIntelligence() int {
	return a.Intelligence
}

// GetWisdom retorna o valor de Sabedoria
func (a *Attributes) GetWisdom() int {
	return a.Wisdom
}

// GetCharisma retorna o valor de Carisma
func (a *Attributes) GetCharisma() int {
	return a.Charisma
}
