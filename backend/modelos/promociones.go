package modelos

import "time"

type Promociones struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idlugar     uint      `json:"idlugar"`
	Nombrelugar string    `json:"nombrelugar"`
	Direccion   string    `json:"direccion"`
	Latitud     float64   `json:"latitud"`
	Longitud    float64   `json:"longitud"`
	Idrubro     uint      `json:"idrubro"`
	Descrirubro string    `json:"descrirubro"`
	Idsubrubro  uint      `json:"idsubrubro"`
	Dsubrubro   string    `json:"dsubrubro"`
	Vencimiento time.Time `json:"vencimiento"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Terminos    string    `json:"terminos"`
	Canticupos  uint      `json:"canticupos"`
	Cuposdispon uint      `json:"cuposdispon"`
	Rutaimg     string    `json:"rutaimg"`
	Idtipopromo uint      `json:"idtipopromo"`
	Estado      string    `json:"estado"`
}

func (Promociones) TableName() string {
	return "promociones"
}
