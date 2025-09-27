package infra

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Msgf("Erro ao conectar no MongoDB: %s", err)
	}
	log.Info().Msg("✅ Connected to MongoDB")

	MongoClient = client
}

func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatal().Msgf("Erro ao desconectar do MongoDB: %v", err)
		} else {
			log.Info().Msg("✅ Disconnected from MongoDB")
		}
	}
}

// CheckMongoDBHealth verifica o status da conexão com MongoDB
func CheckMongoDBHealth() (bool, string, error) {
	if MongoClient == nil {
		return false, "disconnected", nil // Não é um erro fatal, apenas não conectado
	}

	// Criar contexto com timeout curto para health check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Fazer ping para verificar se a conexão está ativa
	err := MongoClient.Ping(ctx, nil)
	if err != nil {
		return false, "error", err
	}

	return true, "connected", nil
}
