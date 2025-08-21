package main

import (
	"athena/infrastructure/api"
	"athena/repository"
	"athena/services"
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

	dbrepo := repository.NewMongoDB(1, 20, 1)
	db, err := dbrepo.InitDB()

	if err != nil {
		log.Fatal("Erro ao iniciar o banco de dados. %w", err)
	}

	tracker_repo := repository.NewTrackerRepository(db)
	// reports_service := services.NewReportService(tracker_repo)
	// if err := reports_service.CreateReport(); err != nil {
	// 	log.Fatal("Erro ao criar a merda do excel %w", err)
	// }

	tracker_service := services.NewTrackerService(tracker_repo, client)
	application := services.NewTimerService(tracker_service, client)
	application.StartApplication(10) // roda a routine de 1 em 1 minuto.
}
