package main

import (
	"flag"
	"log"

	"github.com/bhfonseca/ratatoskr/internal/api"
	"github.com/bhfonseca/ratatoskr/internal/config"
	"github.com/bhfonseca/ratatoskr/internal/database"
)

func main() {
	configFile := flag.String("config", "config.yml", "Arquivo de configuração")
	flag.Parse()

	_, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Erro ao carregar config: %v", err)
	}

	cfg := config.Get()
	if cfg == nil {
		panic(err)
	}

	log.Printf("Iniciando o serviço com o arquivo de configuração: %s", *configFile)

	database.ConnectMongoDB(cfg.Database.MongoURL)
	database.ConnectRedis(cfg.Redis.RedisURL)
	api.ServerStart(cfg)
}
