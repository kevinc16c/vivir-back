package modelos

type Paises struct {
	Idpais     uint   `json:"idpais" gorm:"primary_key"`
	Nombrepais string `json:"nombrepais"`
}

func (Paises) TableName() string {
	return "paises"
}
