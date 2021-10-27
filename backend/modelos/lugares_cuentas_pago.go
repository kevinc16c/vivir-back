package modelos

import "time"

type Lugares_cuentas_pago struct {
	Id           uint      `json:"id" gorm:"primary_key"`
	Idlugar      uint      `json:"idlugar"`
	Userid       string    `json:"userid"`
	Publickey    string    `json:"publickey,omitempty"`
	Accesstoken  string    `json:"accesstoken,omitempty"`
	Refreshtoken string    `json:"refreshtoken,omitempty"`
	Vencimiento  time.Time `json:"vencimiento"`
	Suspendido   uint      `json:"suspendido"`
}

func (Lugares_cuentas_pago) TableName() string {
	return "lugares_cuentas_pago"
}
