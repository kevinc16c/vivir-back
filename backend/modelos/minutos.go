package modelos

type Minutos struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Minutos     string `json:"minutos"`
	Mascaraminu string `json:"mascaraminu"`
}

func (Minutos) TableName() string {
	return "minutos"
}
