package modelos

import "time"

type Usuarios_sesiones struct {
	Id         uint      `json:"id" gorm:"primary_key"`
	Idusuario  uint      `json:"idusuario"`
	Token      string    `json:"token"`
	Fechaalta  time.Time `json:"fechaalta"`
	Fechamodif time.Time `json:"fechamodif"`
}

func (Usuarios_sesiones) TableName() string {
	return "usuarios_app_sesiones"
}
