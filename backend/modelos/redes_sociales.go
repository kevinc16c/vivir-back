package modelos

import "time"

type Redes_sociales struct {
	Id         uint      `json:"id" gorm:"primary_key"`
	Nombrerrss string    `json:"nombrerrss"`
	Rutaimgapp string    `json:"rutaimgapp"`
	Rutaimgweb string    `json:"rutaimgweb"`
	Fechamodif time.Time `json:"fechamodif"`
	Estado     string    `json:"estado"`
}

func (Redes_sociales) TableName() string {
	return "redes_sociales"
}
