package modelos

type Eventos_img struct {
	Id        uint   `json:"id" gorm:"primary_key"`
	Idevento  uint   `json:"idevento"`
	Rutaimg   string `json:"rutaimg"`
	Tituloimg string `json:"tituloimg"`
	Descriimg string `json:"descriimg"`
}

func (Eventos_img) TableName() string {
	return "eventos_img"
}
