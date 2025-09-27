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

type endpointRepository struct {
	col *mongo.Collection
}

func NewEndpointRepository(db *mongo.Database) EndpointRepository {
	return &endpointRepository{
		col: db.Collection("endpoints"),
	}
}

func (r *endpointRepository) Create(ctx context.Context, e *entities.Endpoint) (primitive.ObjectID, error) {
	now := time.Now().UTC()

	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	e.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, e)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return e.ID, nil
}
