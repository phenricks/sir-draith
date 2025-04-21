package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// MongoCardRepository implements CardRepository interface using MongoDB
type MongoCardRepository struct {
	collection *mongo.Collection
}

// NewMongoCardRepository creates a new MongoDB card repository
func NewMongoCardRepository(db *mongo.Database) repositories.CardRepository {
	return &MongoCardRepository{
		collection: db.Collection("cards"),
	}
}

// Create stores a new card in MongoDB
func (r *MongoCardRepository) Create(ctx context.Context, card *entities.Card) error {
	_, err := r.collection.InsertOne(ctx, card)
	return err
}

// Update modifies an existing card in MongoDB
func (r *MongoCardRepository) Update(ctx context.Context, card *entities.Card) error {
	filter := bson.M{"_id": card.ID}
	update := bson.M{"$set": card}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a card from MongoDB
func (r *MongoCardRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// FindByID retrieves a card by its ID from MongoDB
func (r *MongoCardRepository) FindByID(ctx context.Context, id string) (*entities.Card, error) {
	filter := bson.M{"_id": id}
	var card entities.Card
	err := r.collection.FindOne(ctx, filter).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

// FindAll retrieves all cards from MongoDB
func (r *MongoCardRepository) FindAll(ctx context.Context) ([]*entities.Card, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*entities.Card
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// FindByType retrieves all cards of a specific type from MongoDB
func (r *MongoCardRepository) FindByType(ctx context.Context, cardType entities.CardType) ([]*entities.Card, error) {
	filter := bson.M{"type": cardType}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*entities.Card
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// FindByRarity retrieves all cards of a specific rarity from MongoDB
func (r *MongoCardRepository) FindByRarity(ctx context.Context, rarity entities.CardRarity) ([]*entities.Card, error) {
	filter := bson.M{"rarity": rarity}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*entities.Card
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}
