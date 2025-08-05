package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Vehicle struct {
	ID    bson.ObjectID `bson:"_id"`
	Plate string        `bson:"plate"`
}
