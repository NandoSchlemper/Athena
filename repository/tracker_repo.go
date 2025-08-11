package repository

import (
	"athena/domain"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ITrackerRepository interface {
	InsertManyVehicles(vehicles []domain.Dado) error
	InsertOneVehicle(vehicle domain.Dado) error
}

type TrackerRepository struct {
	coll *mongo.Collection
}

// InsertManyVehicles implements IVehicleRepository.
func (v *TrackerRepository) InsertManyVehicles(vehicles []domain.Dado) error {
	_, err := v.coll.InsertMany(context.TODO(), vehicles)
	if err != nil {
		return fmt.Errorf("erro ao inserir varias posições de veículos no DB: %w", err)
	}
	return nil
}

// InsertOneVehicle implements IVehicleRepository.
func (v *TrackerRepository) InsertOneVehicle(vehicle domain.Dado) error {
	_, err := v.coll.InsertOne(context.TODO(), vehicle)
	if err != nil {
		return fmt.Errorf("erro ao inserir uma posição do veículo no DB: %w", err)
	}
	return nil
}

func NewTrackerRepository(db *mongo.Database) ITrackerRepository {
	return &TrackerRepository{coll: db.Collection("vehicles")}
}
