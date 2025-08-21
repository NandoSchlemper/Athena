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
	InsertOneVehicle(vehicle domain.Dado) error
	GetVehicles() ([]domain.Dado, error)
}

type TrackerRepository struct {
	coll *mongo.Collection
}

// GetVehicles implements ITrackerRepository.
func (v TrackerRepository) GetVehicles() ([]domain.Dado, error) {
	findOpt := bson.D{{}} // todos os dados da collection
	var result []domain.Dado

	// obs: Mongo retornar um cursor, não os objetos diretamente, isso ajuda a evitar carregar
	// tudo na memoria de uma vez, bem interessante xD
	cur, err := v.coll.Find(context.Background(), findOpt)
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar todos os dados da collection Vehicles: %w", err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var dado domain.Dado
		if err := cur.Decode(&dado); err != nil {
			return nil, fmt.Errorf("erro ao iterar sobre o cursor: %w", err)
		}
		result = append(result, dado)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("erro ao carregar o cursor: %w", err)
	}

	return result, nil
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
func (v TrackerRepository) InsertOneVehicle(vehicle domain.Dado) error {
	_, err := v.coll.InsertOne(context.TODO(), vehicle)
	if err != nil {
		return fmt.Errorf("erro ao inserir uma posição do veículo no DB: %w", err)
	}
	return nil
}

func NewTrackerRepository(db *mongo.Database) ITrackerRepository {
	return &TrackerRepository{coll: db.Collection("vehicles")}
}
