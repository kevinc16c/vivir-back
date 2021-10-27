package productos_categorias

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
	Registros            int                            `json:"registros,omitempty"`
	Productos_categorias []Modelos.Productos_categorias `json:"categorias,omitempty"`
	Producto_categoria   *Modelos.Productos_categorias  `json:"categoria,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	db = db.Joins("join rubros on rubros.id=productos_categorias.idrubro")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" productos_categorias.id like ? OR descricatprod like ? OR rubros.descrirubro like ?", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("descricatprod")
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
	var categorias []Modelos.Productos_categorias
	db.Preload("Rubro").Offset(offset).Limit(limite).Find(&categorias)
	db.Table("productos_categorias").Count(&registros)
	data := Data{Registros: registros, Productos_categorias: categorias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func CategoriasRubro(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where("idrubro = ? and descricatprod LIKE ? and estado<>'B'", (c.Param("idrubro")), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where("idrubro = ? and estado<>'B'", c.Param("idrubro"))
	}

	db = db.Order("descricatprod")

	// Ejecuto consulta
	var categorias []Modelos.Productos_categorias
	db.Preload("Rubro").Find(&categorias)
	data := Data{Productos_categorias: categorias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetCategoria(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	categorias := new(Modelos.Productos_categorias)
	db.Preload("Rubro").Find(&categorias, id)

	data := Data{Producto_categoria: categorias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	categorias := new(Modelos.Productos_categorias)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(categorias); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&categorias).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Producto_categoria: categorias}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	categorias := new(Modelos.Productos_categorias)
	if err := c.Bind(categorias); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&categorias).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Producto_categoria: categorias}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE productos_categorias SET estado='B' WHERE id = ?", c.Param("id")).Error; err != nil {
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

func Habilitar(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE productos_categorias SET estado='' WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro habilitado con éxito",
	})
}
