package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// MongoDeckRepository implementa a interface DeckRepository usando MongoDB
type MongoDeckRepository struct {
	collection *mongo.Collection
}

// NewMongoDeckRepository cria um novo repositório de decks MongoDB
func NewMongoDeckRepository(db *mongo.Database) repositories.DeckRepository {
	return &MongoDeckRepository{
		collection: db.Collection("decks"),
	}
}

// Create armazena um novo deck no MongoDB
func (r *MongoDeckRepository) Create(ctx context.Context, deck *entities.Deck) error {
	_, err := r.collection.InsertOne(ctx, deck)
	return err
}

// Update atualiza um deck existente no MongoDB
func (r *MongoDeckRepository) Update(ctx context.Context, deck *entities.Deck) error {
	filter := bson.M{"_id": deck.ID}
	update := bson.M{"$set": deck}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete remove um deck do MongoDB
func (r *MongoDeckRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// FindByID busca um deck pelo ID no MongoDB
func (r *MongoDeckRepository) FindByID(ctx context.Context, id string) (*entities.Deck, error) {
	filter := bson.M{"_id": id}
	var deck entities.Deck
	err := r.collection.FindOne(ctx, filter).Decode(&deck)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &deck, nil
}

// FindByUser busca todos os decks de um usuário no MongoDB
func (r *MongoDeckRepository) FindByUser(ctx context.Context, userID string) ([]*entities.Deck, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []*entities.Deck
	if err = cursor.All(ctx, &decks); err != nil {
		return nil, err
	}
	return decks, nil
}

// FindByGuild busca todos os decks de um servidor no MongoDB
func (r *MongoDeckRepository) FindByGuild(ctx context.Context, guildID string) ([]*entities.Deck, error) {
	filter := bson.M{"guild_id": guildID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []*entities.Deck
	if err = cursor.All(ctx, &decks); err != nil {
		return nil, err
	}
	return decks, nil
}

// FindByClass busca todos os decks de uma classe no MongoDB
func (r *MongoDeckRepository) FindByClass(ctx context.Context, class string) ([]*entities.Deck, error) {
	filter := bson.M{"class": class}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []*entities.Deck
	if err = cursor.All(ctx, &decks); err != nil {
		return nil, err
	}
	return decks, nil
}
