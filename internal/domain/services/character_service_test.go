package services

import (
	"context"
	"errors"
	"testing"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repository"
)

// MockCharacterRepository é um mock do repositório de personagens para testes
type MockCharacterRepository struct {
	characters map[string]*entities.Character
}

func NewMockCharacterRepository() *MockCharacterRepository {
	return &MockCharacterRepository{
		characters: make(map[string]*entities.Character),
	}
}

func (m *MockCharacterRepository) Create(ctx context.Context, character *entities.Character) error {
	id := character.ID.Hex()
	if _, exists := m.characters[id]; exists {
		return errors.New("character already exists")
	}
	m.characters[id] = character
	return nil
}

func (m *MockCharacterRepository) Update(ctx context.Context, character *entities.Character) error {
	id := character.ID.Hex()
	if _, exists := m.characters[id]; !exists {
		return repository.ErrCharacterNotFound
	}
	m.characters[id] = character
	return nil
}

func (m *MockCharacterRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.characters[id]; !exists {
		return repository.ErrCharacterNotFound
	}
	delete(m.characters, id)
	return nil
}

func (m *MockCharacterRepository) GetByID(ctx context.Context, id string) (*entities.Character, error) {
	if character, exists := m.characters[id]; exists {
		return character, nil
	}
	return nil, repository.ErrCharacterNotFound
}

func (m *MockCharacterRepository) GetByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error) {
	for _, character := range m.characters {
		if character.UserID == userID && character.GuildID == guildID {
			return character, nil
		}
	}
	return nil, repository.ErrCharacterNotFound
}

func (m *MockCharacterRepository) GetByGuildID(ctx context.Context, guildID string) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		if character.GuildID == guildID {
			result = append(result, character)
		}
	}
	return result, nil
}

func (m *MockCharacterRepository) ListByUser(ctx context.Context, userID string) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		if character.UserID == userID {
			result = append(result, character)
		}
	}
	return result, nil
}

func (m *MockCharacterRepository) ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		if character.GuildID == guildID {
			result = append(result, character)
		}
	}
	return result, nil
}

func (m *MockCharacterRepository) ListActive(ctx context.Context) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		if character.IsActive {
			result = append(result, character)
		}
	}
	return result, nil
}

func (m *MockCharacterRepository) Search(ctx context.Context, query string) ([]*entities.Character, error) {
	// Implementação simplificada para testes
	return m.ListActive(ctx)
}

func (m *MockCharacterRepository) CountByGuild(ctx context.Context, guildID string) (int64, error) {
	var count int64
	for _, character := range m.characters {
		if character.GuildID == guildID {
			count++
		}
	}
	return count, nil
}

func (m *MockCharacterRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	for _, character := range m.characters {
		if character.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *MockCharacterRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		if character.UserID == userID {
			result = append(result, character)
		}
	}
	return result, nil
}

func (m *MockCharacterRepository) List(ctx context.Context) ([]*entities.Character, error) {
	var result []*entities.Character
	for _, character := range m.characters {
		result = append(result, character)
	}
	return result, nil
}

func TestCharacterService_CreateCharacter(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MockCharacterRepository)
		userID   string
		guildID  string
		charName string
		wantErr  bool
	}{
		{
			name:     "should create character successfully",
			setup:    func(r *MockCharacterRepository) {},
			userID:   "123",
			guildID:  "456",
			charName: "Sir Test",
			wantErr:  false,
		},
		{
			name: "should error when character already exists",
			setup: func(r *MockCharacterRepository) {
				char := entities.NewCharacter("123", "456", "Existing")
				r.Create(context.Background(), char)
			},
			userID:   "123",
			guildID:  "456",
			charName: "Sir Test",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			tt.setup(repo)
			service := NewCharacterService(repo)

			character, err := service.CreateCharacter(context.Background(), tt.userID, tt.guildID, tt.charName)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCharacter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if character == nil {
					t.Error("CreateCharacter() returned nil character")
					return
				}
				if character.UserID != tt.userID {
					t.Errorf("CreateCharacter().UserID = %v, want %v", character.UserID, tt.userID)
				}
				if character.GuildID != tt.guildID {
					t.Errorf("CreateCharacter().GuildID = %v, want %v", character.GuildID, tt.guildID)
				}
				if character.Name != tt.charName {
					t.Errorf("CreateCharacter().Name = %v, want %v", character.Name, tt.charName)
				}
			}
		})
	}
}

