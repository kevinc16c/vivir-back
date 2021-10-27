package modelos

import "time"

type Promociones_canje struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idpromocion uint      `json:"idpromocion"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Imagen      string    `json:"imagen"`
	Fecha       time.Time `json:"fecha"`
}

func (Promociones_canje) TableName() string {
	return "promociones_canje"
}
