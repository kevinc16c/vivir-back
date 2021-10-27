package condiciones_iva

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
	Registros      int                       `json:"registros,omitempty"`
	Condicionesiva []Modelos.Condiciones_iva `json:"condiciones_iva,omitempty"`
	Condicioniva   *Modelos.Condiciones_iva  `json:"condicion_iva,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" descriciva LIKE ? AND suspendido=0", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" suspendido=0")
	}

	db = db.Order("descriciva")

	// Ejecuto consulta
	var condiciones []Modelos.Condiciones_iva
	db.Find(&condiciones)
	data := Data{Condicionesiva: condiciones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetCondicion(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	condiciones := new(Modelos.Condiciones_iva)
	db = db.Where(" codigociva=?", id)
	db.Find(&condiciones)

	data := Data{Condicioniva: condiciones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
