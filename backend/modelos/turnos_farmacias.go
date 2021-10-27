package modelos

import "time"

type Turnos_farmacias struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idlugar     uint      `json:"idlugar"`
	Nombrelugar string    `json:"nombrelugar,omitempty"`
	Direccion   string    `json:"direccion,omitempty"`
	Inicioturno time.Time `json:"inicioturno"`
	Finalturno  time.Time `json:"finalturno"`
}

func (Turnos_farmacias) TableName() string {
	return "turnos_farmacias"
}
