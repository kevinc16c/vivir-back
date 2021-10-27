package modelos

type Palabras_clave struct {
	Id           uint   `json:"id" gorm:"primary_key"`
	Palabraclave string `json:"palabraclave"`
	Idrubro      uint   `json:"idrubro"`
	Descrirubro  string `json:"descrirubro"`
	Estado       string `json:"estado"`
}

func (Palabras_clave) TableName() string {
	return "palabras_clave"
}
