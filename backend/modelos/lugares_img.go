package modelos

type Lugares_img struct {
	Id        uint   `json:"id" gorm:"primary_key"`
	Idlugar   uint   `json:"idlugar"`
	Rutaimg   string `json:"rutaimg"`
	Tituloimg string `json:"tituloimg"`
	Descriimg string `json:"descriimg"`
}

func (Lugares_img) TableName() string {
	return "lugares_img"
}
