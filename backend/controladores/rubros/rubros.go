package rubros

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
	Registros int              `json:"registros,omitempty"`
	Rubros    []Modelos.Rubros `json:"rubros,omitempty"`
	Rubro     *Modelos.Rubros  `json:"rubro,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" id like ? ", "%"+c.QueryParam("query")+"%").
			Or(" descrirubro like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("descrirubro")
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
	var rubros []Modelos.Rubros
	db.Offset(offset).Limit(limite).Find(&rubros)
	db.Table("rubros").Count(&registros)
	data := Data{Registros: registros, Rubros: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func RubrosSitiosUtiles(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	db = db.Where("sitioutil=1")

	db = db.Order("descrirubro")

	// Ejecuto consulta
	var rubros []Modelos.Rubros
	db.Find(&rubros)
	data := Data{Rubros: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" descrirubro LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("descrirubro")

	// Ejecuto consulta
	var rubros []Modelos.Rubros
	db.Find(&rubros)
	data := Data{Rubros: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetRubro(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	rubros := new(Modelos.Rubros)
	db.Find(&rubros, id)

	data := Data{Rubro: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	rubros := new(Modelos.Rubros)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(rubros); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&rubros).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Rubro: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	rubros := new(Modelos.Rubros)
	if err := c.Bind(rubros); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&rubros).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Rubro: rubros}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM rubros WHERE id = ?", c.Param("id")).Error; err != nil {
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
