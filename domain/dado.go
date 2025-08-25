package domain

import "time"

type Dado struct {
	ID    string `bson:"id" json:"id"`
	Placa string `bson:"placa" json:"placa"`
	// Lat         float64 `bson:"lat,string" json:"lat,string" no:"3"`
	// Lon         float64 `bson:"lon,string" json:"lon,string" no:"4"`
	Localização string `bson:"localizacao" json:"localizacao"`
	Velocidade  int    `bson:"velocidade,string" json:"velocidade,string" no:"6"`
	// Operacao    Operacao `bson:"web_grupo" json:"web_grupo"`
	// Estado    string `bson:"estado" no:"7"`
	// Motorista string `bson:"motorista" no:"8"`
	Horario time.Time `bson:"data" json:"datagps"`
}
