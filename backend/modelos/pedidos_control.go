package modelos

import "time"

type Pedidos_control struct {
	Id             uint                       `json:"id" gorm:"primary_key"`
	Fechaalta      time.Time                  `json:"fechaalta"`
	Idusuario      uint                       `json:"idusuario"`
	Emailusuario   string                     `json:"emailusuario,omitempty"`
	Usuario        Usuarios                   `json:"usuario" gorm:"ForeignKey:idusuario;AssociationForeignKey:id"`
	Idpropietario  uint                       `json:"idpropietario"`
	Propietario    Personas_humanas_juridicas `json:"propietario" gorm:"ForeignKey:idpropietario;AssociationForeignKey:id"`
	Idlugar        uint                       `json:"idlugar"`
	Lugar          Lugares                    `json:"lugar" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
	Idtiporetiro   uint                       `json:"idtiporetiro"`
	Tiporetiro     Tipos_retiro               `json:"tipo_retiro" gorm:"ForeignKey:idtiporetiro;AssociationForeignKey:id"`
	Idtipodelivery uint                       `json:"idtipodelivery"`
	Tipodelivery   Tipos_delivery             `json:"tipo_delivery" gorm:"ForeignKey:idtipodelivery;AssociationForeignKey:id"`
	Direccion      string                     `json:"direccion"`
	Idtipopago     uint                       `json:"idtipopago"`
	Importe        float64                    `json:"importe"`
	Porcomision    float64                    `jaon:"porcomision"`
	Impcomision    float64                    `json:"impcomision"`
	Impodelivery   float64                    `json:"impodelivery"`
	Impneto        float64                    `json:"impneto"`
	Idestado       uint                       `json:"idestado"`
	Estado         Pedidos_estados            `json:"estado" gorm:"ForeignKey:idestado;AssociationForeignKey:id"`
	Observaciones  string                     `json:"observaciones"`
	Detalle        []Pedidos_detalle          `json:"detalle,omitempty" gorm:"ForeignKey:idpedido;AssociationForeignKey:id"`
	Pago           Pedidos_detalle_pago       `json:"pago,omitempty" gorm:"ForeignKey:idpedido;AssociationForeignKey:id"`
}

func (Pedidos_control) TableName() string {
	return "pedidos_control"
}
