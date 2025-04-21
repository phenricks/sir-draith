package entities

import (
	"testing"

	"sirdraith/internal/domain/gamedata"

	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		guildID    string
		charName   string
		class      string
		wantErr    bool
		assertFunc func(*testing.T, *Character)
	}{
		{
			name:     "should create valid character",
			userID:   "123456789",
			guildID:  "987654321",
			charName: "Sir Test",
			class:    string(gamedata.Warrior),
			wantErr:  false,
			assertFunc: func(t *testing.T, c *Character) {
				assert.Equal(t, "123456789", c.UserID)
				assert.Equal(t, "987654321", c.GuildID)
				assert.Equal(t, "Sir Test", c.Name)
				assert.Equal(t, gamedata.Warrior, c.Class)
				assert.Equal(t, 1, c.Level)
				assert.Equal(t, 0, c.Experience)
				assert.Equal(t, gamedata.StartingGold, c.Gold)
				assert.Empty(t, c.Inventory)
				assert.True(t, c.IsActive)
				assert.False(t, c.CreatedAt.IsZero())
				assert.False(t, c.UpdatedAt.IsZero())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := NewCharacter(tt.userID, tt.guildID, tt.charName, tt.class)
			tt.assertFunc(t, char)
		})
	}
}

func TestCharacter_AddExperience(t *testing.T) {
	tests := []struct {
		name       string
		character  *Character
		expToAdd   int
		wantLevel  int
		wantExp    int
		wantErr    bool
		errMessage string
	}{
		{
			name: "should add experience without level up",
			character: &Character{
				Level:      1,
				Experience: 0,
			},
			expToAdd:  50,
			wantLevel: 1,
			wantExp:   50,
			wantErr:   false,
		},
		{
			name: "should level up once",
			character: &Character{
				Level:      1,
				Experience: 150,
			},
			expToAdd:  50,
			wantLevel: 2,
			wantExp:   200,
			wantErr:   false,
		},
		{
			name: "should not level up at max level",
			character: &Character{
				Level:      gamedata.MaxLevel,
				Experience: 50000,
			},
			expToAdd:  1000,
			wantLevel: gamedata.MaxLevel,
			wantExp:   51000,
			wantErr:   false,
		},
		{
			name: "should return error for negative experience",
			character: &Character{
				Level:      1,
				Experience: 0,
			},
			expToAdd:   -100,
			wantLevel:  1,
			wantExp:    0,
			wantErr:    true,
			errMessage: "experiência inválida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.character.AddExperience(tt.expToAdd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLevel, tt.character.Level)
				assert.Equal(t, tt.wantExp, tt.character.Experience)
			}
		})
	}
}

func TestCharacter_AddItem(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		item      Item
		wantErr   bool
		setup     func(*Character)
		assert    func(*testing.T, *Character, error)
	}{
		{
			name: "should add new item to empty inventory",
			character: &Character{
				Inventory: make([]Item, 0),
			},
			item: Item{
				Name:     "Test Sword",
				Type:     gamedata.Weapon,
				Quantity: 1,
			},
			wantErr: false,
			assert: func(t *testing.T, c *Character, err error) {
				assert.NoError(t, err)
				assert.Len(t, c.Inventory, 1)
				assert.Equal(t, "Test Sword", c.Inventory[0].Name)
				assert.Equal(t, 1, c.Inventory[0].Quantity)
			},
		},
		{
			name: "should stack same items",
			character: &Character{
				Inventory: []Item{{
					Name:     "Test Potion",
					Type:     gamedata.Consumable,
					Quantity: 1,
				}},
			},
			item: Item{
				Name:     "Test Potion",
				Type:     gamedata.Consumable,
				Quantity: 1,
			},
			wantErr: false,
			assert: func(t *testing.T, c *Character, err error) {
				assert.NoError(t, err)
				assert.Len(t, c.Inventory, 1)
				assert.Equal(t, 2, c.Inventory[0].Quantity)
			},
		},
		{
			name: "should fail when inventory is full",
			character: &Character{
				Inventory: make([]Item, gamedata.MaxInventorySize),
			},
			item: Item{
				Name:     "Test Item",
				Quantity: 1,
			},
			wantErr: true,
			assert: func(t *testing.T, c *Character, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "inventário cheio")
			},
		},
		{
			name: "should fail when stack would exceed max quantity",
			character: &Character{
				Inventory: []Item{{
					Name:     "Test Item",
					Quantity: gamedata.MaxItemQuantity - 1,
				}},
			},
			item: Item{
				Name:     "Test Item",
				Quantity: 2,
			},
			wantErr: true,
			assert: func(t *testing.T, c *Character, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "quantidade total excederia o limite")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.character)
			}

			err := tt.character.AddItem(tt.item)
			tt.assert(t, tt.character, err)
		})
	}
}

