package main

import (
	"athena/infrastructure/api"
	"athena/repository"
	"athena/services"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar variaveis de ambiente.")
	}

	cfg := api.NewTrackerAPIConfig(30 * time.Second)
	cfg.SetDefaultTracker()
	client := api.NewTrackerAPIClient(cfg)
	resp, err := client.ListaVeiculos()
	if err != nil {
		log.Fatal("Erro ao receber resposta da aplicação: %w", err)
	}

	if resp != nil {
		fmt.Printf("Resposta recebida")
	}

	dbrepo := repository.NewMongoDB(1, 20, 1)
	db, err := dbrepo.InitDB()

	if err != nil {
		log.Fatal("Erro ao iniciar o banco de dados: %w", err)
	}

	tracker_repo := repository.NewTrackerRepository(db)
	tracker_service := services.NewTrackerService(tracker_repo, client)
	if err = tracker_service.SaveTrackerData(); err != nil {
		log.Fatal("Erro ao salvar os dados através do serviço Tracker: %w", err)
	}

}
