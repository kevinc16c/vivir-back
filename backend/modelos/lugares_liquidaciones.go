package modelos

import "time"

type Lugares_liquidaciones struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idlugar     uint      `json:"idlugar"`
	Lugar       Lugares   `json:"lugar" gorm:"ForeignKey:idlugar;AssociationForeignKey:idlugar"`
	Fecha       time.Time `json:"fecha"`
	Vencimiento time.Time `json:"vencimiento"`
	Observacion string    `json:"observacion"`
	Totalbruto  float64   `json:"totalbruto"`
	Comisefvo   float64   `json:"comisefvo"`
	Comismpgo   float64   `json:"comismpgo"`
	Acobrar     float64   `json:"acobrar"`
}

func (Lugares_liquidaciones) TableName() string {
	return "lugares_liquidaciones"
}
