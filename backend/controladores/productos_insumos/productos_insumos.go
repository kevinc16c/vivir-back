package productos_insumos

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
	Registros         int                         `json:"registros,omitempty"`
	Productos_insumos []Modelos.Productos_insumos `json:"insumos,omitempty"`
	Producto_insumo   *Modelos.Productos_insumos  `json:"insumo,omitempty"`
}

type Set_insumos struct {
	Insumos []Modelos.Productos_insumos
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idproducto = ? AND dinsuprodu like ?", id, "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idproducto = ?", id)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("dinsuprodu")
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
	var insumos []Modelos.Productos_insumos
	db.Offset(offset).Limit(limite).Find(&insumos)
	db.Table("productos_insumos").Count(&registros)
	data := Data{Registros: registros, Productos_insumos: insumos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func SetInsumos(c echo.Context) error {
	db := database.GetDb()

	set_insumos := new(Set_insumos)
	//var insumos []Set_insumos

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(set_insumos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// transaccion
	tr := db.Begin()

	// borro los inusmos anteriores
	if err := tr.Exec("DELETE FROM productos_insumos WHERE idproducto=?", c.Param("id")).Error; err != nil {
		tr.Rollback()
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registros en la tabla
	var insumos []Modelos.Productos_insumos
	insumos = set_insumos.Insumos
	for i := 0; i < len(insumos); i++ {
		if err := tr.Create(&insumos[i]).Error; err != nil {
			tr.Rollback()
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	tr.Commit()

	// Preparo mensaje de retorno
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Sabores, agregados, variedades asignadas con éxito.",
	})
}
