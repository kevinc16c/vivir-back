package tipos_delivery

import (
	"net/http"

	"../../database"
	Modelos "../../modelos"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros     int                      `json:"registros,omitempty"`
	TiposDelivery []Modelos.Tipos_delivery `json:"tiposDelivery,omitempty"`
	TipoDelivery  *Modelos.Tipos_delivery  `json:"tipoDelivery,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" tipodelivery LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var tipos []Modelos.Tipos_delivery
	db.Find(&tipos)
	data := Data{TiposDelivery: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetTipoDelivery(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	tipos := new(Modelos.Tipos_delivery)
	db.Find(&tipos, id)

	data := Data{TipoDelivery: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
