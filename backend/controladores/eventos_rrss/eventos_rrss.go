package eventos_rrss

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
	Registros    int                    `json:"registros,omitempty"`
	Eventos_rrss []Modelos.Eventos_rrss `json:"redes_sociales,omitempty"`
	Evento_rrss  *Modelos.Eventos_rrss  `json:"red_social,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("eventos_rrss.id,eventos_rrss.idevento,eventos_rrss.idrrss,redes_sociales.nombrerrss,eventos_rrss.descriprrss,eventos_rrss.urlrrss")
	db = db.Joins("JOIN redes_sociales ON redes_sociales.id=eventos_rrss.idrrss")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idevento = ? and nombrerrss like ?", (c.Param("id")), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idevento = ? ", c.Param("id"))
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("id")
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
	var redes []Modelos.Eventos_rrss
	db.Offset(offset).Limit(limite).Find(&redes)
	db.Table("eventos_rrss").Count(&registros)
	data := Data{Registros: registros, Eventos_rrss: redes}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetRed(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	redes := new(Modelos.Eventos_rrss)
	db.Find(&redes, id)

	data := Data{Evento_rrss: redes}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	redes := new(Modelos.Eventos_rrss)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(redes); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("Nombrerrss").Create(&redes).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Evento_rrss: redes}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	redes := new(Modelos.Eventos_rrss)
	if err := c.Bind(redes); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("Nombrerrss").Save(&redes).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Evento_rrss: redes}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM eventos_rrss WHERE id = ?", c.Param("id")).Error; err != nil {
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
