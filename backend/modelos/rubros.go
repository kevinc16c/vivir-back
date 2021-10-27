package modelos

type Rubros struct {
	ID           uint    `json:"id" gorm:"primary_key"`
	Descrirubro  string  `json:"descrirubro"`
	Rutaimgrubro string  `json:"rutaimgrubro"`
	Porcomision  float64 `json:"porcomision"`
}

func (Rubros) TableName() string {
	return "rubros"
}