func TestCharacterService_GetCharacterByUserAndGuild(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MockCharacterRepository)
		userID   string
		guildID  string
		wantName string
		wantErr  bool
	}{
		{
			name: "should get character successfully",
			setup: func(r *MockCharacterRepository) {
				char := entities.NewCharacter("123", "456", "Sir Test")
				r.Create(context.Background(), char)
			},
			userID:   "123",
			guildID:  "456",
			wantName: "Sir Test",
			wantErr:  false,
		},
		{
			name:     "should error when character not found",
			setup:    func(r *MockCharacterRepository) {},
			userID:   "123",
			guildID:  "456",
			wantName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			tt.setup(repo)
			service := NewCharacterService(repo)

			character, err := service.GetCharacterByUserAndGuild(context.Background(), tt.userID, tt.guildID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCharacterByUserAndGuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if character == nil {
					t.Error("GetCharacterByUserAndGuild() returned nil character")
					return
				}
				if character.Name != tt.wantName {
					t.Errorf("GetCharacterByUserAndGuild().Name = %v, want %v", character.Name, tt.wantName)
				}
			}
		})
	}
}

func TestCharacterService_ListCharactersByGuild(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*MockCharacterRepository)
		guildID   string
		wantCount int
		wantErr   bool
	}{
		{
			name: "should list characters successfully",
			setup: func(r *MockCharacterRepository) {
				r.Create(context.Background(), entities.NewCharacter("123", "456", "Char1"))
				r.Create(context.Background(), entities.NewCharacter("124", "456", "Char2"))
				r.Create(context.Background(), entities.NewCharacter("125", "789", "Char3"))
			},
			guildID:   "456",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "should return empty list when no characters found",
			setup:     func(r *MockCharacterRepository) {},
			guildID:   "456",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			tt.setup(repo)
			service := NewCharacterService(repo)

			characters, err := service.ListCharactersByGuild(context.Background(), tt.guildID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListCharactersByGuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(characters) != tt.wantCount {
					t.Errorf("ListCharactersByGuild() count = %v, want %v", len(characters), tt.wantCount)
				}
			}
		})
	}
}

func TestCharacterService_EquipItem(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*entities.Character)
		itemName string
		wantErr  bool
	}{
		{
			name: "should equip item successfully",
			setup: func(c *entities.Character) {
				c.AddItem(entities.Item{
					Name: "Espada Longa",
					Slot: "Mão Principal",
				})
			},
			itemName: "Espada Longa",
			wantErr:  false,
		},
		{
			name:     "should error when item not found",
			setup:    func(c *entities.Character) {},
			itemName: "Item Inexistente",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character := entities.NewCharacter("123", "456", "Test")
			tt.setup(character)

			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), character)
			service := NewCharacterService(repo)

			err := service.EquipItem(context.Background(), character, tt.itemName)

			if (err != nil) != tt.wantErr {
				t.Errorf("EquipItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verifica se o item está equipado
				var found bool
				for _, item := range character.Equipment {
					if item.Name == tt.itemName {
						found = true
						if !item.IsEquipped {
							t.Error("EquipItem() item não está marcado como equipado")
						}
						break
					}
				}
				if !found {
					t.Error("EquipItem() item não encontrado nos equipamentos")
				}
			}
		})
	}
}

