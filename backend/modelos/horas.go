package modelos

type Horas struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Hora        string `json:"hora"`
	Mascarahora string `json:"mascarahora"`
}

func (Horas) TableName() string {
	return "horas"
}
