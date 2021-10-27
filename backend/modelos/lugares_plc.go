package modelos

type Lugares_plc struct {
	Id             uint   `json:"id" gorm:"primary_key"`
	Idlugar        uint   `json:"idlugar"`
	Idpalabraclave uint   `json:"idpalabraclave"`
	Palabraclave   string `json:"palabraclave"`
}

func (Lugares_plc) TableName() string {
	return "lugares_plc"
}
