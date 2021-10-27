package modelos

type Productos_insumos struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Idproducto uint   `json:"idproducto"`
	Idinsumo   uint   `json:"idinsumo"`
	Dinsuprodu string `json:"dinsuprodu"`
}

func (Productos_insumos) TableName() string {
	return "productos_insumos"
}