func TestCharacter_EquipItem(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		itemName  string
		wantErr   bool
		setup     func(*Character)
		assert    func(*testing.T, *Character, error)
	}{
		{
			name: "should equip valid item",
			setup: func(c *Character) {
				c.Level = 1
				c.Class = gamedata.Warrior
				c.Inventory = []Item{{
					Name:     "Test Sword",
					Type:     gamedata.Weapon,
					Quantity: 1,
				}}
				c.Equipment = make([]Item, 0)
			},
			itemName: "Test Sword",
			wantErr:  false,
			assert: func(t *testing.T, c *Character, err error) {
				assert.NoError(t, err)
				assert.Len(t, c.Equipment, 1)
				assert.Len(t, c.Inventory, 0)
				assert.True(t, c.Equipment[0].IsEquipped)
			},
		},
		{
			name: "should fail when item not in inventory",
			setup: func(c *Character) {
				c.Inventory = make([]Item, 0)
				c.Equipment = make([]Item, 0)
			},
			itemName: "Nonexistent Item",
			wantErr:  true,
			assert: func(t *testing.T, c *Character, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "item não encontrado")
			},
		},
		{
			name: "should fail when equipment slots are full",
			setup: func(c *Character) {
				c.Equipment = make([]Item, gamedata.MaxEquipmentSize)
				c.Inventory = []Item{{
					Name:     "Test Item",
					Type:     gamedata.Weapon,
					Quantity: 1,
				}}
			},
			itemName: "Test Item",
			wantErr:  true,
			assert: func(t *testing.T, c *Character, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "limite de equipamentos atingido")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &Character{}
			if tt.setup != nil {
				tt.setup(char)
			}

			err := char.EquipItem(tt.itemName)
			tt.assert(t, char, err)
		})
	}
}

func TestCharacter_UnequipItem(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		itemName  string
		wantErr   bool
		setup     func(*Character)
		assert    func(*testing.T, *Character, error)
	}{
		{
			name: "should unequip item successfully",
			setup: func(c *Character) {
				c.Equipment = []Item{{
					Name:       "Test Sword",
					Type:       gamedata.Weapon,
					IsEquipped: true,
				}}
				c.Inventory = make([]Item, 0)
			},
			itemName: "Test Sword",
			wantErr:  false,
			assert: func(t *testing.T, c *Character, err error) {
				assert.NoError(t, err)
				assert.Len(t, c.Equipment, 0)
				assert.Len(t, c.Inventory, 1)
				assert.False(t, c.Inventory[0].IsEquipped)
			},
		},
		{
			name: "should fail when item not equipped",
			setup: func(c *Character) {
				c.Equipment = make([]Item, 0)
			},
			itemName: "Nonexistent Item",
			wantErr:  true,
			assert: func(t *testing.T, c *Character, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "item não encontrado")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &Character{}
			if tt.setup != nil {
				tt.setup(char)
			}

			err := char.UnequipItem(tt.itemName)
			tt.assert(t, char, err)
		})
	}
}

func TestCharacter_Combat(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		setup     func(*Character)
		assert    func(*testing.T, *Character)
	}{
		{
			name: "should take damage correctly",
			setup: func(c *Character) {
				c.Combat.Health = 100
				c.Combat.MaxHealth = 100
			},
			assert: func(t *testing.T, c *Character) {
				isDead := c.TakeDamage(30)
				assert.Equal(t, 70, c.Combat.Health)
				assert.False(t, isDead)

				isDead = c.TakeDamage(80)
				assert.Equal(t, 0, c.Combat.Health)
				assert.True(t, isDead)
			},
		},
		{
			name: "should heal correctly",
			setup: func(c *Character) {
				c.Combat.Health = 50
				c.Combat.MaxHealth = 100
			},
			assert: func(t *testing.T, c *Character) {
				c.Heal(30)
				assert.Equal(t, 80, c.Combat.Health)

				c.Heal(30)
				assert.Equal(t, 100, c.Combat.Health)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &Character{}
			if tt.setup != nil {
				tt.setup(char)
			}

			tt.assert(t, char)
		})
	}
}

func TestCharacter_Status(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		setup     func(*Character)
		assert    func(*testing.T, *Character)
	}{
		{
			name: "should manage status correctly",
			setup: func(c *Character) {
				c.Status = make([]string, 0)
			},
			assert: func(t *testing.T, c *Character) {
				// Add status
				c.AddStatus("Envenenado")
				assert.Contains(t, c.Status, "Envenenado")
				assert.Len(t, c.Status, 1)

				// Add duplicate status
				c.AddStatus("Envenenado")
				assert.Len(t, c.Status, 1)

				// Add different status
				c.AddStatus("Atordoado")
				assert.Len(t, c.Status, 2)

				// Remove status
				removed := c.RemoveStatus("Envenenado")
				assert.True(t, removed)
				assert.NotContains(t, c.Status, "Envenenado")
				assert.Len(t, c.Status, 1)

				// Remove nonexistent status
				removed = c.RemoveStatus("Invisível")
				assert.False(t, removed)
				assert.Len(t, c.Status, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &Character{}
			if tt.setup != nil {
				tt.setup(char)
			}

			tt.assert(t, char)
		})
	}
}

func TestCharacter_AttributeModifiers(t *testing.T) {
	tests := []struct {
		name      string
		character *Character
		want      map[string]int
	}{
		{
			name: "should calculate modifiers correctly",
			character: &Character{
				Attributes: gamedata.Attributes{
					Strength:     16, // +3
					Dexterity:    14, // +2
					Constitution: 12, // +1
					Intelligence: 10, // +0
					Wisdom:       8,  // -1
					Charisma:     6,  // -2
				},
			},
			want: map[string]int{
				"strength":     3,
				"dexterity":    2,
				"constitution": 1,
				"intelligence": 0,
				"wisdom":       -1,
				"charisma":     -2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want["strength"], tt.character.getStrengthModifier())
			assert.Equal(t, tt.want["dexterity"], tt.character.getDexterityModifier())
			assert.Equal(t, tt.want["constitution"], tt.character.getConstitutionModifier())
			assert.Equal(t, tt.want["intelligence"], tt.character.getIntelligenceModifier())
			assert.Equal(t, tt.want["wisdom"], tt.character.getWisdomModifier())
			assert.Equal(t, tt.want["charisma"], tt.character.getCharismaModifier())
		})
	}
}
