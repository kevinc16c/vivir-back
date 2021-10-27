package modelos

type Usuarios_lugares_favoritos struct {
	Id           uint    `json:"id" gorm:"primary_key"`
	Idusuario    uint    `json:"idusuario"`
	Idlugar      uint    `json:"idlugar"`
	Nombrelugar  string  `json:"nombrelugar"`
	Idrubro      uint    `json:"idrubro"`
	Descrirubro  string  `json:"descrirubro"`
	Idsubrubro   uint    `json:"idsubrubro"`
	Dsubrubro    string  `json:"dsubrubro"`
	Direccion    string  `json:"direccion"`
	Idlocalidad  uint    `json:"idlocalidad"`
	Nombrelocali string  `json:"nombrelocali"`
	Idprovincia  uint    `json:"idprovincia,omitempty"`
	Nombrepcia   string  `json:"nombrepcia,omitempty"`
	Rutafoto     string  `json:"rutafoto"`
	Precdelivery float64 `json:"precdelivery"`
	Cpraminima   float64 `json:"cpraminima"`
}

func (Usuarios_lugares_favoritos) TableName() string {
	return "usuarios_lugares_favoritos"
}
