package modelos

type Productos_img struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Idproducto  uint   `json:"idproducto"`
	Rutaimgprod string `json:"rutaimgprod"`
}

func (Productos_img) TableName() string {
	return "productos_img"
}
