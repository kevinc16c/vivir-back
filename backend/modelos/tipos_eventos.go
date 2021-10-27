package modelos

type Tipos_eventos struct {
	Id   uint   `json:"id" gorm:"primary_key"`
	Tipo string `json:"tipo"`
}

func (Tipos_eventos) TableName() string {
	return "tipos_eventos"
}
