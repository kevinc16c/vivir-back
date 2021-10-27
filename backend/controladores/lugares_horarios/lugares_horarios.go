package lugares_horarios

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
	Registros int                        `json:"registros,omitempty"`
	Horarios  []Modelos.Lugares_horarios `json:"horarios,omitempty"`
	Horario   *Modelos.Lugares_horarios  `json:"horario,omitempty"`
}

func GetHorarios(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas")
	db = db.Joins("JOIN dias ON dias.id=lugares_horarios.iddia")

	db = db.Where(" idlugar = ? ", c.Param("id"))

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("id")
	}

	// Ejecuto consultas
	var horarios []Modelos.Lugares_horarios
	db.Find(&horarios)
	//db.Table("lugares_horarios").Count(&registros)
	data := Data{Horarios: horarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetHorario(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	horarios := new(Modelos.Lugares_horarios)

	// Armo select
	db = db.Select("lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas")
	db = db.Joins("JOIN dias ON dias.id=lugares_horarios.iddia")

	db.Find(&horarios, id)

	data := Data{Horario: horarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	horarios := new(Modelos.Lugares_horarios)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(horarios); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("dia").Create(&horarios).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Horario: horarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	horarios := new(Modelos.Lugares_horarios)
	if err := c.Bind(horarios); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("dia").Save(&horarios).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Horario: horarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM lugares_horarios WHERE id = ?", c.Param("id")).Error; err != nil {
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
