package services

import (
	"athena/infrastructure/api"
	"fmt"
	"log"
	"time"
)

type ITimerService interface {
	StartApplication()
}

type TimerService struct {
	trackerService ITrackerService
	api            api.ITrackerAPIClient
}

// StartApplication implements ITimerService.
func (t *TimerService) StartApplication() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		resp, err := t.api.ListaVeiculos()
		if err != nil {
			log.Fatal("Erro ao receber resposta da aplicação: %w", err)
		}

		if resp != nil {
			fmt.Printf("Resposta recebida")
		}

		if err = t.trackerService.SaveTrackerData(); err != nil {
			log.Fatal("Erro ao salvar os dados através do serviço Tracker: %w", err)
		}
	}

}

func NewTimerService(
	trackerService ITrackerService,
	api api.ITrackerAPIClient,
) ITimerService {
	return &TimerService{
		trackerService: trackerService,
		api:            api,
	}
}
