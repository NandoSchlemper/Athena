package domain

type Operacao struct {
	ID   string
	Nome string
}

type Dado struct {
	ID         string   `bson:"id" json:"id"`
	Placa      string   `bson:"placa" json:"placa"`
	Lat        float64  `bson:"lat,string" json:"lat,string"`
	Lon        float64  `bson:"lon,string" json:"lon,string"`
	Velocidade int      `bson:"velocidade,string" json:"velocidade,string"`
	Operacao   Operacao `bson:"web_grupo" json:"web_grupo"`
	Estado     string   `bson:"estado"`
	Motorista  string   `bson:"motorista"`
}

type Response struct {
	Erro      bool   `bson:"erro" json:"erro"`
	Status    int    `bson:"status" json:"status"`
	Mensagem  string `bson:"mensagem" json:"mensagem"`
	Ordem     string `bson:"ordem" json:"ordem"`
	Limit     string `bson:"limit" json:"limit"`
	Pagina    string `bson:"pagina" json:"pagina"`
	QtdResult int    `bson:"qtd_result" json:"qtd_result"`
	Dados     []Dado `bson:"dados" json:"dados"`
}
