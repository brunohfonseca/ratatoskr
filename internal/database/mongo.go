package database

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
