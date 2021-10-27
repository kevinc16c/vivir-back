package modelos

type Tipos_delivery struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Tipodelivery string `json:"tipodelivery"`
}

func (Tipos_delivery) TableName() string {
	return "tipos_delivery"
}
