package modelos

type Condiciones_iva struct {
	Codigociva uint   `json:"codigociva" gorm:"primary_key"`
	Descriciva string `json:"descriciva"`
}

func (Condiciones_iva) TableName() string {
	return "condicion_iva"
}
