package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	once      sync.Once
	client    *mongo.Client
	clientErr error
)

type MongoDBConfig struct {
	URI             string
	MinPoolSize     uint64
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

func GetMongoClient(config MongoDBConfig) (*mongo.Client, error) {
	once.Do(func() {
		uri := config.URI
		if uri == "" {
			uri = os.Getenv("DB_URL")
		}

		if uri == "" {
			clientErr = fmt.Errorf("DB_URL não configurado")
			return
		}

		opts := options.Client().
			ApplyURI(uri).
			SetMinPoolSize(config.MinPoolSize).
			SetMaxPoolSize(config.MaxPoolSize).
			SetMaxConnIdleTime(config.MaxConnIdleTime)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, clientErr = mongo.Connect(opts)
		if clientErr != nil {
			return
		}

		// Verifica a conexão
		clientErr = client.Ping(ctx, nil)
	})

	return client, clientErr
}

func GetDatabase(dbName string) (*mongo.Database, error) {
	config := MongoDBConfig{
		URI:             os.Getenv("DB_URL"),
		MinPoolSize:     10,
		MaxPoolSize:     100,
		MaxConnIdleTime: 30 * time.Minute,
	}

	client, err := GetMongoClient(config)
	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}

// Função para fechar a conexão (usar no shutdown da aplicação)
func CloseMongoConnection() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}
