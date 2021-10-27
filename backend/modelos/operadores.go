package modelos

type Operadores struct {
	ID           uint               `json:"id" gorm:"primary_key"`
	Nickoperador string             `json:"nickoperador"`
	Apynombres   string             `json:"apynombres"`
	Idnivel      uint               `json:"idnivel"`
	Nivel        Operadores_niveles `json:"nivel"  gorm:"ForeignKey:idnivel;AssociationForeignKey:id"`
	Estado       string             `json:"estado"`
}

func (Operadores) TableName() string {
	return "operadores_sistema"
}
