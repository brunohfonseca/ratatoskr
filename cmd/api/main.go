package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunohfonseca/ratatoskr/internal/api"
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/database"
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
		log.Fatalf("configuração não carregada: %v", err)
	}

	log.Printf("Iniciando o serviço com o arquivo de configuração: %s", *configFile)

	database.ConnectMongoDB(cfg.Database.MongoURL)
	database.ConnectRedis(cfg.Redis.RedisURL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	srv := api.ServerStart(cfg)
	// inicia servidor em goroutine
	go func() {
		if cfg.Server.SSL.Enabled {
			if err := srv.ListenAndServeTLS(cfg.Server.SSL.Cert, cfg.Server.SSL.Key); err != nil && err != http.ErrServerClosed {
				log.Fatalf("erro ao iniciar servidor TLS: %v", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("erro ao iniciar servidor: %v", err)
			}
		}
	}()

	fmt.Println("🚀 Servidor iniciado! Pressione Ctrl+C para finalizar.")

	// Aguardar sinal de parada
	<-c
	fmt.Println("\n🛑 Sinal de parada recebido. Finalizando aplicação...")

	database.DisconnectMongoDB()
	database.DisconnectRedis()

	fmt.Println("✅ Aplicação finalizada com sucesso!")
}