func TestCharacterService_UnequipItem(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*entities.Character)
		itemName string
		wantErr  bool
	}{
		{
			name: "should unequip item successfully",
			setup: func(c *entities.Character) {
				c.AddItem(entities.Item{
					Name: "Espada Longa",
					Slot: "Mão Principal",
				})
				c.EquipItem("Espada Longa")
			},
			itemName: "Espada Longa",
			wantErr:  false,
		},
		{
			name:     "should error when item not equipped",
			setup:    func(c *entities.Character) {},
			itemName: "Item Inexistente",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character := entities.NewCharacter("123", "456", "Test")
			tt.setup(character)

			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), character)
			service := NewCharacterService(repo)

			err := service.UnequipItem(context.Background(), character, tt.itemName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnequipItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verifica se o item foi removido dos equipamentos
				for _, item := range character.Equipment {
					if item.Name == tt.itemName {
						t.Error("UnequipItem() item ainda está nos equipamentos")
					}
				}

				// Verifica se o item voltou para o inventário
				var found bool
				for _, item := range character.Inventory {
					if item.Name == tt.itemName {
						found = true
						if item.IsEquipped {
							t.Error("UnequipItem() item ainda está marcado como equipado")
						}
						break
					}
				}
				if !found {
					t.Error("UnequipItem() item não foi adicionado ao inventário")
				}
			}
		})
	}
}

func TestCharacterService_GetCharacter(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockCharacterRepository)
		id      string
		wantErr bool
	}{
		{
			name: "should get character successfully",
			setup: func(r *MockCharacterRepository) {
				char := entities.NewCharacter("123", "456", "Sir Test")
				r.Create(context.Background(), char)
			},
			id:      "existing_id", // Será substituído no teste
			wantErr: false,
		},
		{
			name:    "should error when character not found",
			setup:   func(r *MockCharacterRepository) {},
			id:      "non_existing_id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			tt.setup(repo)
			service := NewCharacterService(repo)

			// Se houver um personagem criado, usa seu ID real
			if len(repo.characters) > 0 {
				for id := range repo.characters {
					tt.id = id
					break
				}
			}

			character, err := service.GetCharacter(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCharacter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && character == nil {
				t.Error("GetCharacter() returned nil character when error not expected")
			}
		})
	}
}

func TestCharacterService_UpdateCharacter(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MockCharacterRepository) *entities.Character
		updateFn func(*entities.Character)
		wantErr  bool
	}{
		{
			name: "should update character successfully",
			setup: func(r *MockCharacterRepository) *entities.Character {
				char := entities.NewCharacter("123", "456", "Sir Test")
				r.Create(context.Background(), char)
				return char
			},
			updateFn: func(c *entities.Character) {
				c.Name = "Sir Test Updated"
				c.Title = "The Brave"
			},
			wantErr: false,
		},
		{
			name: "should error when character not found",
			setup: func(r *MockCharacterRepository) *entities.Character {
				return entities.NewCharacter("123", "456", "Sir Test")
			},
			updateFn: func(c *entities.Character) {},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			character := tt.setup(repo)
			service := NewCharacterService(repo)

			tt.updateFn(character)
			err := service.UpdateCharacter(context.Background(), character)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verifica se as alterações foram salvas
				updated, _ := repo.GetByID(context.Background(), character.ID.Hex())
				if updated.Name != character.Name {
					t.Errorf("UpdateCharacter() name = %v, want %v", updated.Name, character.Name)
				}
				if updated.Title != character.Title {
					t.Errorf("UpdateCharacter() title = %v, want %v", updated.Title, character.Title)
				}
			}
		})
	}
}

func TestCharacterService_DeleteCharacter(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockCharacterRepository) string
		wantErr bool
	}{
		{
			name: "should delete character successfully",
			setup: func(r *MockCharacterRepository) string {
				char := entities.NewCharacter("123", "456", "Sir Test")
				r.Create(context.Background(), char)
				return char.ID.Hex()
			},
			wantErr: false,
		},
		{
			name: "should error when character not found",
			setup: func(r *MockCharacterRepository) string {
				return "non_existing_id"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			id := tt.setup(repo)
			service := NewCharacterService(repo)

			err := service.DeleteCharacter(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verifica se o personagem foi realmente deletado
				_, err := repo.GetByID(context.Background(), id)
				if err == nil {
					t.Error("DeleteCharacter() character still exists after deletion")
				}
			}
		})
	}
}

