package modelos

import "time"

type Registro_usuarios_lugares struct {
	Id           uint      `json:"id" gorm:"primary_key"`
	Idusuario    uint      `json:"idusuario"`
	Usuario      Usuarios  `json:"usuario" gorm:"ForeignKey:idusuario;AssociationForeignKey:id"`
	Idlugar      uint      `json:"idlugar"`
	Nombrelugar  string    `json:"nombrelugar"`
	Idrubro      uint      `json:"idrubro"`
	Descrirubro  string    `json:"descrirubro"`
	Idsubrubro   uint      `json:"idsubrubro"`
	Dsubrubro    string    `json:"dsubrubro"`
	Direccion    string    `json:"direccion"`
	Telefono     string    `json:"telefono"`
	Celular      string    `json:"celular"`
	E_mail       string    `json:"e_mail"`
	Idlocalidad  uint      `json:"idlocalidad"`
	Nombrelocali string    `json:"nombrelocali"`
	Fechayhora   time.Time `json:"fechayhora"`
}

func (Registro_usuarios_lugares) TableName() string {
	return "registros_usuarios_lugares"
}
