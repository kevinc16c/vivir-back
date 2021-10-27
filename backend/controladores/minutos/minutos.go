package minutos

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
	Registros int               `json:"registros,omitempty"`
	Minutos   []Modelos.Minutos `json:"minutos,omitempty"`
	Minuto    *Modelos.Minutos  `json:"minuto,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Ejecuto consulta
	var minutos []Modelos.Minutos
	db.Find(&minutos)
	data := Data{Minutos: minutos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetMinuto(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	minutos := new(Modelos.Minutos)
	db = db.Where(" id=?", id)
	db.Find(&minutos)

	data := Data{Minuto: minutos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
