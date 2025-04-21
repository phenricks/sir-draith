package gamedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateExpForLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		wantExp int
	}{
		{
			name:    "level 1 requires 100 exp",
			level:   1,
			wantExp: 100,
		},
		{
			name:    "level 2 requires 200 exp",
			level:   2,
			wantExp: 200,
		},
		{
			name:    "level 10 requires 1000 exp",
			level:   10,
			wantExp: 1000,
		},
		{
			name:    "max level requires max exp",
			level:   MaxLevel,
			wantExp: MaxLevel * 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateExpForLevel(tt.level)
			assert.Equal(t, tt.wantExp, got)
		})
	}
}

func TestCalculateMaxHealth(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		constitution int
		wantHealth   int
	}{
		{
			name:         "level 1 with 10 constitution",
			level:        1,
			constitution: 10,
			wantHealth:   35, // (10 base + 5 per level + 20 from constitution)
		},
		{
			name:         "level 5 with 14 constitution",
			level:        5,
			constitution: 14,
			wantHealth:   55, // (10 base + 25 from level + 28 from constitution)
		},
		{
			name:         "level 10 with 18 constitution",
			level:        10,
			constitution: 18,
			wantHealth:   110, // (10 base + 50 from level + 36 from constitution)
		},
		{
			name:         "minimum constitution",
			level:        1,
			constitution: MinAttributeValue,
			wantHealth:   31, // (10 base + 5 from level + 16 from constitution)
		},
		{
			name:         "maximum constitution",
			level:        1,
			constitution: MaxAttributeValue,
			wantHealth:   55, // (10 base + 5 from level + 40 from constitution)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMaxHealth(tt.level, tt.constitution)
			assert.Equal(t, tt.wantHealth, got)
		})
	}
}

func TestValidateEquipment(t *testing.T) {
	tests := []struct {
		name            string
		characterLevel  int
		characterClass  CharacterClass
		itemType        ItemType
		requiredLevel   int
		requiredClasses []CharacterClass
		wantErr         bool
		errContains     string
	}{
		{
			name:            "valid weapon for warrior",
			characterLevel:  5,
			characterClass:  Warrior,
			itemType:        Weapon,
			requiredLevel:   3,
			requiredClasses: []CharacterClass{Warrior, Paladin},
			wantErr:         false,
		},
		{
			name:           "invalid item type",
			characterLevel: 5,
			characterClass: Warrior,
			itemType:       Quest,
			requiredLevel:  1,
			wantErr:        true,
			errContains:    "não pode ser equipado",
		},
		{
			name:           "insufficient level",
			characterLevel: 3,
			characterClass: Warrior,
			itemType:       Weapon,
			requiredLevel:  5,
			wantErr:        true,
			errContains:    "nível insuficiente",
		},
		{
			name:            "wrong class",
			characterLevel:  5,
			characterClass:  Warrior,
			itemType:        Weapon,
			requiredLevel:   1,
			requiredClasses: []CharacterClass{Mage, Warlock},
			wantErr:         true,
			errContains:     "não pode equipar este item",
		},
		{
			name:            "no class restriction",
			characterLevel:  5,
			characterClass:  Warrior,
			itemType:        Armor,
			requiredLevel:   1,
			requiredClasses: []CharacterClass{},
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEquipment(tt.characterLevel, tt.characterClass, tt.itemType, tt.requiredLevel, tt.requiredClasses)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
