package repository

import (
	"athena/domain"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ITrackerRepository interface {
	InsertManyVehicles(ctx context.Context, vehicles []domain.Dado) error
	InsertOneVehicle(ctx context.Context, vehicle domain.Dado) error
	GetVehicles(ctx context.Context) ([]domain.Dado, error)
}

type TrackerRepository struct {
	coll *mongo.Collection
}

// GetVehicles implements ITrackerRepository.
func (v TrackerRepository) GetVehicles(ctx context.Context) ([]domain.Dado, error) {
	var results []domain.Dado

	// obs: Mongo retornar um cursor, não os objetos diretamente, isso ajuda a evitar carregar
	// tudo na memoria de uma vez, bem interessante xD
	cur, err := v.coll.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar todos os dados da collection Vehicles: %w", err)
	}
	defer cur.Close(context.Background())

	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("erro ao decodificar veículos: %w", err)
	}

	return results, nil
}

// InsertManyVehicles implements IVehicleRepository.
func (v TrackerRepository) InsertManyVehicles(ctx context.Context, vehicles []domain.Dado) error {
	_, err := v.coll.InsertMany(ctx, vehicles)
	if err != nil {
		return fmt.Errorf("erro ao inserir varias posições de veículos no DB: %w", err)
	}
	return nil
}

// InsertOneVehicle implements IVehicleRepository.
func (v *TrackerRepository) InsertOneVehicle(ctx context.Context, vehicle domain.Dado) error {
	_, err := v.coll.InsertOne(ctx, vehicle)
	if err != nil {
		return fmt.Errorf("erro ao inserir veículo: %w", err)
	}
	return nil
}

func NewTrackerRepository() (ITrackerRepository, error) {
	db, err := GetDatabase("AthenaDB")
	if err != nil {
		return nil, fmt.Errorf("erro ao se connectar ao db: %w", err)
	}

	return &TrackerRepository{coll: db.Collection("vehicles")}, nil
}
