package entities

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Character representa um personagem no sistema
type Character struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`         // ID do usuário Discord
	GuildID     string             `bson:"guild_id"`        // ID do servidor Discord
	Name        string             `bson:"name"`            // Nome do personagem
	Title       string             `bson:"title,omitempty"` // Título nobiliárquico
	Description string             `bson:"description"`     // História/Descrição
	Level       int                `bson:"level"`           // Nível do personagem
	Experience  int                `bson:"experience"`      // Experiência atual
	Gold        int                `bson:"gold"`            // Moedas de ouro

	// Atributos base
	Attributes Attributes `bson:"attributes"`

	// Características medievais
	Class      string   `bson:"class"`      // Classe (Cavaleiro, Mago, etc)
	Background string   `bson:"background"` // Origem (Nobre, Plebeu, etc)
	Skills     []string `bson:"skills"`     // Habilidades especiais
	Equipment  []Item   `bson:"equipment"`  // Itens equipados
	Inventory  []Item   `bson:"inventory"`  // Inventário completo
	Status     []string `bson:"status"`     // Estados atuais (Ferido, Envenenado, etc)

	// Estatísticas de combate
	Combat Combat `bson:"combat"`

	// Dados de auditoria
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsActive  bool      `bson:"is_active"`
}

// Item representa um item no inventário ou equipado
type Item struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `bson:"name"`        // Nome do item
	Description string `bson:"description"` // Descrição
	Quantity    int    `bson:"quantity"`    // Quantidade
	Type        string `bson:"type"`        // Tipo (Arma, Armadura, Consumível, etc)
	Value       int    `bson:"value"`       // Valor em moedas de ouro
	Effects     string `bson:"effects"`     // Efeitos especiais
	Slot        string `bson:"slot"`        // Slot de equipamento (Mão Principal, Armadura, etc)
	IsEquipped  bool   `bson:"is_equipped"` // Se está equipado
	Rarity      string `bson:"rarity"`      // Raridade do item
}

// Attributes representa os atributos base do personagem
type Attributes struct {
	Strength     int `bson:"strength"`     // Força
	Dexterity    int `bson:"dexterity"`    // Destreza
	Constitution int `bson:"constitution"` // Constituição
	Intelligence int `bson:"intelligence"` // Inteligência
	Wisdom       int `bson:"wisdom"`       // Sabedoria
	Charisma     int `bson:"charisma"`     // Carisma
}

// Combat representa as estatísticas de combate do personagem
type Combat struct {
	Health     int  `bson:"health"`       // Pontos de vida atual
	MaxHealth  int  `bson:"max_health"`   // Pontos de vida máximo
	Armor      int  `bson:"armor"`        // Classe de armadura
	Initiative int  `bson:"initiative"`   // Iniciativa em combate
	IsInCombat bool `bson:"is_in_combat"` // Se está em combate
}

// NewCharacter cria um novo personagem com valores padrão
func NewCharacter(userID, guildID, name string) *Character {
	now := time.Now()
	char := &Character{
		UserID:     userID,
		GuildID:    guildID,
		Name:       name,
		Level:      1,
		Experience: 0,
		Gold:       100, // Moedas iniciais
		CreatedAt:  now,
		UpdatedAt:  now,
		IsActive:   true,
		Skills:     make([]string, 0),
		Equipment:  make([]Item, 0),
		Inventory:  make([]Item, 0),
		Status:     make([]string, 0),
	}

	// Atributos base iniciais
	char.Attributes.Strength = 10
	char.Attributes.Dexterity = 10
	char.Attributes.Constitution = 10
	char.Attributes.Intelligence = 10
	char.Attributes.Wisdom = 10
	char.Attributes.Charisma = 10

	// Configuração inicial de combate
	char.Combat.MaxHealth = 10 + char.getConstitutionModifier()
	char.Combat.Health = char.Combat.MaxHealth
	char.Combat.Armor = 10 + char.getDexterityModifier()
	char.Combat.Initiative = char.getDexterityModifier()

	return char
}

// Métodos auxiliares para cálculos de atributos

func (c *Character) getAttributeModifier(value int) int {
	return (value - 10) / 2
}

func (c *Character) getStrengthModifier() int {
	return c.getAttributeModifier(c.Attributes.Strength)
}

func (c *Character) getDexterityModifier() int {
	return c.getAttributeModifier(c.Attributes.Dexterity)
}

func (c *Character) getConstitutionModifier() int {
	return c.getAttributeModifier(c.Attributes.Constitution)
}

func (c *Character) getIntelligenceModifier() int {
	return c.getAttributeModifier(c.Attributes.Intelligence)
}

func (c *Character) getWisdomModifier() int {
	return c.getAttributeModifier(c.Attributes.Wisdom)
}

func (c *Character) getCharismaModifier() int {
	return c.getAttributeModifier(c.Attributes.Charisma)
}

// AddExperience adiciona experiência e verifica level up
func (c *Character) AddExperience(exp int) bool {
	c.Experience += exp
	leveledUp := false

	// Verifica level up (100 exp por nível)
	for c.Experience >= c.getNextLevelExperience() {
		c.levelUp()
		leveledUp = true
	}

	return leveledUp
}

// getNextLevelExperience retorna a experiência necessária para o próximo nível
func (c *Character) getNextLevelExperience() int {
	return c.Level * 100
}

// levelUp aumenta o nível do personagem e atualiza atributos
func (c *Character) levelUp() {
	c.Level++

	// Aumenta pontos de vida máximos
	healthIncrease := 5 + c.getConstitutionModifier()
	if healthIncrease < 1 {
		healthIncrease = 1
	}
	c.Combat.MaxHealth += healthIncrease
	c.Combat.Health = c.Combat.MaxHealth

	// Atualiza outras estatísticas baseadas em nível
	c.Combat.Initiative = c.getDexterityModifier() + (c.Level / 4)
	c.Combat.Armor = 10 + c.getDexterityModifier() + (c.Level / 5)
}

// AddItem adiciona um item ao inventário
func (c *Character) AddItem(item Item) error {
	// Procura se o item já existe no inventário
	for i, existingItem := range c.Inventory {
		if existingItem.Name == item.Name && !existingItem.IsEquipped {
			c.Inventory[i].Quantity += item.Quantity
			return nil
		}
	}
	// Se não existe, adiciona novo item
	c.Inventory = append(c.Inventory, item)
	return nil
}

// RemoveItem remove uma quantidade de um item do inventário
func (c *Character) RemoveItem(itemName string, quantity int) error {
	for i, item := range c.Inventory {
		if item.Name == itemName && !item.IsEquipped {
			if item.Quantity < quantity {
				return fmt.Errorf("quantidade insuficiente do item")
			}
			c.Inventory[i].Quantity -= quantity
			if c.Inventory[i].Quantity == 0 {
				// Remove o item se a quantidade chegar a 0
				c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			}
			return nil
		}
	}
	return fmt.Errorf("item não encontrado ou está equipado")
}

// EquipItem equipa um item do inventário
func (c *Character) EquipItem(itemName string) error {
	// Verifica se o item existe no inventário
	var itemToEquip *Item
	var itemIndex int
	for i, item := range c.Inventory {
		if item.Name == itemName && !item.IsEquipped {
			itemToEquip = &item
			itemIndex = i
			break
		}
	}

	if itemToEquip == nil {
		return fmt.Errorf("item não encontrado no inventário ou já está equipado")
	}

	// Verifica se já existe um item equipado no mesmo slot
	for i, equippedItem := range c.Equipment {
		if equippedItem.Slot == itemToEquip.Slot {
			// Desequipa o item atual
			equippedItem.IsEquipped = false
			c.Inventory = append(c.Inventory, equippedItem)
			c.Equipment = append(c.Equipment[:i], c.Equipment[i+1:]...)
			break
		}
	}

	// Equipa o novo item
	itemToEquip.IsEquipped = true
	c.Equipment = append(c.Equipment, *itemToEquip)
	// Remove do inventário
	c.Inventory = append(c.Inventory[:itemIndex], c.Inventory[itemIndex+1:]...)

	return nil
}

// UnequipItem desequipa um item
func (c *Character) UnequipItem(itemName string) error {
	for i, item := range c.Equipment {
		if item.Name == itemName {
			// Desequipa o item
			item.IsEquipped = false
			c.Inventory = append(c.Inventory, item)
			c.Equipment = append(c.Equipment[:i], c.Equipment[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item não encontrado nos equipamentos")
}

// AddStatus adiciona um status ao personagem
func (c *Character) AddStatus(status string) {
	// Verifica se o status já existe
	for _, s := range c.Status {
		if s == status {
			return
		}
	}
	c.Status = append(c.Status, status)
}

// RemoveStatus remove um status do personagem
func (c *Character) RemoveStatus(status string) bool {
	for i, s := range c.Status {
		if s == status {
			c.Status = append(c.Status[:i], c.Status[i+1:]...)
			return true
		}
	}
	return false
}

// TakeDamage aplica dano ao personagem
func (c *Character) TakeDamage(damage int) bool {
	c.Combat.Health -= damage
	if c.Combat.Health < 0 {
		c.Combat.Health = 0
	}
	return c.Combat.Health == 0
}

// Heal cura o personagem
func (c *Character) Heal(amount int) {
	c.Combat.Health += amount
	if c.Combat.Health > c.Combat.MaxHealth {
		c.Combat.Health = c.Combat.MaxHealth
	}
}
