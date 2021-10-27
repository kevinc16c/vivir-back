package modelos

type Usuarios struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Apellido    string `json:"apellido"`
	Nombres     string `json:"nombres"`
	Numedocume  string `json:"numedocume"`
	Celular     string `json:"celular"`
	Correoelec  string `json:"correoelec"`
	Contrasena  string `json:"contrasena,omitempty"`
	Rutaimgusu  string `json:"rutaimgusu"`
	Asignarpass uint   `json:"asignarpass"`
}

func (Usuarios) TableName() string {
	return "usuarios_app"
}
