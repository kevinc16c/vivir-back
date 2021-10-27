package modelos

type Operadores_niveles struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Niveloper string `json:"niveloper"`
}

func (Operadores_niveles) TableName() string {
	return "operadores_sistema_niveles"
}
