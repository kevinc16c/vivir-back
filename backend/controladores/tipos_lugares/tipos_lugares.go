package tipos_lugares

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
	TiposLugares []Modelos.Tipos_lugares `json:"tiposLugares,omitempty"`
	TipoLugar    *Modelos.Tipos_lugares  `json:"tipoLugar,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" tipolugar LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var tipos []Modelos.Tipos_lugares
	db.Find(&tipos)
	data := Data{TiposLugares: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetTipoLugar(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	tipos := new(Modelos.Tipos_lugares)
	db.Find(&tipos, id)

	data := Data{TipoLugar: tipos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
