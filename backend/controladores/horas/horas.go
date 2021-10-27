package horas

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
	Registros int             `json:"registros,omitempty"`
	Horas     []Modelos.Horas `json:"horas,omitempty"`
	Hora      *Modelos.Horas  `json:"hora,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Ejecuto consulta
	var horas []Modelos.Horas
	db.Find(&horas)
	data := Data{Horas: horas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetHora(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	horas := new(Modelos.Horas)
	db = db.Where(" id=?", id)
	db.Find(&horas)

	data := Data{Hora: horas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
