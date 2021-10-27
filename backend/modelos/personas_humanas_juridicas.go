package modelos

import "time"

type Personas_humanas_juridicas struct {
	Id           uint            `json:"id" gorm:"primary_key"`
	Razonsocial  string          `json:"razonsocial"`
	Nofantasia   string          `json:"nofantasia"`
	Idcondiva    uint            `json:"idcondiva"`
	Condicioniva Condiciones_iva `json:"iva" gorm:"ForeignKey:idcondiva;AssociationForeignKey:codigociva"`
	Numerocuit   string          `json:"numerocuit"`
	Direccion    string          `json:"direccion"`
	Idlocalidad  uint            `json:"idlocalidad"`
	Localidad    Localidades     `json:"localidad" gorm:"ForeignKey:idlocalidad;AssociationForeignKey:idlocalidad"`
	Telefono     string          `json:"telefono"`
	Telefono2    string          `json:"telefono2"`
	Celular1     string          `json:"celular1"`
	Celular2     string          `json:"celular2"`
	Celular3     string          `json:"celular3"`
	Email        string          `json:"email"`
	Estado       string          `json:"estado"`
	Fechaestado  *time.Time      `json:"fechaestado"`
	Cambiarpass  uint            `json:"cambiarpass"`
}

func (Personas_humanas_juridicas) TableName() string {
	return "personas_humanas_juridicas"
}
