package modelos

import "time"

type Eventos struct {
	Id           uint           `json:"id" gorm:"primary_key"`
	Idtipoeven   uint           `json:"idtipoeven"`
	Tipo         Tipos_eventos  `json:"tipo" gorm:"ForeignKey:idtipoeven;AssociationForeignKey:id"`
	Idlocalidad  uint           `json:"idlocalidad"`
	Localidad    Localidades    `json:"localidad" gorm:"ForeignKey:idlocalidad;AssociationForeignKey:idlocalidad"`
	Latitud      float64        `json:"latitud"`
	Longitud     float64        `json:"longitud"`
	Evento       string         `json:"evento"`
	Descripcion  string         `json:"descripcion"`
	Direccion    string         `json:"direccion"`
	Telefono     string         `json:"telefono"`
	Celular      string         `json:"celular"`
	E_mail       string         `json:"e_mail"`
	Sitioweb     string         `json:"sitioweb"`
	Fechaalta    time.Time      `json:"fechaalta,omitempty"`
	Inicioevento time.Time      `json:"inicioevento"`
	Finalevento  time.Time      `json:"finalevento"`
	Fechamodif   *time.Time     `json:"fechamodif,omitempty"`
	Rutafoto     string         `json:"rutafoto"`
	Distancia    float64        `json:"distancia,omitempty"`
	Redes        []Eventos_rrss `json:"redes,omitempty" gorm:"ForeignKey:idevento;AssociationForeignKey:id"`
	Imagenes     []Eventos_img  `json:"imagenes,omitempty" gorm:"ForeignKey:idevento;AssociationForeignKey:id"`
}

func (Eventos) TableName() string {
	return "eventos"
}
