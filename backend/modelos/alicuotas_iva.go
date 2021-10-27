package modelos

type Alicuotas_iva struct {
	Codigoalic uint    `json:"codigoalic" gorm:"primary_key"`
	Descrialic string  `json:"descrialic"`
	Alicuoalic float64 `json:"alicuoalic"`
}

func (Alicuotas_iva) TableName() string {
	return "alicuotas_iva"
}
