package modelos

import "time"

type Valoraciones struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idlugar     uint      `json:"idlugar"`
	Nombrelugar string    `json:"nombrelugar"`
	Idusuario   uint      `json:"idusuario"`
	Apellido    string    `json:"apellido"`
	Nombres     string    `json:"nombres"`
	Rutaimgusu  string    `json:"rutaimgusu"`
	Puntuacion  float64   `json:"puntuacion"`
	Fecha       time.Time `json:"fecha"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Fechamodif  time.Time `json:"fechamodif,omitempty"`
}

func (Valoraciones) TableName() string {
	return "valoraciones"
}
