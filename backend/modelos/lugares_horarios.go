package modelos

type Lugares_horarios struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Idlugar    uint   `json:"idlugar"`
	Iddia      uint   `json:"iddia"`
	Dia        string `json:"dia"`
	Lughorades string `json:"Lughorades"`
	Lughorahas string `json:"Lughorahas"`
}

func (Lugares_horarios) TableName() string {
	return "lugares_horarios"
}
