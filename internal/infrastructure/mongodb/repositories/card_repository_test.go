package repositories

import (
	"context"
	"testing"
	"time"

	"sirdraith/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
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

func TestMongoCardRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewMongoCardRepository(db)

	t.Run("Create", func(t *testing.T) {
		card := entities.NewCard(
			"card1",
			"Test Card",
			entities.CardTypeSpell,
			entities.CardRarityCommon,
			"Test Description",
			2,
			0,
			0,
			[]string{"effect1"},
			[]string{"keyword1"},
			"image.jpg",
		)

		err := repo.Create(context.Background(), card)
		assert.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		card := entities.NewCard(
			"card2",
			"Test Card 2",
			entities.CardTypeCreature,
			entities.CardRarityRare,
			"Test Description 2",
			3,
			2,
			2,
			[]string{"effect2"},
			[]string{"keyword2"},
			"image2.jpg",
		)

		err := repo.Create(context.Background(), card)
		assert.NoError(t, err)

		found, err := repo.FindByID(context.Background(), "card2")
		assert.NoError(t, err)
		assert.Equal(t, card.ID, found.ID)
		assert.Equal(t, card.Name, found.Name)
	})

	t.Run("Update", func(t *testing.T) {
		card := entities.NewCard(
			"card3",
			"Test Card 3",
			entities.CardTypeSpell,
			entities.CardRarityUncommon,
			"Test Description 3",
			1,
			0,
			0,
			[]string{"effect3"},
			[]string{"keyword3"},
			"image3.jpg",
		)

		err := repo.Create(context.Background(), card)
		assert.NoError(t, err)

		card.Name = "Updated Card 3"
		card.UpdatedAt = time.Now().Unix()

		err = repo.Update(context.Background(), card)
		assert.NoError(t, err)

		found, err := repo.FindByID(context.Background(), "card3")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Card 3", found.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		card := entities.NewCard(
			"card4",
			"Test Card 4",
			entities.CardTypeArtifact,
			entities.CardRarityEpic,
			"Test Description 4",
			4,
			0,
			0,
			[]string{"effect4"},
			[]string{"keyword4"},
			"image4.jpg",
		)

		err := repo.Create(context.Background(), card)
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), "card4")
		assert.NoError(t, err)

		_, err = repo.FindByID(context.Background(), "card4")
		assert.Error(t, err)
	})

	t.Run("FindAll", func(t *testing.T) {
		cards, err := repo.FindAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, cards)
	})

	t.Run("FindByType", func(t *testing.T) {
		spellCards, err := repo.FindByType(context.Background(), entities.CardTypeSpell)
		assert.NoError(t, err)
		for _, card := range spellCards {
			assert.Equal(t, entities.CardTypeSpell, card.Type)
		}
	})

	t.Run("FindByRarity", func(t *testing.T) {
		rareCards, err := repo.FindByRarity(context.Background(), entities.CardRarityRare)
		assert.NoError(t, err)
		for _, card := range rareCards {
			assert.Equal(t, entities.CardRarityRare, card.Rarity)
		}
	})
}
