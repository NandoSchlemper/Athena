package repository

import (
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IMongoDB interface {
	InitDB() (*mongo.Database, error)
}

type MongoDB struct {
	Uri             string
	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnIdleTime uint64
}

// InitDB implements IMongoDB.
func (m *MongoDB) InitDB() (*mongo.Database, error) {
	m.Uri = os.Getenv("DB_URL")

	client, err := mongo.Connect(options.Client().ApplyURI(m.Uri).SetMaxPoolSize(m.MaxPoolSize).SetMinPoolSize(m.MinPoolSize).SetMaxConnIdleTime(time.Duration(m.MaxConnIdleTime)))

	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar o db, bruh")
	}

	DB := client.Database("AthenaDB")
	return DB, nil
}

func NewMongoDB(minpool, maxpool, maxconn uint64) IMongoDB {
	return &MongoDB{
		MinPoolSize:     minpool,
		MaxPoolSize:     maxpool,
		MaxConnIdleTime: maxconn,
	}
}
