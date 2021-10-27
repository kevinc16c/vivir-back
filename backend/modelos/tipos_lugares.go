package modelos

type Tipos_lugares struct {
	Id        uint   `json:"id" gorm:"primary_key"`
	Tipolugar string `json:"tipolugar"`
}

func (Tipos_lugares) TableName() string {
	return "tipos_lugares"
}
