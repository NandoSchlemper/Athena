package services

import (
	"athena/infrastructure/api"
	"athena/repository"
	"fmt"
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

	if err = t.repository.InsertManyVehicles(r.Dados); err != nil {
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
