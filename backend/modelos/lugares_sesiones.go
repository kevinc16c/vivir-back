package modelos

import "time"

type Lugares_sesiones struct {
	Id         uint      `json:"id" gorm:"primary_key"`
	Idlugar    uint      `json:"idlugar"`
	Token      string    `json:"token"`
	Fechaalta  time.Time `json:"fechaalta"`
	Fechamodif time.Time `json:"fechamodif"`
}

func (Lugares_sesiones) TableName() string {
	return "lugares_sesiones"
}
