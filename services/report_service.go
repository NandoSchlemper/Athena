package services

import "athena/repository"

type IReportService interface {
	SendReportsByEmail()
}

type ReportService struct {
	databaseRepository repository.ITrackerRepository
}
