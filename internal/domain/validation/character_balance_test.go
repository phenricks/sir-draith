package validation

import (
	"testing"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"

	"github.com/stretchr/testify/assert"
)

func TestCalculateExpForLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{
			name:     "nível 1 deve retornar 0",
			level:    1,
			expected: 0,
		},
		{
			name:     "nível 2 deve retornar 100",
			level:    2,
			expected: 100,
		},
		{
			name:     "nível 3 deve retornar 150",
			level:    3,
			expected: 150,
		},
		{
			name:     "nível máximo + 1 deve retornar -1",
			level:    gamedata.MaxLevel + 1,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateExpForLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMaxHealth(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		constitution int
		expected     int
	}{
		{
			name:         "nível 1 com constituição 10 deve ter 10 HP",
			level:        1,
			constitution: 10,
			expected:     10,
		},
		{
			name:         "nível 1 com constituição 14 deve ter 12 HP",
			level:        1,
			constitution: 14,
			expected:     12,
		},
		{
			name:         "nível 2 com constituição 12 deve ter 22 HP",
			level:        2,
			constitution: 12,
			expected:     22,
		},
		{
			name:         "nível inválido deve retornar 0",
			level:        0,
			constitution: 10,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMaxHealth(tt.level, tt.constitution)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateAttackBonus(t *testing.T) {
	tests := []struct {
		name       string
		class      gamedata.CharacterClass
		level      int
		attributes *gamedata.Attributes
		expected   int
	}{
		{
			name:  "guerreiro nível 1 com força 14",
			class: gamedata.Warrior,
			level: 1,
			attributes: &gamedata.Attributes{
				Strength: 14,
			},
			expected: 4, // +2 proficiência, +2 força
		},
		{
			name:  "mago nível 1 com inteligência 16",
			class: gamedata.Mage,
			level: 1,
			attributes: &gamedata.Attributes{
				Intelligence: 16,
			},
			expected: 5, // +2 proficiência, +3 inteligência
		},
		{
			name:  "arqueiro nível 4 com destreza 18",
			class: gamedata.Ranger,
			level: 4,
			attributes: &gamedata.Attributes{
				Dexterity: 18,
			},
			expected: 6, // +2 proficiência, +4 destreza
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAttackBonus(tt.class, tt.level, tt.attributes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateDefense(t *testing.T) {
	tests := []struct {
		name       string
		class      gamedata.CharacterClass
		attributes *gamedata.Attributes
		expected   int
	}{
		{
			name:  "guerreiro com força 16 e destreza 14",
			class: gamedata.Warrior,
			attributes: &gamedata.Attributes{
				Strength:  16,
				Dexterity: 14,
			},
			expected: 13, // 10 + 3 (força)
		},
		{
			name:  "mago com inteligência 18",
			class: gamedata.Mage,
			attributes: &gamedata.Attributes{
				Intelligence: 18,
			},
			expected: 14, // 10 + 4 (inteligência)
		},
		{
			name:  "arqueiro com destreza 20",
			class: gamedata.Ranger,
			attributes: &gamedata.Attributes{
				Dexterity: 20,
			},
			expected: 15, // 10 + 5 (destreza)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDefense(tt.class, tt.attributes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateCharacterProgression(t *testing.T) {
	tests := []struct {
		name        string
		character   *entities.Character
		expectError bool
	}{
		{
			name: "personagem válido",
			character: &entities.Character{
				Level:      1,
				Experience: 0,
				Attributes: gamedata.Attributes{
					Constitution: 14,
				},
				Combat: entities.Combat{
					MaxHealth: 12,
					Health:    12,
				},
			},
			expectError: false,
		},
		{
			name: "nível inválido",
			character: &entities.Character{
				Level:      0,
				Experience: 0,
			},
			expectError: true,
		},
		{
			name: "experiência insuficiente",
			character: &entities.Character{
				Level:      2,
				Experience: 50,
			},
			expectError: true,
		},
		{
			name: "HP atual maior que máximo",
			character: &entities.Character{
				Level:      1,
				Experience: 0,
				Attributes: gamedata.Attributes{
					Constitution: 10,
				},
				Combat: entities.Combat{
					MaxHealth: 10,
					Health:    11,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCharacterProgression(tt.character)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
