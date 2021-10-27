package dias

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
	Registros int            `json:"registros,omitempty"`
	Dias      []Modelos.Dias `json:"dias,omitempty"`
	Dia       *Modelos.Dias  `json:"dia,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" dia LIKE ?", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var dias []Modelos.Dias
	db.Find(&dias)
	data := Data{Dias: dias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetDia(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	dias := new(Modelos.Dias)
	db = db.Where(" id=?", id)
	db.Find(&dias)

	data := Data{Dia: dias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
