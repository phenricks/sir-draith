package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sirdraith/internal/domain/entities"
)

const characterCollection = "characters"

// CharacterRepository implementa a interface repositories.CharacterRepository
type CharacterRepository struct {
	db *mongo.Database
}

// NewCharacterRepository cria uma nova instância do repositório
func NewCharacterRepository(db *mongo.Database) *CharacterRepository {
	return &CharacterRepository{
		db: db,
	}
}

// Create cria um novo personagem
func (r *CharacterRepository) Create(ctx context.Context, character *entities.Character) error {
	character.CreatedAt = time.Now()
	character.UpdatedAt = time.Now()
	character.IsActive = true

	result, err := r.db.Collection(characterCollection).InsertOne(ctx, character)
	if err != nil {
		return fmt.Errorf("erro ao criar personagem: %w", err)
	}

	character.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update atualiza um personagem existente
func (r *CharacterRepository) Update(ctx context.Context, character *entities.Character) error {
	character.UpdatedAt = time.Now()

	filter := bson.M{"_id": character.ID}
	_, err := r.db.Collection(characterCollection).ReplaceOne(ctx, filter, character)
	if err != nil {
		return fmt.Errorf("erro ao atualizar personagem: %w", err)
	}

	return nil
}

// Delete marca um personagem como inativo
func (r *CharacterRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	filter := bson.M{"_id": objectID}
	result, err := r.db.Collection(characterCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("erro ao deletar personagem: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("personagem não encontrado")
	}

	return nil
}

// GetByID busca um personagem pelo ID
func (r *CharacterRepository) GetByID(ctx context.Context, id string) (*entities.Character, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	var character entities.Character
	filter := bson.M{"_id": objectID, "is_active": true}
	err = r.db.Collection(characterCollection).FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("personagem não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar personagem: %w", err)
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
	err := r.db.Collection(characterCollection).FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("personagem não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar personagem: %w", err)
	}

	return &character, nil
}

// ListByUser lista todos os personagens de um usuário
func (r *CharacterRepository) ListByUser(ctx context.Context, userID string) ([]*entities.Character, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}
	return r.find(ctx, filter)
}

// ListByGuild lista todos os personagens de um servidor
func (r *CharacterRepository) ListByGuild(ctx context.Context, guildID string) ([]*entities.Character, error) {
	filter := bson.M{
		"guild_id":  guildID,
		"is_active": true,
	}
	return r.find(ctx, filter)
}

// ListActive lista todos os personagens ativos
func (r *CharacterRepository) ListActive(ctx context.Context) ([]*entities.Character, error) {
	filter := bson.M{"is_active": true}
	return r.find(ctx, filter)
}

// Search busca personagens por nome ou título
func (r *CharacterRepository) Search(ctx context.Context, query string) ([]*entities.Character, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}},
			{"title": bson.M{"$regex": primitive.Regex{Pattern: query, Options: "i"}}},
		},
		"is_active": true,
	}
	return r.find(ctx, filter)
}

// CountByGuild conta o número de personagens em um servidor
func (r *CharacterRepository) CountByGuild(ctx context.Context, guildID string) (int64, error) {
	filter := bson.M{
		"guild_id":  guildID,
		"is_active": true,
	}
	return r.db.Collection(characterCollection).CountDocuments(ctx, filter)
}

// CountByUser conta o número de personagens de um usuário
func (r *CharacterRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}
	return r.db.Collection(characterCollection).CountDocuments(ctx, filter)
}

// find é um método auxiliar para buscar múltiplos personagens
func (r *CharacterRepository) find(ctx context.Context, filter bson.M) ([]*entities.Character, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(characterCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar personagens: %w", err)
	}
	defer cursor.Close(ctx)

	var characters []*entities.Character
	if err = cursor.All(ctx, &characters); err != nil {
		return nil, fmt.Errorf("erro ao decodificar personagens: %w", err)
	}

	return characters, nil
}

// GetByUserID lista todos os personagens de um usuário
func (r *CharacterRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Character, error) {
	return r.ListByUser(ctx, userID)
}

// GetByGuildID lista todos os personagens de um servidor
func (r *CharacterRepository) GetByGuildID(ctx context.Context, guildID string) ([]*entities.Character, error) {
	return r.ListByGuild(ctx, guildID)
}

// List lista todos os personagens ativos
func (r *CharacterRepository) List(ctx context.Context) ([]*entities.Character, error) {
	return r.ListActive(ctx)
}
