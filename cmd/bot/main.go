package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sirdraith/internal/infrastructure/discord"
)

func main() {
	log.Println("Iniciando Sir Draith - Bot Medieval...")

	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Arquivo .env não encontrado: %v\n", err)
	}

	// Obtém o token do Discord das variáveis de ambiente
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Token do Discord não encontrado nas variáveis de ambiente")
	}

	// Obtém a URL do MongoDB e credenciais
	mongoURL := os.Getenv("MONGODB_URL")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPass := os.Getenv("MONGODB_PASS")
	mongoAuthSource := os.Getenv("MONGODB_AUTH_SOURCE")
	if mongoAuthSource == "" {
		mongoAuthSource = "admin"
	}

	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017" // URL padrão
	}

	// Configura as opções de conexão do MongoDB
	opts := options.Client().ApplyURI(mongoURL)
	if mongoUser != "" && mongoPass != "" {
		opts.SetAuth(options.Credential{
			Username:   mongoUser,
			Password:   mongoPass,
			AuthSource: mongoAuthSource,
		})
	}

	// Conecta ao MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Verifica a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Erro ao verificar conexão com MongoDB: %v", err)
	}

	// Seleciona o banco de dados
	db := client.Database("sirdraith")

	// Cria e configura o cliente Discord
	discordClient, err := discord.NewClient(token, db)
	if err != nil {
		log.Fatalf("Erro ao criar cliente Discord: %v", err)
	}

	// Conecta ao Discord
	if err := discordClient.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao Discord: %v", err)
	}

	// Configura graceful shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Desconecta do Discord
	if err := discordClient.Disconnect(); err != nil {
		log.Printf("Erro ao desconectar do Discord: %v", err)
	}

	log.Println("Encerrando Sir Draith...")
} 