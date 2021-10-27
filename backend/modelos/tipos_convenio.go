package modelos

type Tipos_convenio struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Desctconv string `json:"desctconv"`
}

func (Tipos_convenio) TableName() string {
	return "tipos_convenio"
}
