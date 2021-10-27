package provincias

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
	Registros  int                  `json:"registros,omitempty"`
	Provincias []Modelos.Provincias `json:"provincias,omitempty"`
	Provincia  *Modelos.Provincias  `json:"provincia,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idprovincia like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nombrepcia like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("nombrepcia")
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
	var provincias []Modelos.Provincias
	db.Preload("Pais").Offset(offset).Limit(limite).Find(&provincias)
	db.Table("provincias").Count(&registros)
	data := Data{Registros: registros, Provincias: provincias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" nombrepcia LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("nombrepcia")

	// Ejecuto consulta
	var provincias []Modelos.Provincias
	db.Preload("Pais").Find(&provincias)
	data := Data{Provincias: provincias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetProvincia(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("idprovincia")

	provincias := new(Modelos.Provincias)
	db.Preload("Pais").Find(&provincias, id)

	data := Data{Provincia: provincias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	provincias := new(Modelos.Provincias)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(provincias); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&provincias).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Provincia: provincias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	provincias := new(Modelos.Provincias)
	if err := c.Bind(provincias); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&provincias).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Provincia: provincias}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}
