package services

import (
	"athena/repository"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type IReportService interface {
	CreateReport() error
}

type ReportService struct {
	databaseRepository repository.ITrackerRepository
}

// CreateReport implements IReportService.
func (r *ReportService) CreateReport() error {
	res, err := r.databaseRepository.GetVehicles()
	sheet := "Sheet1"

	if err != nil {
		return fmt.Errorf("erro ao criar report baseado no banco: %w", err)
	}

	f := excelize.NewFile()

	f.SetCellValue(sheet, "A1", "Placa")
	f.SetCellValue(sheet, "B1", "Localização")

	for i, v := range res {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.Placa)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Localização)
	}

	if err := f.SaveAs("output.xlsx"); err != nil {
		return fmt.Errorf("erro ao salvar arquivo excel: %w", err)
	}

	return nil
}

func NewReportService(db repository.ITrackerRepository) IReportService {
	return &ReportService{
		databaseRepository: db,
	}
}
