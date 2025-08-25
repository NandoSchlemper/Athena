package utils

import (
	"athena/domain"
	"slices"
)

// Retorna um slice de todos os Dados que tiverem sua Velocidade = 0
func ValidateSave(r *domain.Response) []domain.Dado {
	formated_data := slices.DeleteFunc(r.Dados, func(d domain.Dado) bool {
		return d.Velocidade != 0
	})

	return formated_data
}
