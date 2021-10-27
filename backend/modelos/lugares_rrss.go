package modelos

type Lugares_rrss struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Idlugar     uint   `json:"idlugar"`
	Idrrss      uint   `json:"idrrss"`
	Nombrerrss  string `json:"nombrerrss"`
	Descriprrss string `json:"descriprrss"`
	Urlrrss     string `json:"urlrrss"`
}

func (Lugares_rrss) TableName() string {
	return "lugares_rrss"
}