func TestCharacterService_AddGold(t *testing.T) {
	tests := []struct {
		name      string
		character *entities.Character
		amount    int
		wantGold  int
		wantErr   bool
	}{
		{
			name:      "should add gold successfully",
			character: entities.NewCharacter("123", "456", "Sir Test"),
			amount:    50,
			wantGold:  150, // 100 (inicial) + 50
			wantErr:   false,
		},
		{
			name:      "should remove gold successfully",
			character: entities.NewCharacter("123", "456", "Sir Test"),
			amount:    -50,
			wantGold:  50, // 100 (inicial) - 50
			wantErr:   false,
		},
		{
			name:      "should not allow negative gold",
			character: entities.NewCharacter("123", "456", "Sir Test"),
			amount:    -150,
			wantGold:  0, // Não pode ficar negativo
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), tt.character)
			service := NewCharacterService(repo)

			err := service.AddGold(context.Background(), tt.character, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddGold() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.character.Gold != tt.wantGold {
				t.Errorf("AddGold() gold = %v, want %v", tt.character.Gold, tt.wantGold)
			}
		})
	}
}

func TestCharacterService_TakeDamage(t *testing.T) {
	tests := []struct {
		name         string
		character    *entities.Character
		damage       int
		wantHealth   int
		wantDefeated bool
		wantErr      bool
	}{
		{
			name:         "should take damage successfully",
			character:    entities.NewCharacter("123", "456", "Sir Test"),
			damage:       5,
			wantHealth:   5, // 10 (inicial) - 5
			wantDefeated: false,
			wantErr:      false,
		},
		{
			name:         "should be defeated when health reaches 0",
			character:    entities.NewCharacter("123", "456", "Sir Test"),
			damage:       15,
			wantHealth:   0,
			wantDefeated: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), tt.character)
			service := NewCharacterService(repo)

			err := service.TakeDamage(context.Background(), tt.character, tt.damage)

			if (err != nil) != tt.wantErr {
				t.Errorf("TakeDamage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.character.Combat.Health != tt.wantHealth {
				t.Errorf("TakeDamage() health = %v, want %v", tt.character.Combat.Health, tt.wantHealth)
			}

			defeated := tt.character.Combat.Health == 0
			if defeated != tt.wantDefeated {
				t.Errorf("TakeDamage() defeated = %v, want %v", defeated, tt.wantDefeated)
			}
		})
	}
}

func TestCharacterService_Heal(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*entities.Character)
		healAmount int
		wantHealth int
		wantErr    bool
	}{
		{
			name: "should heal damage successfully",
			setup: func(c *entities.Character) {
				c.TakeDamage(5) // Reduz a vida para 5
			},
			healAmount: 3,
			wantHealth: 8, // 5 + 3
			wantErr:    false,
		},
		{
			name: "should not exceed max health when healing",
			setup: func(c *entities.Character) {
				c.TakeDamage(2) // Reduz a vida para 8
			},
			healAmount: 5,
			wantHealth: 10, // Não pode exceder o máximo
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character := entities.NewCharacter("123", "456", "Sir Test")
			tt.setup(character)

			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), character)
			service := NewCharacterService(repo)

			err := service.Heal(context.Background(), character, tt.healAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("Heal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if character.Combat.Health != tt.wantHealth {
				t.Errorf("Heal() health = %v, want %v", character.Combat.Health, tt.wantHealth)
			}
		})
	}
}

