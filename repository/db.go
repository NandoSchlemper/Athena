package repository

import (
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IMongoDB interface {
	InitDB() *mongo.Database
}

type MongoDB struct {
	Uri             string
	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnIdleTime uint64
}

// InitDB implements IMongoDB.
func (m *MongoDB) InitDB() *mongo.Database {
	m.Uri = os.Getenv("DB_URL")

	client, _ := mongo.Connect(options.Client().ApplyURI(m.Uri).SetMaxPoolSize(m.MaxPoolSize).SetMinPoolSize(m.MinPoolSize).SetMaxConnIdleTime(time.Duration(m.MaxConnIdleTime)))

	DB := client.Database("AthenaDB")
	return DB
}

func NewMongoDB(minpool, maxpool, maxconn uint64) IMongoDB {
	return &MongoDB{
		MinPoolSize:     minpool,
		MaxPoolSize:     maxpool,
		MaxConnIdleTime: maxconn,
	}
}
