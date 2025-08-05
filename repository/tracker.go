package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type IVehicleRepo interface{}

// Inserir veiculos baseados no retorno da API

type VehicleRepo struct {
	coll *mongo.Collection
}

func NewVehicleRepo(db *mongo.Database) IVehicleRepo {
	return &VehicleRepo{coll: db.Collection("vehicles")}
}
