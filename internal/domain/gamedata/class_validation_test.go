package gamedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateClassRequirements(t *testing.T) {
	tests := []struct {
		name        string
		class       CharacterClass
		attributes  *Attributes
		wantErr     bool
		errContains string
	}{
		{
			name:  "valid warrior attributes",
			class: Warrior,
			attributes: &Attributes{
				Strength:     13,
				Constitution: 12,
			},
			wantErr: false,
		},
		{
			name:  "invalid warrior strength",
			class: Warrior,
			attributes: &Attributes{
				Strength:     12,
				Constitution: 12,
			},
			wantErr:     true,
			errContains: "força insuficiente",
		},
		{
			name:  "valid mage attributes",
			class: Mage,
			attributes: &Attributes{
				Intelligence: 13,
				Wisdom:       12,
			},
			wantErr: false,
		},
		{
			name:  "invalid mage intelligence",
			class: Mage,
			attributes: &Attributes{
				Intelligence: 12,
				Wisdom:       12,
			},
			wantErr:     true,
			errContains: "inteligência insuficiente",
		},
		{
			name:        "invalid class",
			class:       "invalid",
			attributes:  &Attributes{},
			wantErr:     true,
			errContains: "não encontrada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClassRequirements(tt.class, tt.attributes)

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

func TestValidateClassEquipment(t *testing.T) {
	tests := []struct {
		name        string
		class       CharacterClass
		itemType    ItemType
		wantErr     bool
		errContains string
	}{
		{
			name:     "warrior can use weapon",
			class:    Warrior,
			itemType: Weapon,
			wantErr:  false,
		},
		{
			name:     "warrior can use armor",
			class:    Warrior,
			itemType: Armor,
			wantErr:  false,
		},
		{
			name:        "warrior cannot use accessory",
			class:       Warrior,
			itemType:    Accessory,
			wantErr:     true,
			errContains: "não pode usar equipamentos",
		},
		{
			name:     "mage can use weapon",
			class:    Mage,
			itemType: Weapon,
			wantErr:  false,
		},
		{
			name:        "mage cannot use armor",
			class:       Mage,
			itemType:    Armor,
			wantErr:     true,
			errContains: "não pode usar equipamentos",
		},
		{
			name:        "invalid class",
			class:       "invalid",
			itemType:    Weapon,
			wantErr:     true,
			errContains: "não encontrada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClassEquipment(tt.class, tt.itemType)

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

func TestGetClassPrimaryAttribute(t *testing.T) {
	tests := []struct {
		name  string
		class CharacterClass
		want  string
	}{
		{
			name:  "warrior primary attribute is strength",
			class: Warrior,
			want:  "strength",
		},
		{
			name:  "mage primary attribute is intelligence",
			class: Mage,
			want:  "intelligence",
		},
		{
			name:  "ranger primary attribute is dexterity",
			class: Ranger,
			want:  "dexterity",
		},
		{
			name:  "cleric primary attribute is wisdom",
			class: Cleric,
			want:  "wisdom",
		},
		{
			name:  "paladin primary attribute is charisma",
			class: Paladin,
			want:  "charisma",
		},
		{
			name:  "invalid class returns empty string",
			class: "invalid",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetClassPrimaryAttribute(tt.class)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateClassBonus(t *testing.T) {
	tests := []struct {
		name       string
		class      CharacterClass
		attributes *Attributes
		want       int
	}{
		{
			name:  "warrior with high strength",
			class: Warrior,
			attributes: &Attributes{
				Strength: 16,
			},
			want: 5, // Base bonus 3 (16-10)/2 + 2 warrior bonus
		},
		{
			name:  "mage with high intelligence",
			class: Mage,
			attributes: &Attributes{
				Intelligence: 16,
			},
			want: 4, // Base bonus 3 (16-10)/2 + 1 mage bonus
		},
		{
			name:  "ranger with high dexterity",
			class: Ranger,
			attributes: &Attributes{
				Dexterity: 16,
			},
			want: 4, // Base bonus 3 (16-10)/2 + 1 precision bonus
		},
		{
			name:  "barbarian with high constitution",
			class: Barbarian,
			attributes: &Attributes{
				Strength:     16,
				Constitution: 16,
			},
			want: 4, // Base bonus 3 (16-10)/2 + 1 constitution bonus
		},
		{
			name:  "invalid class returns 0",
			class: "invalid",
			attributes: &Attributes{
				Strength: 16,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateClassBonus(tt.class, tt.attributes)
			assert.Equal(t, tt.want, got)
		})
	}
}
