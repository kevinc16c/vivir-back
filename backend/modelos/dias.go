package modelos

type Dias struct {
	Id  uint   `json:"id" gorm:"primary_key"`
	Dia string `json:"dia"`
}

func (Dias) TableName() string {
	return "dias"
}
