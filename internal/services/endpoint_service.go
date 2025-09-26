package services

import (
	"context"
	"errors"
	"time"

	mongo2 "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/mongo"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EndpointService struct{}

// NewEndpointService cria uma nova instância do EndpointService
func NewEndpointService() *EndpointService {
	return &EndpointService{}
}

// getCollection retorna a collection de endpoints
func (s *EndpointService) getCollection() *mongo.Collection {
	return mongo2.MongoClient.Database("ratatoskr").Collection("endpoints")
}

// CreateEndpoint cria um novo endpoint
func (s *EndpointService) CreateEndpoint(endpoint *models.Endpoint) (*models.Endpoint, error) {
	// Validações básicas
	if endpoint.Name == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if endpoint.Domain == "" {
		return nil, errors.New("domínio é obrigatório")
	}

	// Verificar se já existe
	exists, err := s.ExistsByName(endpoint.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("já existe um endpoint com este nome")
	}

	// Preparar dados
	now := time.Now()
	endpoint.ID = primitive.NewObjectID()
	endpoint.CreatedAt = now
	endpoint.UpdatedAt = now
	endpoint.Status = models.StatusUnknown
	endpoint.Enabled = true

	// Valores padrão
	if endpoint.Timeout == 0 {
		endpoint.Timeout = 30 * time.Second
	}
	if endpoint.Interval == 0 {
		endpoint.Interval = 5 * time.Minute
	}

	// Inserir
	ctx := context.Background()
	_, err = s.getCollection().InsertOne(ctx, endpoint)
	if err != nil {
		log.Error().Err(err).Str("name", endpoint.Name).Msg("Erro ao criar endpoint")
		return nil, err
	}

	log.Info().Str("id", endpoint.ID.Hex()).Str("name", endpoint.Name).Msg("Endpoint criado")
	return endpoint, nil
}

// GetByID busca endpoint por ID
func (s *EndpointService) GetByID(id string) (*models.Endpoint, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	var endpoint models.Endpoint
	ctx := context.Background()
	err = s.getCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&endpoint)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("endpoint não encontrado")
		}
		return nil, err
	}

	return &endpoint, nil
}

// Update atualiza um endpoint
func (s *EndpointService) Update(id string, endpoint *models.Endpoint) (*models.Endpoint, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	// Validações
	if endpoint.Name == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if endpoint.Domain == "" {
		return nil, errors.New("domínio é obrigatório")
	}

	// Verificar se existe
	_, err = s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update
	endpoint.UpdatedAt = time.Now()
	ctx := context.Background()
	_, err = s.getCollection().UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": endpoint})
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Erro ao atualizar endpoint")
		return nil, err
	}

	// Retornar atualizado
	return s.GetByID(id)
}

// Delete remove um endpoint
func (s *EndpointService) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	// Verificar se existe
	endpoint, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Deletar
	ctx := context.Background()
	_, err = s.getCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Erro ao deletar endpoint")
		return err
	}

	log.Info().Str("id", id).Str("name", endpoint.Name).Msg("Endpoint removido")
	return nil
}

// List retorna todos os endpoints
func (s *EndpointService) List() ([]*models.Endpoint, error) {
	ctx := context.Background()
	cursor, err := s.getCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var endpoints []*models.Endpoint
	for cursor.Next(ctx) {
		var endpoint models.Endpoint
		if err := cursor.Decode(&endpoint); err != nil {
			log.Error().Err(err).Msg("Erro ao decodificar endpoint")
			continue
		}
		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, nil
}

// ExistsByName verifica se existe endpoint com o nome
func (s *EndpointService) ExistsByName(name string) (bool, error) {
	ctx := context.Background()
	count, err := s.getCollection().CountDocuments(ctx, bson.M{"name": name})
	return count > 0, err
}

// UpdateStatus atualiza apenas o status
func (s *EndpointService) UpdateStatus(id string, status models.EndpointStatus, responseTime time.Duration, errorMsg string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	update := bson.M{
		"$set": bson.M{
			"status":        status,
			"response_time": responseTime,
			"error_message": errorMsg,
			"last_check":    time.Now(),
			"updated_at":    time.Now(),
		},
	}

	ctx := context.Background()
	_, err = s.getCollection().UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
