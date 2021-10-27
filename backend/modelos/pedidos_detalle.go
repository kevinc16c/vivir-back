package modelos

import "time"

type Pedidos_detalle struct {
	Id           uint       `json:"id" gorm:"primary_key"`
	Idpedido     uint       `json:"idpedido"`
	Idproducto   uint       `json:"idproducto"`
	Producto     Productos  `json:"producto" gorm:"ForeignKey:idproducto;AssociationForeignKey:id"`
	Cantidad     uint       `json:"cantidad"`
	Punitario    float64    `json:"punitario"`
	Subtotal     float64    `json:"subtotal"`
	Desde        *time.Time `json:"desde"`
	Hasta        *time.Time `json:"hasta"`
	Descvariedad string     `json:"descvariedad"`
	Notas        string     `json:"notas"`
}

func (Pedidos_detalle) TableName() string {
	return "pedidos_detalle"
}
