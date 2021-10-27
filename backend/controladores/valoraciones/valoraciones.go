package valoraciones

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
	Registros       int                    `json:"registros,omitempty"`
	PuntuacionLugar float64                `json:"puntuacion,omitempty"`
	Valoraciones    []Modelos.Valoraciones `json:"valoraciones,omitempty"`
	Valoracion      *Modelos.Valoraciones  `json:"valoracion,omitempty"`
}

type Puntuacion struct {
	Puntuacion float64 `json:"puntuacion,omitempty"`
}

type Valoracion struct {
	Id          uint      `json:"id"`
	Idlugar     uint      `json:"idlugar"`
	Idusuario   uint      `json:"idusuario"`
	Puntuacion  float64   `json:"puntuacion"`
	Fecha       time.Time `json:"fecha"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Fechamodif  time.Time `json:"fechamodif,omitempty"`
}

func (Valoracion) TableName() string {
	return "valoraciones"
}

func PaginacionValoracionLugar(c echo.Context) error {
	db := database.GetDb()
	idlugar := c.Param("id")

	db = db.Select("valoraciones.id,valoraciones.idlugar,lugares.nombrelugar,valoraciones.idusuario,usuarios_app.apellido,usuarios_app.nombres,usuarios_app.rutaimgusu,valoraciones.puntuacion,valoraciones.fecha,valoraciones.titulo,valoraciones.descripcion,valoraciones.fechamodif")
	db = db.Joins("JOIN lugares ON lugares.idlugar=valoraciones.idlugar").
		Joins("JOIN usuarios_app ON usuarios_app.id=valoraciones.idusuario")

	db = db.Where(" valoraciones.idlugar = ?", idlugar)

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fecha DESC")
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
	var valoraciones []Modelos.Valoraciones
	db.Offset(offset).Limit(limite).Find(&valoraciones)
	db.Table("valoraciones").Count(&registros)

	var puntuacion = new(Puntuacion)
	db2 := database.GetDb()
	db2.Raw("SELECT IFNULL(ROUND(SUM(puntuacion)/COUNT(id), 2), 0.00) AS puntuacion FROM valoraciones WHERE idlugar=?", idlugar).
		Scan(&puntuacion)

	data := Data{Registros: registros, PuntuacionLugar: puntuacion.Puntuacion, Valoraciones: valoraciones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetValoracionUsuarioLugar(c echo.Context) error {
	db := database.GetDb()
	idusuario := c.Param("idusuario")
	idlugar := c.Param("idlugar")

	valoraciones := new(Modelos.Valoraciones)

	db = db.Select("valoraciones.id,valoraciones.idlugar,lugares.nombrelugar,valoraciones.idusuario,usuarios_app.apellido,usuarios_app.nombres,valoraciones.puntuacion,valoraciones.fecha,valoraciones.titulo,valoraciones.descripcion,valoraciones.fechamodif")
	db = db.Joins("JOIN lugares ON lugares.idlugar=valoraciones.idlugar").
		Joins("JOIN usuarios_app ON usuarios_app.id=valoraciones.idusuario")

	db.Where("valoraciones.idusuario=? AND valoraciones.idlugar=?", idusuario, idlugar).First(&valoraciones)

	data := Data{Valoracion: valoraciones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	valoracion := new(Valoracion)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(valoracion); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia opiniones del usuario para el lugar
	var valoraciones []Modelos.Valoraciones
	db.Where("idusuario = ? AND idlugar = ?", valoracion.Idusuario, valoracion.Idlugar).First(&valoraciones)
	if len(valoraciones) > 0 {

		// Actualiza el registro
		valoracion.Fechamodif = time.Now() //utils.GetNow()
		if err := db.Exec("UPDATE valoraciones SET puntuacion=?, descripcion=?, fechamodif=? WHERE idusuario=? AND idlugar=?", valoracion.Puntuacion, valoracion.Descripcion, valoracion.Fechamodif, valoracion.Idusuario, valoracion.Idlugar).Error; err != nil {
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		// Preparo mensaje de retorno
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Tu opinión fue actualizada correctamente. ",
		})

	} else {

		// Inserta registro en la tabla
		valoracion.Fecha = time.Now() //utils.GetNow()
		if err := db.Omit("Fechamodif").Create(&valoracion).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		// Preparo mensaje de retorno
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Tu opinión fue registrada con éxito",
		})
	}

}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM valoraciones WHERE id = ?", c.Param("id")).Error; err != nil {
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
