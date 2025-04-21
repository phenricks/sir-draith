package gamedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSkillModifier(t *testing.T) {
	tests := []struct {
		name        string
		skill       Skill
		attributes  *Attributes
		proficiency *SkillProficiency
		level       int
		want        int
	}{
		{
			name:  "athletics with strength 16, proficient",
			skill: Athletics,
			attributes: &Attributes{
				Strength: 16,
			},
			proficiency: &SkillProficiency{
				Skill:        Athletics,
				IsProficient: true,
			},
			level: 1,
			want:  5, // (16-10)/2 + 2 (proficiency at level 1)
		},
		{
			name:  "athletics with strength 16, not proficient",
			skill: Athletics,
			attributes: &Attributes{
				Strength: 16,
			},
			proficiency: &SkillProficiency{
				Skill:        Athletics,
				IsProficient: false,
			},
			level: 1,
			want:  3, // (16-10)/2
		},
		{
			name:  "stealth with dexterity 18, proficient, level 5",
			skill: Stealth,
			attributes: &Attributes{
				Dexterity: 18,
			},
			proficiency: &SkillProficiency{
				Skill:        Stealth,
				IsProficient: true,
				Bonus:        1, // Bônus adicional
			},
			level: 5,
			want:  7, // (18-10)/2 + 3 (proficiency at level 5) + 1 (bonus)
		},
		{
			name:  "arcana with intelligence 14, proficient",
			skill: Arcana,
			attributes: &Attributes{
				Intelligence: 14,
			},
			proficiency: &SkillProficiency{
				Skill:        Arcana,
				IsProficient: true,
			},
			level: 1,
			want:  4, // (14-10)/2 + 2 (proficiency at level 1)
		},
		{
			name:  "invalid skill returns 0",
			skill: "invalid",
			attributes: &Attributes{
				Strength: 16,
			},
			proficiency: nil,
			level:       1,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSkillModifier(tt.skill, tt.attributes, tt.proficiency, tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetSkillsForClass(t *testing.T) {
	tests := []struct {
		name      string
		class     CharacterClass
		wantCount int
		wantSkill Skill // Uma perícia que deve estar presente
	}{
		{
			name:      "warrior skills",
			class:     Warrior,
			wantCount: 4,
			wantSkill: Athletics,
		},
		{
			name:      "mage skills",
			class:     Mage,
			wantCount: 4,
			wantSkill: Arcana,
		},
		{
			name:      "ranger skills",
			class:     Ranger,
			wantCount: 5,
			wantSkill: Survival,
		},
		{
			name:      "invalid class returns empty slice",
			class:     "invalid",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSkillsForClass(tt.class)
			assert.Len(t, got, tt.wantCount)

			if tt.wantCount > 0 {
				found := false
				for _, skill := range got {
					if skill == tt.wantSkill {
						found = true
						break
					}
				}
				assert.True(t, found, "expected skill not found")
			}
		})
	}
}

func TestValidateSkillProficiency(t *testing.T) {
	tests := []struct {
		name  string
		class CharacterClass
		skill Skill
		want  bool
	}{
		{
			name:  "warrior can use athletics",
			class: Warrior,
			skill: Athletics,
			want:  true,
		},
		{
			name:  "warrior cannot use arcana",
			class: Warrior,
			skill: Arcana,
			want:  false,
		},
		{
			name:  "mage can use arcana",
			class: Mage,
			skill: Arcana,
			want:  true,
		},
		{
			name:  "ranger can use survival",
			class: Ranger,
			skill: Survival,
			want:  true,
		},
		{
			name:  "invalid class returns false",
			class: "invalid",
			skill: Athletics,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSkillProficiency(tt.class, tt.skill)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSkillBaseAttribute(t *testing.T) {
	tests := []struct {
		name     string
		skill    Skill
		wantAttr string
	}{
		{
			name:     "athletics uses strength",
			skill:    Athletics,
			wantAttr: "strength",
		},
		{
			name:     "stealth uses dexterity",
			skill:    Stealth,
			wantAttr: "dexterity",
		},
		{
			name:     "arcana uses intelligence",
			skill:    Arcana,
			wantAttr: "intelligence",
		},
		{
			name:     "perception uses wisdom",
			skill:    Perception,
			wantAttr: "wisdom",
		},
		{
			name:     "persuasion uses charisma",
			skill:    Persuasion,
			wantAttr: "charisma",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SkillBaseAttribute[tt.skill]
			assert.Equal(t, tt.wantAttr, got)
		})
	}
}