func TestCharacterService_SearchCharacters(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*MockCharacterRepository)
		query     string
		wantCount int
		wantErr   bool
	}{
		{
			name: "should find characters by name",
			setup: func(r *MockCharacterRepository) {
				r.Create(context.Background(), entities.NewCharacter("123", "456", "Sir Test"))
				r.Create(context.Background(), entities.NewCharacter("124", "456", "Sir Test II"))
				r.Create(context.Background(), entities.NewCharacter("125", "456", "Knight"))
			},
			query:     "Sir",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "should return empty when no matches",
			setup: func(r *MockCharacterRepository) {
				r.Create(context.Background(), entities.NewCharacter("123", "456", "Knight"))
			},
			query:     "Wizard",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			tt.setup(repo)
			service := NewCharacterService(repo)

			characters, err := service.SearchCharacters(context.Background(), tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchCharacters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(characters) != tt.wantCount {
				t.Errorf("SearchCharacters() count = %v, want %v", len(characters), tt.wantCount)
			}
		})
	}
}

func TestCharacterService_AddExperience(t *testing.T) {
	tests := []struct {
		name        string
		character   *entities.Character
		exp         int
		wantLevel   int
		wantExp     int
		wantLevelUp bool
		wantErr     bool
	}{
		{
			name:        "should add experience without level up",
			character:   entities.NewCharacter("123", "456", "Sir Test"),
			exp:         50,
			wantLevel:   1,
			wantExp:     50,
			wantLevelUp: false,
			wantErr:     false,
		},
		{
			name:        "should level up with sufficient experience",
			character:   entities.NewCharacter("123", "456", "Sir Test"),
			exp:         150,
			wantLevel:   2,
			wantExp:     50,
			wantLevelUp: true,
			wantErr:     false,
		},
		{
			name: "should level up multiple times",
			character: func() *entities.Character {
				c := entities.NewCharacter("123", "456", "Sir Test")
				c.Experience = 90
				return c
			}(),
			exp:         210,
			wantLevel:   3,
			wantExp:     0,
			wantLevelUp: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), tt.character)
			service := NewCharacterService(repo)

			err := service.AddExperience(context.Background(), tt.character, tt.exp)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddExperience() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.character.Level != tt.wantLevel {
				t.Errorf("AddExperience() level = %v, want %v", tt.character.Level, tt.wantLevel)
			}

			if tt.character.Experience != tt.wantExp {
				t.Errorf("AddExperience() experience = %v, want %v", tt.character.Experience, tt.wantExp)
			}

			// Verifica se o personagem foi atualizado no repositório
			updated, _ := repo.GetByID(context.Background(), tt.character.ID.Hex())
			if updated.Level != tt.wantLevel {
				t.Errorf("AddExperience() saved level = %v, want %v", updated.Level, tt.wantLevel)
			}
		})
	}
}

func TestCharacterService_AddItem(t *testing.T) {
	tests := []struct {
		name      string
		character *entities.Character
		item      entities.Item
		wantQty   int
		wantErr   bool
	}{
		{
			name:      "should add new item",
			character: entities.NewCharacter("123", "456", "Sir Test"),
			item: entities.Item{
				Name:     "Poção de Cura",
				Quantity: 1,
			},
			wantQty: 1,
			wantErr: false,
		},
		{
			name: "should stack existing item",
			character: func() *entities.Character {
				c := entities.NewCharacter("123", "456", "Sir Test")
				c.AddItem(entities.Item{Name: "Poção de Cura", Quantity: 1})
				return c
			}(),
			item: entities.Item{
				Name:     "Poção de Cura",
				Quantity: 2,
			},
			wantQty: 3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), tt.character)
			service := NewCharacterService(repo)

			err := service.AddItem(context.Background(), tt.character, tt.item)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verifica se o item foi adicionado corretamente
			var found bool
			for _, item := range tt.character.Inventory {
				if item.Name == tt.item.Name {
					found = true
					if item.Quantity != tt.wantQty {
						t.Errorf("AddItem() quantity = %v, want %v", item.Quantity, tt.wantQty)
					}
					break
				}
			}
			if !found && !tt.wantErr {
				t.Errorf("AddItem() item não encontrado no inventário")
			}

			// Verifica se o personagem foi atualizado no repositório
			updated, _ := repo.GetByID(context.Background(), tt.character.ID.Hex())
			found = false
			for _, item := range updated.Inventory {
				if item.Name == tt.item.Name {
					found = true
					if item.Quantity != tt.wantQty {
						t.Errorf("AddItem() saved quantity = %v, want %v", item.Quantity, tt.wantQty)
					}
					break
				}
			}
			if !found && !tt.wantErr {
				t.Errorf("AddItem() item não encontrado no inventário salvo")
			}
		})
	}
}

