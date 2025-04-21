package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sirdraith/internal/domain/entities"
)

func setupCommandTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	db := client.Database("test_db")

	return db, func() {
		err := db.Drop(ctx)
		assert.NoError(t, err)
		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	}
}

func TestMongoCommandRepository(t *testing.T) {
	db, cleanup := setupCommandTestDB(t)
	defer cleanup()

	repo := NewMongoCommandRepository(db)
	ctx := context.Background()

	// Clean up after tests
	defer func() {
		db.Collection("commands").Drop(ctx)
	}()

	t.Run("Create and FindByID", func(t *testing.T) {
		command := &entities.Command{
			ID:          "cmd1",
			GuildID:     "guild1",
			Name:        "test",
			Description: "test command",
			Type:        entities.TextCommand,
			Response:    "test response",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		err := repo.Create(ctx, command)
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, command.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, command.ID, found.ID)
		assert.Equal(t, command.Name, found.Name)
		assert.Equal(t, command.Description, found.Description)
		assert.Equal(t, command.Type, found.Type)
		assert.Equal(t, command.Response, found.Response)
		assert.Equal(t, command.CreatedBy, found.CreatedBy)
		assert.Equal(t, command.Enabled, found.Enabled)
		assert.NotZero(t, found.CreatedAt)
		assert.NotZero(t, found.UpdatedAt)
	})

	t.Run("FindByName", func(t *testing.T) {
		command := &entities.Command{
			ID:          "cmd2",
			GuildID:     "guild1",
			Name:        "unique",
			Description: "unique command",
			Type:        entities.CustomCommand,
			Response:    "unique response",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		err := repo.Create(ctx, command)
		assert.NoError(t, err)

		found, err := repo.FindByName(ctx, command.GuildID, command.Name)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, command.ID, found.ID)
		assert.Equal(t, command.Name, found.Name)
	})

	t.Run("ListByGuild", func(t *testing.T) {
		command1 := &entities.Command{
			ID:          "cmd3",
			GuildID:     "guild2",
			Name:        "test1",
			Description: "test command 1",
			Type:        entities.TextCommand,
			Response:    "response 1",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		command2 := &entities.Command{
			ID:          "cmd4",
			GuildID:     "guild2",
			Name:        "test2",
			Description: "test command 2",
			Type:        entities.CustomCommand,
			Response:    "response 2",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		err := repo.Create(ctx, command1)
		assert.NoError(t, err)
		err = repo.Create(ctx, command2)
		assert.NoError(t, err)

		commands, err := repo.ListByGuild(ctx, "guild2")
		assert.NoError(t, err)
		assert.Len(t, commands, 2)
	})

	t.Run("Update", func(t *testing.T) {
		command := &entities.Command{
			ID:          "cmd5",
			GuildID:     "guild1",
			Name:        "update",
			Description: "update command",
			Type:        entities.TextCommand,
			Response:    "update response",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		err := repo.Create(ctx, command)
		assert.NoError(t, err)

		command.Description = "updated description"
		command.Response = "updated response"
		command.Enabled = false

		err = repo.Update(ctx, command)
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, command.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "updated description", found.Description)
		assert.Equal(t, "updated response", found.Response)
		assert.False(t, found.Enabled)
		assert.True(t, found.UpdatedAt.After(found.CreatedAt))
	})

	t.Run("Delete", func(t *testing.T) {
		command := &entities.Command{
			ID:          "cmd6",
			GuildID:     "guild1",
			Name:        "delete",
			Description: "delete command",
			Type:        entities.TextCommand,
			Response:    "delete response",
			CreatedBy:   "user1",
			Enabled:     true,
		}

		err := repo.Create(ctx, command)
		assert.NoError(t, err)

		err = repo.Delete(ctx, command.ID)
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, command.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}
