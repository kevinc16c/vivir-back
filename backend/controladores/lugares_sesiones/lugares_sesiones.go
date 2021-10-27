package lugares_sesiones

import (
	"net/http"
	"time"

	"../../database"
	Modelos "../../modelos"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	lugares_sesion := new(Modelos.Lugares_sesiones)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(lugares_sesion); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia del email
	lugar2 := new(Modelos.Lugares_sesiones)
	db.Where("idlugar = ? AND token = ?", lugares_sesion.Idlugar, lugares_sesion.Token).First(&lugar2)

	if lugar2.Id > 0 {

		// Actualizo token
		lugares_sesion.Fechamodif = time.Now() //utils.GetNow()
		if err := db.Exec("UPDATE lugares_sesiones SET token=?, fechamodif=? WHERE idlugar=? AND token=?", lugares_sesion.Token, lugares_sesion.Fechamodif, lugares_sesion.Idlugar, lugares_sesion.Token).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

	} else {

		// Inserta registro en la tabla
		lugares_sesion.Fechaalta = time.Now() //utils.GetNow()
		if err := db.Omit("fechamodif").Create(&lugares_sesion).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "exito",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	lugares_sesion := new(Modelos.Lugares_sesiones)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(lugares_sesion); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	if err := db.Exec("DELETE FROM lugares_sesiones WHERE idlugar=? AND token=?", lugares_sesion.Idlugar, lugares_sesion.Token).Error; err != nil {
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
