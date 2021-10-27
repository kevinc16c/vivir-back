package modelos

type Pedidos_estados struct {
	Id     uint   `json:"id" gorm:"primary_key"`
	Estado string `json:"estado"`
}

func (Pedidos_estados) TableName() string {
	return "pedidos_estados"
}
