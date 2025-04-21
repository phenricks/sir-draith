package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"sirdraith/internal/domain/entities"
)

const characterCollection = "characters"

// CharacterRepository implementa a interface repositories.CharacterRepository
type CharacterRepository struct {
	collection *mongo.Collection
}

// NewCharacterRepository cria uma nova instância do repositório
func NewCharacterRepository(db *mongo.Database) *CharacterRepository {
	return &CharacterRepository{
		collection: db.Collection(characterCollection),
	}
}

// Create cria um novo personagem
func (r *CharacterRepository) Create(ctx context.Context, character *entities.Character) error {
	character.ID = primitive.NewObjectID()
	character.CreatedAt = time.Now()
	character.UpdatedAt = time.Now()
	character.IsActive = true

	result, err := r.collection.InsertOne(ctx, character)
	if err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}

	character.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update atualiza um personagem existente
func (r *CharacterRepository) Update(ctx context.Context, character *entities.Character) error {
	character.UpdatedAt = time.Now()

	filter := bson.M{"_id": character.ID}
	result, err := r.collection.ReplaceOne(ctx, filter, character)
	if err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("character not found")
	}

	return nil
}

// Delete marca um personagem como inativo
func (r *CharacterRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"deleted_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("character not found")
	}

	return nil
}

// GetByID busca um personagem pelo ID
func (r *CharacterRepository) GetByID(ctx context.Context, id string) (*entities.Character, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var character entities.Character
	filter := bson.M{"_id": objectID, "is_active": true}

	err = r.collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	return &character, nil
}

// GetByUserAndGuild busca um personagem pelo ID do usuário e do servidor
func (r *CharacterRepository) GetByUserAndGuild(ctx context.Context, userID, guildID string) (*entities.Character, error) {
	var character entities.Character
	filter := bson.M{
		"user_id":   userID,
		"guild_id":  guildID,
		"is_active": true,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	return &character, nil
}

// ListByUser lista todos os personagens de um usuário
func (r *CharacterRepository) ListByUser(ctx context.Context, userID string) ([]*entities.Character, error) {
	filter := bson.M{"user_id": userID, "is_active": true}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list characters: %w", err)
	}
	defer cursor.Close(ctx)

	var characters []*entities.Character
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, fmt.Errorf("failed to decode characters: %w", err)
	}

	return characters, nil
}

// ListByGuild lista todos os personagens de um servidor
func (r *CharacterRepository) ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error) {
	filter := bson.M{"guild_id": guildID, "is_active": true}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list characters: %w", err)
	}
	defer cursor.Close(ctx)

	var characters []*entities.Character
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, fmt.Errorf("failed to decode characters: %w", err)
	}

	return characters, nil
}

// List returns all active characters
func (r *CharacterRepository) List(ctx context.Context) ([]*entities.Character, error) {
	filter := bson.M{"is_deleted": false}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list characters: %w", err)
	}
	defer cursor.Close(ctx)

	var characters []*entities.Character
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, fmt.Errorf("failed to decode characters: %w", err)
	}

	return characters, nil
}

// Search busca personagens por nome ou título
func (r *CharacterRepository) Search(ctx context.Context, query string) ([]*entities.Character, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"class": bson.M{"$regex": query, "$options": "i"}},
		},
		"is_active": true,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search characters: %w", err)
	}
	defer cursor.Close(ctx)

	var characters []*entities.Character
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, fmt.Errorf("failed to decode characters: %w", err)
	}

	return characters, nil
}

// CountByGuild conta o número de personagens em um servidor
func (r *CharacterRepository) CountByGuild(ctx context.Context, guildID string) (int64, error) {
	filter := bson.M{"guild_id": guildID, "is_active": true}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count characters: %w", err)
	}
	return count, nil
}

// CountByUser conta o número de personagens de um usuário
func (r *CharacterRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"user_id": userID, "is_active": true}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count characters: %w", err)
	}
	return count, nil
}

// GetByUserID lista todos os personagens de um usuário
func (r *CharacterRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Character, error) {
	return r.ListByUser(ctx, userID)
}

// GetByGuildID lista todos os personagens de um servidor
func (r *CharacterRepository) GetByGuildID(ctx context.Context, guildID string) ([]*entities.Character, error) {
	return r.ListByGuild(ctx, guildID)
}
