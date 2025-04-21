package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/gamedata"
)

type CharacterRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repository *CharacterRepository
	ctx        context.Context
}

func TestCharacterRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CharacterRepositoryTestSuite))
}

func (s *CharacterRepositoryTestSuite) SetupSuite() {
	// Conectar ao MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://sirdraith:sirdraith123@mongodb:27017"))
	if err != nil {
		s.T().Fatal(err)
	}

	// Usar banco de dados de teste
	s.db = client.Database("sirdraith_test")
	s.repository = NewCharacterRepository(s.db)
	s.ctx = ctx
}

func (s *CharacterRepositoryTestSuite) TearDownSuite() {
	// Limpar banco de dados de teste
	if err := s.db.Drop(s.ctx); err != nil {
		s.T().Fatal(err)
	}
}

func (s *CharacterRepositoryTestSuite) SetupTest() {
	// Limpar coleção antes de cada teste
	if err := s.db.Collection(characterCollection).Drop(s.ctx); err != nil {
		s.T().Fatal(err)
	}
}

func (s *CharacterRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name      string
		character *entities.Character
		wantErr   bool
	}{
		{
			name: "should create valid character",
			character: &entities.Character{
				UserID:  "123456789",
				GuildID: "987654321",
				Name:    "Sir Test",
				Class:   gamedata.Warrior,
				Level:   1,
			},
			wantErr: false,
		},
		{
			name: "should create character with minimum attributes",
			character: &entities.Character{
				UserID:  "123456789",
				GuildID: "987654321",
				Name:    "Test Minimum",
				Class:   gamedata.Warrior,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repository.Create(s.ctx, tt.character)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.NotEmpty(tt.character.ID)
			s.True(tt.character.IsActive)
			s.False(tt.character.CreatedAt.IsZero())
			s.False(tt.character.UpdatedAt.IsZero())

			// Verificar se o personagem foi salvo corretamente
			var saved entities.Character
			err = s.db.Collection(characterCollection).FindOne(s.ctx, bson.M{"_id": tt.character.ID}).Decode(&saved)
			s.NoError(err)
			s.Equal(tt.character.UserID, saved.UserID)
			s.Equal(tt.character.GuildID, saved.GuildID)
			s.Equal(tt.character.Name, saved.Name)
			s.Equal(tt.character.Class, saved.Class)
			s.Equal(tt.character.Level, saved.Level)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestUpdate() {
	// Criar personagem para teste
	character := &entities.Character{
		UserID:  "123456789",
		GuildID: "987654321",
		Name:    "Sir Test",
		Class:   gamedata.Warrior,
		Level:   1,
	}
	err := s.repository.Create(s.ctx, character)
	s.NoError(err)

	tests := []struct {
		name      string
		character *entities.Character
		updates   func(*entities.Character)
		wantErr   bool
	}{
		{
			name:      "should update character name",
			character: character,
			updates: func(c *entities.Character) {
				c.Name = "Sir Test Updated"
			},
			wantErr: false,
		},
		{
			name:      "should update character level",
			character: character,
			updates: func(c *entities.Character) {
				c.Level = 2
				c.Experience = 200
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			oldUpdatedAt := tt.character.UpdatedAt

			// Aplicar atualizações
			tt.updates(tt.character)
			time.Sleep(time.Millisecond) // Garantir que UpdatedAt será diferente

			err := s.repository.Update(s.ctx, tt.character)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.True(tt.character.UpdatedAt.After(oldUpdatedAt))

			// Verificar se o personagem foi atualizado corretamente
			var saved entities.Character
			err = s.db.Collection(characterCollection).FindOne(s.ctx, bson.M{"_id": tt.character.ID}).Decode(&saved)
			s.NoError(err)
			s.Equal(tt.character.Name, saved.Name)
			s.Equal(tt.character.Level, saved.Level)
			s.Equal(tt.character.Experience, saved.Experience)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestDelete() {
	// Criar personagem para teste
	character := &entities.Character{
		UserID:  "123456789",
		GuildID: "987654321",
		Name:    "Sir Test",
		Class:   gamedata.Warrior,
		Level:   1,
	}
	err := s.repository.Create(s.ctx, character)
	s.NoError(err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "should delete existing character",
			id:      character.ID.Hex(),
			wantErr: false,
		},
		{
			name:    "should return error for invalid id",
			id:      "invalid-id",
			wantErr: true,
		},
		{
			name:    "should return error for non-existent character",
			id:      primitive.NewObjectID().Hex(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repository.Delete(s.ctx, tt.id)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)

			// Verificar se o personagem foi marcado como inativo
			var saved entities.Character
			err = s.db.Collection(characterCollection).FindOne(s.ctx, bson.M{"_id": character.ID}).Decode(&saved)
			s.NoError(err)
			s.False(saved.IsActive)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestGetByID() {
	// Criar personagem para teste
	character := &entities.Character{
		UserID:  "123456789",
		GuildID: "987654321",
		Name:    "Sir Test",
		Class:   gamedata.Warrior,
		Level:   1,
	}
	err := s.repository.Create(s.ctx, character)
	s.NoError(err)

	tests := []struct {
		name    string
		id      string
		want    *entities.Character
		wantErr bool
	}{
		{
			name:    "should get existing character",
			id:      character.ID.Hex(),
			want:    character,
			wantErr: false,
		},
		{
			name:    "should return error for invalid id",
			id:      "invalid-id",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error for non-existent character",
			id:      primitive.NewObjectID().Hex(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.GetByID(s.ctx, tt.id)

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
				return
			}

			s.NoError(err)
			s.NotNil(got)
			s.Equal(tt.want.ID, got.ID)
			s.Equal(tt.want.UserID, got.UserID)
			s.Equal(tt.want.GuildID, got.GuildID)
			s.Equal(tt.want.Name, got.Name)
			s.Equal(tt.want.Class, got.Class)
			s.Equal(tt.want.Level, got.Level)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestGetByUserAndGuild() {
	// Criar personagens para teste
	character1 := &entities.Character{
		UserID:  "123456789",
		GuildID: "987654321",
		Name:    "Sir Test 1",
		Class:   gamedata.Warrior,
		Level:   1,
	}
	err := s.repository.Create(s.ctx, character1)
	s.NoError(err)

	character2 := &entities.Character{
		UserID:  "123456789",
		GuildID: "987654322",
		Name:    "Sir Test 2",
		Class:   gamedata.Mage,
		Level:   1,
	}
	err = s.repository.Create(s.ctx, character2)
	s.NoError(err)

	tests := []struct {
		name    string
		userID  string
		guildID string
		want    *entities.Character
		wantErr bool
	}{
		{
			name:    "should get character by user and guild",
			userID:  "123456789",
			guildID: "987654321",
			want:    character1,
			wantErr: false,
		},
		{
			name:    "should get different character by user and guild",
			userID:  "123456789",
			guildID: "987654322",
			want:    character2,
			wantErr: false,
		},
		{
			name:    "should return error for non-existent character",
			userID:  "123456789",
			guildID: "987654323",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.GetByUserAndGuild(s.ctx, tt.userID, tt.guildID)

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
				return
			}

			s.NoError(err)
			s.NotNil(got)
			s.Equal(tt.want.ID, got.ID)
			s.Equal(tt.want.UserID, got.UserID)
			s.Equal(tt.want.GuildID, got.GuildID)
			s.Equal(tt.want.Name, got.Name)
			s.Equal(tt.want.Class, got.Class)
			s.Equal(tt.want.Level, got.Level)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestListByUser() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test 1",
			Class:   gamedata.Warrior,
			Level:   1,
		},
		{
			UserID:  "123456789",
			GuildID: "987654322",
			Name:    "Sir Test 2",
			Class:   gamedata.Mage,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654321",
			Name:    "Sir Test 3",
			Class:   gamedata.Warrior,
			Level:   1,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		userID  string
		want    int
		wantErr bool
	}{
		{
			name:    "should list characters by user",
			userID:  "123456789",
			want:    2,
			wantErr: false,
		},
		{
			name:    "should return empty list for non-existent user",
			userID:  "000000000",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.ListByUser(s.ctx, tt.userID)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Len(got, tt.want)

			if tt.want > 0 {
				for _, c := range got {
					s.Equal(tt.userID, c.UserID)
				}
			}
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestSearch() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test Knight",
			Class:   gamedata.Warrior,
			Level:   1,
		},
		{
			UserID:  "123456789",
			GuildID: "987654322",
			Name:    "Lady Test Mage",
			Class:   gamedata.Mage,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654321",
			Name:    "Master Test",
			Class:   gamedata.Warrior,
			Level:   1,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		query   string
		want    int
		wantErr bool
	}{
		{
			name:    "should find characters by name",
			query:   "Knight",
			want:    1,
			wantErr: false,
		},
		{
			name:    "should find characters by partial name",
			query:   "Test",
			want:    3,
			wantErr: false,
		},
		{
			name:    "should return empty list for non-matching query",
			query:   "NonExistent",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.Search(s.ctx, tt.query)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Len(got, tt.want)

			if tt.want > 0 {
				for _, c := range got {
					s.Contains(c.Name, "Test")
				}
			}
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestCountByGuild() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test 1",
			Class:   gamedata.Warrior,
			Level:   1,
		},
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test 2",
			Class:   gamedata.Mage,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654322",
			Name:    "Sir Test 3",
			Class:   gamedata.Warrior,
			Level:   1,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		guildID string
		want    int64
		wantErr bool
	}{
		{
			name:    "should count characters in guild",
			guildID: "987654321",
			want:    2,
			wantErr: false,
		},
		{
			name:    "should return zero for empty guild",
			guildID: "000000000",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.CountByGuild(s.ctx, tt.guildID)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Equal(tt.want, got)
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestListByGuild() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test 1",
			Class:   gamedata.Warrior,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654321",
			Name:    "Sir Test 2",
			Class:   gamedata.Mage,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654322",
			Name:    "Sir Test 3",
			Class:   gamedata.Warrior,
			Level:   1,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		guildID string
		want    int
		wantErr bool
	}{
		{
			name:    "should list characters by guild",
			guildID: "987654321",
			want:    2,
			wantErr: false,
		},
		{
			name:    "should return empty list for non-existent guild",
			guildID: "000000000",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.ListByGuild(s.ctx, tt.guildID)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Len(got, tt.want)

			if tt.want > 0 {
				for _, c := range got {
					s.Equal(tt.guildID, c.GuildID)
				}
			}
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestListActive() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:   "123456789",
			GuildID:  "987654321",
			Name:     "Sir Test 1",
			Class:    gamedata.Warrior,
			Level:    1,
			IsActive: true,
		},
		{
			UserID:   "987654321",
			GuildID:  "987654321",
			Name:     "Sir Test 2",
			Class:    gamedata.Mage,
			Level:    1,
			IsActive: true,
		},
		{
			UserID:   "987654321",
			GuildID:  "987654322",
			Name:     "Sir Test 3",
			Class:    gamedata.Warrior,
			Level:    1,
			IsActive: true,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "should list only active characters",
			want:    2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.ListActive(s.ctx)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Len(got, tt.want)

			for _, c := range got {
				s.True(c.IsActive)
			}
		})
	}
}

func (s *CharacterRepositoryTestSuite) TestCountByUser() {
	// Criar personagens para teste
	characters := []*entities.Character{
		{
			UserID:  "123456789",
			GuildID: "987654321",
			Name:    "Sir Test 1",
			Class:   gamedata.Warrior,
			Level:   1,
		},
		{
			UserID:  "123456789",
			GuildID: "987654322",
			Name:    "Sir Test 2",
			Class:   gamedata.Mage,
			Level:   1,
		},
		{
			UserID:  "987654321",
			GuildID: "987654321",
			Name:    "Sir Test 3",
			Class:   gamedata.Warrior,
			Level:   1,
		},
	}

	for _, c := range characters {
		err := s.repository.Create(s.ctx, c)
		s.NoError(err)
	}

	tests := []struct {
		name    string
		userID  string
		want    int64
		wantErr bool
	}{
		{
			name:    "should count characters by user",
			userID:  "123456789",
			want:    2,
			wantErr: false,
		},
		{
			name:    "should return zero for non-existent user",
			userID:  "000000000",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repository.CountByUser(s.ctx, tt.userID)

			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Equal(tt.want, got)
		})
	}
}
