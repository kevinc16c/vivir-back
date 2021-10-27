package modelos

type Insumos struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	Idlugar     uint   `json:"idlugar"`
	Descripcion string `json:"descripcion"`
	Suspendido  uint   `json:"suspendido"`
}

func (Insumos) TableName() string {
	return "insumos"
}
