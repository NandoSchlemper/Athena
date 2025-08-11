package services

import (
	"athena/infrastructure/api"
	"athena/repository"
	"athena/utils"
	"context"
	"fmt"
	"time"
)

type ITrackerService interface {
	SaveTrackerData() error
}

type TrackerService struct {
	repository repository.ITrackerRepository
	api        api.ITrackerAPIClient
}

// SaveTrackerData implements ITrackerService.
func (t *TrackerService) SaveTrackerData() error {
	r, err := t.api.ListaVeiculos()
	if err != nil {
		return fmt.Errorf("erro na resposta da api: %w", err)
	}

	formatted := utils.ValidateSave(r)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = t.repository.InsertManyVehicles(ctx, formatted); err != nil {
		return fmt.Errorf("erro ao inserir os veículos")
	}
	return nil
}

func NewTrackerService(
	repo repository.ITrackerRepository,
	trackerapi api.ITrackerAPIClient,
) ITrackerService {
	return &TrackerService{
		repository: repo,
		api:        trackerapi,
	}
}
