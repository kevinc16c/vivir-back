package insumos

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
	Registros int               `json:"registros,omitempty"`
	Insumos   []Modelos.Insumos `json:"insumos,omitempty"`
	Insumo    *Modelos.Insumos  `json:"insumo,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idlugar = ? AND (id like ? OR descripcion like ?)", id, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idlugar = ?", id)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("descripcion")
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
	var insumos []Modelos.Insumos
	db.Offset(offset).Limit(limite).Find(&insumos)
	db.Table("insumos").Count(&registros)
	data := Data{Registros: registros, Insumos: insumos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetInsumo(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	insumos := new(Modelos.Insumos)
	db.Find(&insumos, id)

	data := Data{Insumo: insumos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	insumos := new(Modelos.Insumos)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(insumos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&insumos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Insumo: insumos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	insumos := new(Modelos.Insumos)
	if err := c.Bind(insumos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&insumos).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Insumo: insumos}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	// transaccion
	tr := db.Begin()

	if err := tr.Exec("DELETE FROM productos_insumos WHERE idinsumo = ?", c.Param("id")).Error; err != nil {
		tr.Rollback()
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	if err := tr.Exec("DELETE FROM insumos WHERE id = ?", c.Param("id")).Error; err != nil {
		tr.Rollback()
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	tr.Commit()

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro eliminado con éxito",
	})
}

func Suspender(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE insumos SET suspendido=1 WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro marcado como suspendido con éxito",
	})
}

func Habilitar(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE insumos SET suspendido=0 WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro marcado como habilitado con éxito",
	})
}
