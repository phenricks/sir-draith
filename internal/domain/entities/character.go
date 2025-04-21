package entities

import (
	"fmt"
	"time"

	"sirdraith/internal/domain/gamedata"

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
	Attributes gamedata.Attributes `bson:"attributes"`

	// Características medievais
	Class      gamedata.CharacterClass     `bson:"class"`      // Classe (Cavaleiro, Mago, etc)
	Background string                      `bson:"background"` // Origem (Nobre, Plebeu, etc)
	Skills     []gamedata.SkillProficiency `bson:"skills"`     // Perícias do personagem
	Equipment  []Item                      `bson:"equipment"`  // Itens equipados
	Inventory  []Item                      `bson:"inventory"`  // Inventário completo
	Status     []string                    `bson:"status"`     // Estados atuais (Ferido, Envenenado, etc)

	// Estatísticas de combate
	Combat Combat `bson:"combat"`

	// Dados de auditoria
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsActive  bool      `bson:"is_active"`
}

// Item representa um item no inventário ou equipado
type Item struct {
	ID              string                    `bson:"_id,omitempty"`
	Name            string                    `bson:"name"`             // Nome do item
	Description     string                    `bson:"description"`      // Descrição
	Quantity        int                       `bson:"quantity"`         // Quantidade
	Type            gamedata.ItemType         `bson:"type"`             // Tipo (Arma, Armadura, Consumível, etc)
	Value           int                       `bson:"value"`            // Valor em moedas de ouro
	Effects         string                    `bson:"effects"`          // Efeitos especiais
	Slot            gamedata.EquipmentSlot    `bson:"slot"`             // Slot de equipamento
	IsEquipped      bool                      `bson:"is_equipped"`      // Se está equipado
	Rarity          gamedata.ItemRarity       `bson:"rarity"`           // Raridade do item
	Stats           gamedata.ItemStats        `bson:"stats"`            // Estatísticas do item
	RequiredLevel   int                       `bson:"required_level"`   // Nível mínimo para usar
	RequiredClasses []gamedata.CharacterClass `bson:"required_classes"` // Classes que podem usar
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
func NewCharacter(userID, guildID, name, class string) *Character {
	now := time.Now()
	char := &Character{
		UserID:     userID,
		GuildID:    guildID,
		Name:       name,
		Class:      gamedata.CharacterClass(class),
		Level:      1,
		Experience: 0,
		Gold:       gamedata.StartingGold,
		Inventory:  make([]Item, 0),
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	char.updateCombatConfig()
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
func (c *Character) AddExperience(exp int) error {
	if exp < 0 {
		return fmt.Errorf("experiência inválida")
	}

	c.Experience += exp
	for {
		nextLevelExp := c.getNextLevelExperience()
		if nextLevelExp == -1 || c.Experience < nextLevelExp {
			break
		}
		if err := c.levelUp(); err != nil {
			return err
		}
	}
	return nil
}

// getNextLevelExperience retorna a experiência necessária para o próximo nível
func (c *Character) getNextLevelExperience() int {
	if c.Level >= gamedata.MaxLevel {
		return -1 // Indica que não há próximo nível
	}
	return gamedata.CalculateExpForLevel(c.Level + 1)
}

// levelUp aumenta o nível do personagem e atualiza atributos
func (c *Character) levelUp() error {
	if c.Level >= gamedata.MaxLevel {
		return fmt.Errorf("nível máximo atingido")
	}
	c.Level++
	c.updateCombatConfig()
	return nil
}

func (c *Character) updateCombatConfig() {
	c.Combat.MaxHealth = gamedata.CalculateMaxHealth(c.Level, c.Attributes.Constitution)
	c.Combat.Health = c.Combat.MaxHealth
	c.Combat.Armor = 10 + c.getDexterityModifier()
	c.Combat.Initiative = c.getDexterityModifier() + (c.Level / 4)
}

// AddItem adiciona um item ao inventário
func (c *Character) AddItem(item Item) error {
	if len(c.Inventory) >= gamedata.MaxInventorySize {
		return fmt.Errorf("inventário cheio (limite de %d itens)", gamedata.MaxInventorySize)
	}

	for i, existingItem := range c.Inventory {
		if existingItem.Name == item.Name && !existingItem.IsEquipped {
			newQuantity := existingItem.Quantity + item.Quantity
			if newQuantity > gamedata.MaxItemQuantity {
				return fmt.Errorf("quantidade total excederia o limite de %d", gamedata.MaxItemQuantity)
			}
			c.Inventory[i].Quantity = newQuantity
			return nil
		}
	}

	c.Inventory = append(c.Inventory, item)
	return nil
}

// RemoveItem remove uma quantidade de um item do inventário
func (c *Character) RemoveItem(itemName string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantidade deve ser maior que zero")
	}

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
	if len(c.Equipment) >= gamedata.MaxEquipmentSize {
		return fmt.Errorf("limite de equipamentos atingido (%d itens)", gamedata.MaxEquipmentSize)
	}

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

	if err := gamedata.ValidateEquipment(c.Level, c.Class, itemToEquip.Type, itemToEquip.RequiredLevel, itemToEquip.RequiredClasses); err != nil {
		return fmt.Errorf("não é possível equipar o item: %w", err)
	}

	itemToEquip.IsEquipped = true
	c.Equipment = append(c.Equipment, *itemToEquip)
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

// CalculateMaxHealth calcula a vida máxima do personagem
func (c *Character) CalculateMaxHealth() int {
	return gamedata.CalculateMaxHealth(c.Level, c.Attributes.Constitution)
}

// InitializeSkills inicializa as perícias do personagem com base na classe
func (c *Character) InitializeSkills() {
	// Obter perícias disponíveis para a classe
	classSkills := gamedata.GetSkillsForClass(c.Class)

	// Inicializar slice de perícias
	c.Skills = make([]gamedata.SkillProficiency, 0, len(classSkills))

	// Adicionar cada perícia com proficiência
	for _, skill := range classSkills {
		c.Skills = append(c.Skills, gamedata.SkillProficiency{
			Skill:        skill,
			IsProficient: true,
		})
	}
}

// GetSkillModifier retorna o modificador total de uma perícia
func (c *Character) GetSkillModifier(skill gamedata.Skill) int {
	// Encontrar a proficiência da perícia
	var proficiency *gamedata.SkillProficiency
	for i, p := range c.Skills {
		if p.Skill == skill {
			proficiency = &c.Skills[i]
			break
		}
	}

	return gamedata.CalculateSkillModifier(skill, &c.Attributes, proficiency, c.Level)
}

// AddSkillProficiency adiciona proficiência em uma perícia
func (c *Character) AddSkillProficiency(skill gamedata.Skill) error {
	// Verificar se a perícia pode ser usada pela classe
	if !gamedata.ValidateSkillProficiency(c.Class, skill) {
		return fmt.Errorf("a classe %s não pode ter proficiência em %s", c.Class, skill)
	}

	// Verificar se já tem a perícia
	for i, p := range c.Skills {
		if p.Skill == skill {
			c.Skills[i].IsProficient = true
			return nil
		}
	}

	// Adicionar nova perícia
	c.Skills = append(c.Skills, gamedata.SkillProficiency{
		Skill:        skill,
		IsProficient: true,
	})
	return nil
}

// RemoveSkillProficiency remove a proficiência em uma perícia
func (c *Character) RemoveSkillProficiency(skill gamedata.Skill) {
	for i, p := range c.Skills {
		if p.Skill == skill {
			c.Skills[i].IsProficient = false
			break
		}
	}
}

// AddSkillBonus adiciona um bônus a uma perícia
func (c *Character) AddSkillBonus(skill gamedata.Skill, bonus int) error {
	// Encontrar a perícia
	for i, p := range c.Skills {
		if p.Skill == skill {
			c.Skills[i].Bonus += bonus
			return nil
		}
	}

	// Se não encontrou, adiciona nova perícia com o bônus
	c.Skills = append(c.Skills, gamedata.SkillProficiency{
		Skill: skill,
		Bonus: bonus,
	})
	return nil
}
