package operadores_niveles

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
	Registros          int                          `json:"registros,omitempty"`
	Operadores_niveles []Modelos.Operadores_niveles `json:"niveles,omitempty"`
	Operador_nivel     *Modelos.Operadores_niveles  `json:"nivel,omitempty"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" niveloper LIKE ?", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var operadores_niveles []Modelos.Operadores_niveles
	db.Find(&operadores_niveles)
	data := Data{Operadores_niveles: operadores_niveles}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetNivel(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	operadores_niveles := new(Modelos.Operadores_niveles)
	db.Find(&operadores_niveles, id)

	data := Data{Operador_nivel: operadores_niveles}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
