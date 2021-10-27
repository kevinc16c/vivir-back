package modelos

type Tipos_retiro struct {
	Id   uint   `json:"id" gorm:"primary_key"`
	Tipo string `json:"tipo"`
}

func (Tipos_retiro) TableName() string {
	return "tipos_retiro"
}
