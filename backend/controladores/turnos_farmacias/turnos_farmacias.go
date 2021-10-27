package turnos_farmacias

import (
	"net/http"
	"time"

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
	Registros int                        `json:"registros,omitempty"`
	Turnos    []Modelos.Turnos_farmacias `json:"turnos,omitempty"`
	Turno     *Modelos.Turnos_farmacias  `json:"turno,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()
	fecha := time.Now() //utils.GetNow()
	hoy := fecha.Format("2006-01-02")

	db = db.Select("turnos_farmacias.id,lugares.nombrelugar,lugares.direccion,turnos_farmacias.inicioturno,turnos_farmacias.finalturno")
	db = db.Joins("JOIN lugares ON lugares.idlugar=turnos_farmacias.idlugar")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" DATE(finalturno)>=? AND (turnos_farmacias.id like ? OR nombrelugar like ? OR direccion like ?)", hoy, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" DATE(finalturno)>=?", hoy)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("nombrelugar")
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
	var turnos []Modelos.Turnos_farmacias
	db.Offset(offset).Limit(limite).Find(&turnos)
	db.Table("turnos_farmacias").Count(&registros)
	data := Data{Registros: registros, Turnos: turnos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetTurno(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	db = db.Select("turnos_farmacias.id,turnos_farmacias.idlugar,lugares.nombrelugar,lugares.direccion,turnos_farmacias.inicioturno,turnos_farmacias.finalturno")
	db = db.Joins("JOIN lugares ON lugares.idlugar=turnos_farmacias.idlugar")

	turnos := new(Modelos.Turnos_farmacias)
	db.Find(&turnos, id)

	data := Data{Turno: turnos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	turnos := new(Modelos.Turnos_farmacias)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(turnos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("nombrelugar", "direccion").Create(&turnos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Turno: turnos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM turnos_farmacias WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Turno dado de baja con éxito",
	})
}
