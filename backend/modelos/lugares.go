package modelos

import "time"

type Lugares struct {
	Idlugar                    uint                       `json:"idlugar" gorm:"primary_key"`
	Idpropietario              uint                       `json:"idpropietario"`
	Personas_humanas_juridicas Personas_humanas_juridicas `json:"propietario" gorm:"ForeignKey:idpropietario;AssociationForeignKey:id"`
	Idrubro                    uint                       `json:"idrubro"`
	Idsubrubro                 uint                       `json:"idsubrubro"`
	Idtipolugar                uint                       `json:"idtipolugar"`
	Idtipoconv                 uint                       `json:"idtipoconv"`
	Direccion                  string                     `json:"direccion"`
	Idlocalidad                uint                       `json:"idlocalidad"`
	Localidad                  Localidades                `json:"localidad" gorm:"ForeignKey:idlocalidad;AssociationForeignKey:idlocalidad"`
	Latitud                    float64                    `json:"latitud"`
	Longitud                   float64                    `json:"longitud"`
	Nombrelugar                string                     `json:"nombrelugar"`
	Telefono                   string                     `json:"telefono"`
	Celular                    string                     `json:"celular"`
	E_mail                     string                     `json:"e_mail"`
	Sitioweb                   string                     `json:"sitioweb"`
	Describreve                string                     `json:"describreve"`
	Idtipodelivery             uint                       `json:"idtipodelivery"`
	Precdelivery               float64                    `json:"precdelivery"`
	Conpddiferido              uint                       `json:"conpddiferido"`
	Cpraminima                 float64                    `json:"cpraminima"`
	Impoabono                  float64                    `json:"impoabono"`
	Porcomision                float64                    `json:"porcomision"`
	Activo                     uint                       `json:"activo"`
	Fechaalta                  time.Time                  `json:"fechaalta"`
	Vencimiento                time.Time                  `json:"vencimiento"`
	Fechamodif                 *time.Time                 `json:"fechamodif"`
	Fechaestado                *time.Time                 `json:"fechaestado"`
	Estado                     string                     `json:"estado"`
	Horarios                   []Lugares_horarios         `json:"horarios" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
	Redes                      []Lugares_rrss             `json:"redes" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
	Imagenes                   []Lugares_img              `json:"imagenes" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
}

func (Lugares) TableName() string {
	return "lugares"
}
