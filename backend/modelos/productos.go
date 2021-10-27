package modelos

import "time"

type Productos struct {
	Id           uint       `json:"id" gorm:"primary_key"`
	Idlugar      uint       `json:"idlugar"`
	Codintprod   string     `json:"codintprod"`
	Idcategprod  uint       `json:"idcategprod"`
	Descriprod   string     `json:"descriprod"`
	Desextprod   string     `json:"desextprod"`
	Prunitprod   float64    `json:"prunitprod"`
	Aliivaprod   float64    `json:"Aliivaprod"`
	Suspendido   uint       `json:"suspendido"`
	Agregados    uint       `json:"agregados"`
	Ocupdesde    *time.Time `json:"ocupdesde"`
	Ocuphasta    *time.Time `json:"ocuphasta"`
	Controlstock uint       `json:"controlstock"`
	Estado       string     `json:"estado"`
	Fechabaja    *time.Time `json:"fechabaja"`
}

//Lugares      Lugares              `json:"lugar" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
//Categoria    Productos_categorias `json:"categoria" gorm:"ForeignKey:idcategprod;AssociationForeignKey:id"`

func (Productos) TableName() string {
	return "productos"
}
