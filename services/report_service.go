package services

import (
	"athena/domain"
	"athena/repository"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type IReportService interface {
	CreateReport(ctx context.Context, filename ...string) (string, error)
}

type ReportService struct {
	databaseRepository repository.ITrackerRepository
}

func (r *ReportService) CreateReport(ctx context.Context, filename ...string) (string, error) {
	return r.createReport(ctx, r.databaseRepository.GetVehicles, filename...)
}

func (r *ReportService) createReport(ctx context.Context, dataFunc func(context.Context) ([]domain.Dado, error), filename ...string) (string, error) {
	// Gera nome do arquivo com timestamp
	outputFile := r.generateFilename(filename...)

	// Obtém dados do repositório
	res, err := dataFunc(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter dados do banco: %w", err)
	}

	if len(res) == 0 {
		return "", fmt.Errorf("nenhum dado encontrado para gerar relatório")
	}

	// Cria arquivo Excel
	f := excelize.NewFile()
	sheet := "Veículos"

	// Configura sheet
	if err := r.setupSheet(f, sheet); err != nil {
		return "", fmt.Errorf("erro ao configurar planilha: %w", err)
	}

	// Preenche dados
	if err := r.fillData(f, sheet, res); err != nil {
		return "", fmt.Errorf("erro ao preencher dados: %w", err)
	}

	// Salva arquivo
	if err := f.SaveAs(outputFile); err != nil {
		return "", fmt.Errorf("erro ao salvar arquivo Excel: %w", err)
	}

	return outputFile, nil
}

func (r *ReportService) generateFilename(filenames ...string) string {
	if len(filenames) > 0 && filenames[0] != "" {
		return filenames[0]
	}

	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("relatorio_veiculos_%s.xlsx", timestamp)
}

func (r *ReportService) setupSheet(f *excelize.File, sheet string) error {
	// Define estilos
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	// Configura cabeçalhos
	headers := map[string]string{
		"A1": "Placa",
		"B1": "Localização",
		"C1": "Horário",
		"D1": "Velocidade",
	}

	for cell, value := range headers {
		f.SetCellValue(sheet, cell, value)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 15) // Placa
	f.SetColWidth(sheet, "B", "B", 30) // Localização
	f.SetColWidth(sheet, "C", "C", 20) // Horário
	f.SetColWidth(sheet, "D", "G", 12) // Demais campos

	return nil
}

// fillData preenche os dados na planilha
func (r *ReportService) fillData(f *excelize.File, sheet string, data []domain.Dado) error {
	for i, v := range data {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.safeString(v.Placa))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.safeString(v.Localização))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.safeTime(v.Horario))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.safeInt(v.Velocidade))
	}
	return nil
}

// Funções auxiliares para tratamento seguro de dados
func (r *ReportService) safeString(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func (r *ReportService) safeTime(value time.Time) string {
	if value.IsZero() {
		return "N/A"
	}
	return value.Format("02/01/2006 15:04:05")
}

func (r *ReportService) safeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

var (
	reportInstance *ReportService
	reportOnce     sync.Once
)

func NewReportService(db repository.ITrackerRepository) IReportService {
	reportOnce.Do(func() {
		reportInstance = &ReportService{
			databaseRepository: db,
		}
	})
	return reportInstance
}
