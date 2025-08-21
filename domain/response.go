package domain

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
