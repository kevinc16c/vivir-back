package tipos_eventos

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
	Registros    int                     `json:"registros,omitempty"`
	TiposEventos []Modelos.Tipos_eventos `json:"tiposEventos,omitempty"`
	TipoEvento   *Modelos.Tipos_eventos  `json:"tipoEvento,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" tipo LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var tipos []Modelos.Tipos_eventos
	db.Find(&tipos)
	data := Data{TiposEventos: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetTipoEvento(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	tipos := new(Modelos.Tipos_eventos)
	db.Find(&tipos, id)

	data := Data{TipoEvento: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
