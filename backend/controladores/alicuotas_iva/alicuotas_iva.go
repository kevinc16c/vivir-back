package alicuotas_iva

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
	Registros int                     `json:"registros,omitempty"`
	Alicuotas []Modelos.Alicuotas_iva `json:"alicuotas_iva,omitempty"`
	Alicuota  *Modelos.Alicuotas_iva  `json:"alicuota_iva,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("codigoalic, descrialic, alicuoalic")
	db = db.Where(" suspenfact=0")

	// Ejecuto consulta
	var alicuotas []Modelos.Alicuotas_iva
	db.Find(&alicuotas)
	data := Data{Alicuotas: alicuotas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetAlicuota(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	alicuotas := new(Modelos.Alicuotas_iva)
	db = db.Select("codigoalic, descrialic, alicuoalic")
	db = db.Where(" codigoalic=?", id)
	db.Find(&alicuotas)

	data := Data{Alicuota: alicuotas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
