package modelos

type Notificacion_mercado_pago struct {
	Id             uint                           `json:"id"`
	Type           string                         `json:"type"`
	Application_id uint                           `json:"application_id"`
	User_id        uint                           `json:"user_id"`
	Action         string                         `json:"action"`
	Data           Data_notificacion_mercado_pago `json:"data"`
}

type Data_notificacion_mercado_pago struct {
	Id uint `json:"id"`
}
