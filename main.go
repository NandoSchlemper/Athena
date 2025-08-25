package main

import (
	"athena/infrastructure/api"
	"athena/repository"
	"athena/services"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erro ao carregar variáveis de ambiente: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	cfg := api.NewTrackerAPIConfig(30 * time.Second)
	cfg.SetDefaultTracker()
	client := api.NewTrackerAPIClient(cfg)

	trackerRepo, err := repository.NewTrackerRepository()
	if err != nil {
		log.Fatalf("Erro ao criar repositório: %v", err)
	}

	trackerService := services.NewTrackerService(trackerRepo, client)
	timerService := services.NewTimerService(trackerService, client)

	log.Println("Iniciando Athena Application...")
	timerService.StartApplication(ctx, 1)

	<-ctx.Done() // aguarda o evento de shutdown
	log.Println("Recebido sinal de desligamento...")

	timerService.StopApplication()

	if err := repository.CloseMongoConnection(); err != nil {
		log.Printf("Erro ao fechar conexão MongoDB: %v", err)
	} // basicamente da o close na conexão

	log.Println("Aplicação finalizada corretamente")
}
