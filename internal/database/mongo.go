package database

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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

// Registry para structs automáticas
var registeredModels []interface{}

// RegisterModel - Registra uma struct para sincronização automática
func RegisterModel(models ...interface{}) {
	registeredModels = append(registeredModels, models...)
}

// AutoSync - Sincroniza automaticamente todas as structs registradas
func AutoSync(mongoURL string) {
	if len(registeredModels) == 0 {
		log.Warn().Msg("⚠️  Nenhuma model registrada para sincronização")
		return
	}

	// Extrair database name da URL ou usar padrão
	databaseName := extractDatabaseFromURL(mongoURL)
	if databaseName == "" {
		databaseName = "ratatoskr" // database padrão
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := MongoClient.Database(databaseName)
	log.Info().Msgf("🔄 Iniciando sincronização automática no database: %s", databaseName)

	for _, model := range registeredModels {
		syncModel(ctx, db, model)
	}

	log.Info().Msgf("✅ Sincronização automática concluída - %d models processadas", len(registeredModels))
}

// Extrai o nome do database da URL do MongoDB
func extractDatabaseFromURL(mongoURL string) string {
	// Exemplo: mongodb://localhost:27017/ratatoskr → ratatoskr
	// Exemplo: mongodb://user:pass@host:port/dbname → dbname

	// Procurar pela última barra
	parts := strings.Split(mongoURL, "/")
	if len(parts) >= 4 {
		dbPart := parts[len(parts)-1]
		// Remover query parameters se existir
		if strings.Contains(dbPart, "?") {
			dbPart = strings.Split(dbPart, "?")[0]
		}
		if dbPart != "" {
			return dbPart
		}
	}

	return "" // Retorna vazio se não conseguir extrair
}

func syncModel(ctx context.Context, db *mongo.Database, model interface{}) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Gerar nome da collection automaticamente
	collectionName := getCollectionName(modelType.Name())
	collection := db.Collection(collectionName)

	log.Info().Msgf("📄 Processando model: %s → collection: %s", modelType.Name(), collectionName)

	// Criar índices baseado nas tags
	indexes := extractIndexes(modelType)
	if len(indexes) > 0 {
		_, err := collection.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			log.Error().Err(err).Msgf("Erro ao criar índices para %s", collectionName)
		} else {
			log.Info().Msgf("✅ %d índices criados para '%s'", len(indexes), collectionName)
		}
	}

	// Verificar se tem TTL
	ttlIndex := extractTTL(modelType)
	if ttlIndex != nil {
		_, err := collection.Indexes().CreateOne(ctx, *ttlIndex)
		if err != nil {
			log.Error().Err(err).Msgf("Erro ao criar TTL para %s", collectionName)
		} else {
			log.Info().Msgf("⏰ TTL configurado para '%s'", collectionName)
		}
	}
}

// Converte nome da struct para nome da collection
func getCollectionName(structName string) string {
	// ServiceHealthHistory → service_health_histories
	name := camelToSnake(structName)

	// Pluralizar (regra simples)
	if strings.HasSuffix(name, "y") {
		name = strings.TrimSuffix(name, "y") + "ies"
	} else if strings.HasSuffix(name, "s") {
		name = name + "es"
	} else {
		name = name + "s"
	}

	return name
}

// Converte CamelCase para snake_case
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}

		// Converter para minúscula apenas se for maiúscula
		if 'A' <= r && r <= 'Z' {
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Extrai índices das tags da struct
func extractIndexes(modelType reflect.Type) []mongo.IndexModel {
	var indexes []mongo.IndexModel

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Buscar tag 'index'
		indexTag := field.Tag.Get("index")
		if indexTag == "" {
			continue
		}

		bsonName := field.Tag.Get("bson")
		if bsonName == "" {
			bsonName = strings.ToLower(field.Name)
		} else {
			// Extrair nome do campo da tag bson (ex: "name,omitempty" → "name")
			bsonName = strings.Split(bsonName, ",")[0]
		}

		// Criar índice baseado na tag
		index := createIndexFromTag(bsonName, indexTag)
		if index != nil {
			indexes = append(indexes, *index)
		}
	}

	return indexes
}

// Cria índice baseado na tag
func createIndexFromTag(fieldName, tag string) *mongo.IndexModel {
	parts := strings.Split(tag, ",")

	opts := options.Index()
	keys := bson.D{{fieldName, 1}} // Default ascending

	for _, part := range parts {
		part = strings.TrimSpace(part)

		switch part {
		case "unique":
			opts.SetUnique(true)
		case "desc", "-1":
			keys = bson.D{{fieldName, -1}}
		case "text":
			keys = bson.D{{fieldName, "text"}}
		}
	}

	indexName := "idx_" + fieldName
	if opts.Unique != nil && *opts.Unique {
		indexName += "_unique"
	}
	opts.SetName(indexName)

	return &mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}
}

// Extrai configuração de TTL das tags
func extractTTL(modelType reflect.Type) *mongo.IndexModel {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		ttlTag := field.Tag.Get("ttl")
		if ttlTag == "" {
			continue
		}

		bsonName := field.Tag.Get("bson")
		if bsonName == "" {
			bsonName = strings.ToLower(field.Name)
		} else {
			bsonName = strings.Split(bsonName, ",")[0]
		}

		// Converter TTL (ex: "90d", "30d", "7d")
		var seconds int32
		if strings.HasSuffix(ttlTag, "d") {
			days := strings.TrimSuffix(ttlTag, "d")
			if days == "90" {
				seconds = 90 * 24 * 60 * 60
			} else if days == "30" {
				seconds = 30 * 24 * 60 * 60
			} else if days == "7" {
				seconds = 7 * 24 * 60 * 60
			}
		}

		if seconds > 0 {
			return &mongo.IndexModel{
				Keys: bson.D{{bsonName, 1}},
				Options: options.Index().
					SetExpireAfterSeconds(seconds).
					SetName("idx_" + bsonName + "_ttl"),
			}
		}
	}
	return nil
}
