package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sirdraith/internal/domain/model"
	"sirdraith/internal/domain/repository"
)

const configCollection = "guild_configs"

// MongoConfigRepository implementa ConfigRepository usando MongoDB
type MongoConfigRepository struct {
	db *mongo.Database
}

// NewConfigRepository cria uma nova instância do repositório de configurações
func NewConfigRepository(db *mongo.Database) repository.ConfigRepository {
	return &MongoConfigRepository{db: db}
}

// collection retorna a coleção de configurações
func (r *MongoConfigRepository) collection() *mongo.Collection {
	return r.db.Collection(configCollection)
}

// GetGuildConfig implementa repository.ConfigRepository
func (r *MongoConfigRepository) GetGuildConfig(guildID string) (*model.GuildConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config model.GuildConfig
	err := r.collection().FindOne(ctx, bson.M{"_id": guildID}).Decode(&config)
	if err == mongo.ErrNoDocuments {
		return model.NewGuildConfig(guildID), nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configuração do servidor: %w", err)
	}

	return &config, nil
}

// UpdateGuildPrefix implementa repository.ConfigRepository
func (r *MongoConfigRepository) UpdateGuildPrefix(guildID string, newPrefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"prefix":     newPrefix,
			"updated_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": guildID}, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao atualizar prefixo do servidor: %w", err)
	}

	return nil
}

// EnsureGuildConfig implementa repository.ConfigRepository
func (r *MongoConfigRepository) EnsureGuildConfig(guildID string) (*model.GuildConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := model.NewGuildConfig(guildID)
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection().ReplaceOne(ctx, bson.M{"_id": guildID}, config, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao garantir configuração do servidor: %w", err)
	}

	return config, nil
}

// UpdateGuildConfig atualiza a configuração completa de um servidor
func (r *MongoConfigRepository) UpdateGuildConfig(guildID string, config *model.GuildConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"prefix":          config.Prefix,
			"welcome_channel": config.WelcomeChannel,
			"goodbye_channel": config.GoodbyeChannel,
			"updated_at":      config.UpdatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": guildID}, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao atualizar configuração do servidor: %w", err)
	}

	return nil
}
