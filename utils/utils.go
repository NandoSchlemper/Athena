package utils

import (
	"athena/domain"
	"slices"
)

func ValidateSave(r *domain.Response) []domain.Dado {
	formated_data := slices.DeleteFunc(r.Dados, func(d domain.Dado) bool {
		return d.Velocidade != 0
	})

	return formated_data
}
