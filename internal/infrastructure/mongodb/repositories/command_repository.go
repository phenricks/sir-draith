package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"sirdraith/internal/domain/entities"
	"sirdraith/internal/domain/repositories"
)

// MongoCommandRepository implements the CommandRepository interface using MongoDB
type MongoCommandRepository struct {
	collection *mongo.Collection
}

// NewMongoCommandRepository creates a new MongoCommandRepository instance
func NewMongoCommandRepository(db *mongo.Database) repositories.CommandRepository {
	return &MongoCommandRepository{
		collection: db.Collection("commands"),
	}
}

// Create stores a new command in the database
func (r *MongoCommandRepository) Create(ctx context.Context, command *entities.Command) error {
	command.CreatedAt = time.Now()
	command.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, command)
	return err
}

// FindByID retrieves a command by its ID
func (r *MongoCommandRepository) FindByID(ctx context.Context, id string) (*entities.Command, error) {
	var command entities.Command
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&command)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &command, nil
}

// FindByName retrieves a command by its name within a guild
func (r *MongoCommandRepository) FindByName(ctx context.Context, guildID, name string) (*entities.Command, error) {
	var command entities.Command
	err := r.collection.FindOne(ctx, bson.M{
		"guild_id": guildID,
		"name":     name,
	}).Decode(&command)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &command, nil
}

// ListByGuild retrieves all commands for a specific guild
func (r *MongoCommandRepository) ListByGuild(ctx context.Context, guildID string) ([]*entities.Command, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"guild_id": guildID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commands []*entities.Command
	if err = cursor.All(ctx, &commands); err != nil {
		return nil, err
	}
	return commands, nil
}

// Update updates an existing command in the database
func (r *MongoCommandRepository) Update(ctx context.Context, command *entities.Command) error {
	command.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": command.ID}, command)
	return err
}

// Delete removes a command from the database
func (r *MongoCommandRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
