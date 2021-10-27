package modelos

type Eventos_rrss struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Idevento    uint   `json:"idevento"`
	Idrrss      uint   `json:"idrrss"`
	Nombrerrss  string `json:"nombrerrss"`
	Descriprrss string `json:"descriprrss"`
	Rutaimgapp  string `json:"rutaimgapp"`
	Urlrrss     string `json:"urlrrss"`
}

func (Eventos_rrss) TableName() string {
	return "eventos_rrss"
}
