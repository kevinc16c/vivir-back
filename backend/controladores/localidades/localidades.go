package localidades

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
	Registros int `json:"registros,omitempty"`
	//Localidades      []Modelos.Localidades `json:"localidades,omitempty"`
	LocalidadesLista []LocalidadesLista   `json:"localidades,omitempty"`
	Localidad        *Modelos.Localidades `json:"localidad,omitempty"`
}

type LocalidadesLista struct {
	Idlocalidad  uint   `json:"idlocalidad" gorm:"primary_key"`
	Nombrelocali string `json:"nombrelocali"`
	Idprovincia  uint   `json:"idprovincia"`
	Nombrepcia   string `json:"nombrepcia"`
	Idpais       uint   `json:"idpais"`
	Nombrepais   string `json:"nombrepais"`
}

func (LocalidadesLista) TableName() string {
	return "localidades"
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("localidades.idlocalidad, localidades.nombrelocali, localidades.idprovincia, provincias.nombrepcia, provincias.idpais, paises.nombrepais")
	db = db.Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idlocalidad like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nombrelocali like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("nombrelocali")
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
	var localidades []LocalidadesLista
	db.Offset(offset).Limit(limite).Find(&localidades)
	db.Table("localidades").Count(&registros)
	data := Data{Registros: registros, LocalidadesLista: localidades}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("localidades.idlocalidad, localidades.nombrelocali, localidades.idprovincia, provincias.nombrepcia, provincias.idpais, paises.nombrepais")
	db = db.Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" nombrelocali LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("nombrelocali")

	// Ejecuto consulta
	var localidades []LocalidadesLista
	db.Find(&localidades)
	data := Data{LocalidadesLista: localidades}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetLocalidad(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	localidades := new(Modelos.Localidades)
	db.Preload("Provincia").Find(&localidades, id)

	data := Data{Localidad: localidades}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	localidades := new(Modelos.Localidades)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(localidades); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&localidades).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Localidad: localidades}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	localidades := new(Modelos.Localidades)
	if err := c.Bind(localidades); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&localidades).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Localidad: localidades}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}
