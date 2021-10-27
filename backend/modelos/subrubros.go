package modelos

type Subrubros struct {
	Id           uint    `json:"id" gorm:"primary_key"`
	Dsubrubro    string  `json:"dsubrubro"`
	Idrubro      uint    `json:"idrubro"`
	Rubro        Rubros  `json:"rubro" gorm:"ForeignKey:idrubro;AssociationForeignKey:id"`
	Porcomision  float64 `json:"porcomision"`
	Rimgsubrubro string  `json:"rimgsubrubro"`
}

func (Subrubros) TableName() string {
	return "subrubros"
}
