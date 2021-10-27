package modelos

type Provincias struct {
	Idprovincia uint   `json:"idprovincia" gorm:"primary_key"`
	Nombrepcia  string `json:"nombrepcia"`
	Idpais      uint   `json:"idpais"`
	Pais        Paises `json:"pais" gorm:"ForeignKey:idpais;AssociationForeignKey:idpais"`
}

func (Provincias) TableName() string {
	return "provincias"
}
