package repositories

import (
	"context"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EndpointRepository interface {
	Create(ctx context.Context, e *entities.Endpoint) (primitive.ObjectID, error)
}

type endpointRepositoryMongo struct {
	col *mongo.Collection
}

func NewEndpointRepositoryMongo(db *mongo.Database) EndpointRepository {
	return &endpointRepositoryMongo{
		col: db.Collection("endpoints"),
	}
}

func (r *endpointRepositoryMongo) Create(ctx context.Context, e *entities.Endpoint) (primitive.ObjectID, error) {
	now := time.Now().UTC()
	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}
	e.CreatedAt = now
	e.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, e)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return e.ID, nil
}
