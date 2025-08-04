package main

import (
	"athena/infrastructure/api"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar variaveis de ambiente.")
	}

	cfg := api.NewWrsatAPIConfig(30 * time.Second)
	cfg.SetAPIVariables()
	client := api.NewWrsatAPIClient(cfg)
	resp, err := client.ListaVeiculos()
	if err != nil {
		log.Fatal("Erro ao receber resposta da aplicação")
	}

	if resp != nil {
		fmt.Printf("Resposta recebida")
	}
}
