package modelos

type Localidades struct {
	Idlocalidad  uint       `json:"idlocalidad" gorm:"primary_key"`
	Nombrelocali string     `json:"nombrelocali"`
	Idprovincia  uint       `json:"idprovincia"`
	Provincia    Provincias `json:"provincia" gorm:"ForeignKey:idprovincia;AssociationForeignKey:idprovincia"`
}

func (Localidades) TableName() string {
	return "localidades"
}
