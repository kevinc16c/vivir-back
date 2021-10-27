package lugares_plc

import (
	"net/http"

	"../../database"
	Modelos "../../modelos"
	"../../utils"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros   int                   `json:"registros,omitempty"`
	Lugares_plc []Modelos.Lugares_plc `json:"palabras_clave,omitempty"`
	Lugar_plc   *Modelos.Lugares_plc  `json:"palabra_clave,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("lugares_plc.id,lugares_plc.idlugar,lugares_plc.idpalabraclave,palabras_clave.palabraclave")
	db = db.Joins("JOIN palabras_clave ON palabras_clave.id=lugares_plc.idpalabraclave")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idlugar = ? and palabraclave like ?", c.Param("id"), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idlugar = ? ", c.Param("id"))
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("palabraclave")
	}

	// Preparo paginacion
	var pagina uint = 1
	var limite uint = 10
	var offset uint = 0
	var registros int = 0
	if c.QueryParam("limite") != "" {
		limite = utils.ParseInt(c.QueryParam("limite"))
	}
	if c.QueryParam("pagina") != "" {
		pagina = utils.ParseInt(c.QueryParam("pagina"))
	}
	offset = limite * (pagina - 1)

	// Ejecuto consultas
	var palabras []Modelos.Lugares_plc
	db.Offset(offset).Limit(limite).Find(&palabras)
	db.Table("lugares_plc").Count(&registros)
	data := Data{Registros: registros, Lugares_plc: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPalabraClave(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	palabras := new(Modelos.Lugares_plc)
	db.Find(&palabras, id)

	data := Data{Lugar_plc: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	palabras := new(Modelos.Lugares_plc)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(palabras); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("palabraclave").Create(&palabras).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Lugar_plc: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM lugares_plc WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro dado de baja con éxito",
	})
}