func TestCharacterService_RemoveItem(t *testing.T) {
	tests := []struct {
		name      string
		character *entities.Character
		itemName  string
		quantity  int
		wantQty   int
		wantErr   bool
	}{
		{
			name: "should remove item partially",
			character: func() *entities.Character {
				c := entities.NewCharacter("123", "456", "Sir Test")
				c.AddItem(entities.Item{Name: "Poção de Cura", Quantity: 3})
				return c
			}(),
			itemName: "Poção de Cura",
			quantity: 2,
			wantQty:  1,
			wantErr:  false,
		},
		{
			name: "should remove item completely",
			character: func() *entities.Character {
				c := entities.NewCharacter("123", "456", "Sir Test")
				c.AddItem(entities.Item{Name: "Poção de Cura", Quantity: 1})
				return c
			}(),
			itemName: "Poção de Cura",
			quantity: 1,
			wantQty:  0,
			wantErr:  false,
		},
		{
			name:      "should error on non-existent item",
			character: entities.NewCharacter("123", "456", "Sir Test"),
			itemName:  "Item Inexistente",
			quantity:  1,
			wantQty:   0,
			wantErr:   true,
		},
		{
			name: "should error on insufficient quantity",
			character: func() *entities.Character {
				c := entities.NewCharacter("123", "456", "Sir Test")
				c.AddItem(entities.Item{Name: "Poção de Cura", Quantity: 1})
				return c
			}(),
			itemName: "Poção de Cura",
			quantity: 2,
			wantQty:  1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockCharacterRepository()
			repo.Create(context.Background(), tt.character)
			service := NewCharacterService(repo)

			err := service.RemoveItem(context.Background(), tt.character, tt.itemName, tt.quantity)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verifica se o item foi removido corretamente
				var found bool
				for _, item := range tt.character.Inventory {
					if item.Name == tt.itemName {
						found = true
						if item.Quantity != tt.wantQty {
							t.Errorf("RemoveItem() quantity = %v, want %v", item.Quantity, tt.wantQty)
						}
						break
					}
				}
				if found && tt.wantQty == 0 {
					t.Errorf("RemoveItem() item ainda está no inventário quando deveria ter sido removido")
				}
				if !found && tt.wantQty > 0 {
					t.Errorf("RemoveItem() item não encontrado no inventário quando deveria ter quantidade %d", tt.wantQty)
				}

				// Verifica se o personagem foi atualizado no repositório
				updated, _ := repo.GetByID(context.Background(), tt.character.ID.Hex())
				found = false
				for _, item := range updated.Inventory {
					if item.Name == tt.itemName {
						found = true
						if item.Quantity != tt.wantQty {
							t.Errorf("RemoveItem() saved quantity = %v, want %v", item.Quantity, tt.wantQty)
						}
						break
					}
				}
				if found && tt.wantQty == 0 {
					t.Errorf("RemoveItem() item ainda está no inventário salvo quando deveria ter sido removido")
				}
				if !found && tt.wantQty > 0 {
					t.Errorf("RemoveItem() item não encontrado no inventário salvo quando deveria ter quantidade %d", tt.wantQty)
				}
			}
		})
	}
}
