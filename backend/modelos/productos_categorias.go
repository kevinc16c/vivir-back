package modelos

type Productos_categorias struct {
	Id            uint   `json:"id" gorm:"primary_key"`
	Descricatprod string `json:"descricatprod"`
	Idrubro       uint   `json:"idrubro"`
	Rubro         Rubros `json:"rubro" gorm:"ForeignKey:idrubro;AssociationForeignKey:id"`
	Estado        string `json:"estado"`
}

func (Productos_categorias) TableName() string {
	return "productos_categorias"
}
