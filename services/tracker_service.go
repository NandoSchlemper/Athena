package services

import "athena/repository"

type ITrackerService interface{}

type TrackerService struct {
	repository repository.ITrackerRepository
}

func NewTrackerService(repo repository.ITrackerRepository) ITrackerService {
	return &TrackerService{
		repository: repo,
	}
}
