package usuarios_sesiones

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

	usuario_sesion := new(Modelos.Usuarios_sesiones)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(usuario_sesion); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia del email
	usuario2 := new(Modelos.Usuarios_sesiones)
	db.Where("idusuario = ? AND token = ?", usuario_sesion.Idusuario, usuario_sesion.Token).First(&usuario2)

	if usuario2.Id > 0 {

		// Actualizo token
		usuario_sesion.Fechamodif = time.Now() //utils.GetNow()
		if err := db.Exec("UPDATE usuarios_app_sesiones SET token=?, fechamodif=? WHERE idusuario=? AND token=?", usuario_sesion.Token, usuario_sesion.Fechamodif, usuario_sesion.Idusuario, usuario_sesion.Token).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

	} else {

		// Inserta registro en la tabla
		usuario_sesion.Fechaalta = time.Now() //utils.GetNow()
		if err := db.Omit("fechamodif").Create(&usuario_sesion).Error; err != nil {
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

	usuario_sesion := new(Modelos.Usuarios_sesiones)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(usuario_sesion); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	if err := db.Exec("DELETE FROM usuarios_app_sesiones WHERE idusuario=? AND token=?", usuario_sesion.Idusuario, usuario_sesion.Token).Error; err != nil {
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
