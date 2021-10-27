package tipos_convenio

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
	TiposConvenio []Modelos.Tipos_convenio `json:"tiposConvenio,omitempty"`
	TipoConvenio  *Modelos.Tipos_convenio  `json:"tipoConvenio,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" desctconv LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var convenios []Modelos.Tipos_convenio
	db.Find(&convenios)
	data := Data{TiposConvenio: convenios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetTipoConvenio(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	convenios := new(Modelos.Tipos_convenio)
	db.Find(&convenios, id)

	data := Data{TipoConvenio: convenios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
