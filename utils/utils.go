package utils

import (
	"athena/domain"
	"slices"
	"time"
)

// Retorna um slice de todos os Dados que tiverem sua Velocidade = 0
func ValidateSave(r *domain.Response) []domain.Dado {
	formated_data := slices.DeleteFunc(r.Dados, func(d domain.Dado) bool {
		return d.Velocidade != 0
	})

	// adicionando o horario basicamente, kekw
	for i := range formated_data {
		formated_data[i].Horario = time.Now()
	}

	return formated_data
}
