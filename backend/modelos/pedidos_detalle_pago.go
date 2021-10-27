package modelos

type Pedidos_detalle_pago struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Idpedido   uint   `json:"idpedido"`
	Idpreferen string `json:"idpreferen"`
	Idpago     uint   `json:"idpago"`
	Tipo       string `json:"tipo"`
	Metodo     string `json:"metodo"`
	Estado     string `json:"estado"`
	Detalle    string `json:"detalle"`
}

func (Pedidos_detalle_pago) TableName() string {
	return "pedidos_detalle_pago"
}
