package subrubros

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
	Registros int                 `json:"registros,omitempty"`
	Subrubros []Modelos.Subrubros `json:"subrubros,omitempty"`
	Subrubro  *Modelos.Subrubros  `json:"subrubro,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Joins("join rubros on rubros.id=subrubros.idrubro")
		db = db.Where(" subrubros.id like ? OR dsubrubro like ? OR rubros.descrirubro like ?", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("dsubrubro")
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
	var subrubros []Modelos.Subrubros
	db.Preload("Rubro").Offset(offset).Limit(limite).Find(&subrubros)
	db.Table("subrubros").Count(&registros)
	data := Data{Registros: registros, Subrubros: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where("dsubrubro LIKE ?", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("dsubrubro")

	// Ejecuto consulta
	var subrubros []Modelos.Subrubros
	db.Find(&subrubros)
	data := Data{Subrubros: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func SubrubrosRubro(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where("idrubro = ? and dsubrubro LIKE ?", (c.Param("idrubro")), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where("idrubro = ?", c.Param("idrubro"))
	}

	db = db.Order("dsubrubro")

	// Ejecuto consulta
	var subrubros []Modelos.Subrubros
	db.Find(&subrubros)
	data := Data{Subrubros: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetSubrubro(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	subrubros := new(Modelos.Subrubros)
	db.Preload("Rubro").Find(&subrubros, id)

	data := Data{Subrubro: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	subrubros := new(Modelos.Subrubros)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(subrubros); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&subrubros).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Subrubro: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	subrubros := new(Modelos.Subrubros)
	if err := c.Bind(subrubros); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&subrubros).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Subrubro: subrubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}
