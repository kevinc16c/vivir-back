package palabras_clave

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
	Registros      int                      `json:"registros,omitempty"`
	Palabras_clave []Modelos.Palabras_clave `json:"palabras_clave,omitempty"`
	Palabra_clave  *Modelos.Palabras_clave  `json:"palabra_clave,omitempty"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("palabras_clave.id,palabras_clave.palabraclave,palabras_clave.idrubro,rubros.descrirubro,palabras_clave.estado")
	db = db.Joins("JOIN rubros ON rubros.id=palabras_clave.idrubro")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" id like ? ", "%"+c.QueryParam("query")+"%").
			Or(" palabraclave like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("palabraclave")
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
	var palabras []Modelos.Palabras_clave
	db.Offset(offset).Limit(limite).Find(&palabras)
	db.Table("palabras_clave").Count(&registros)
	data := Data{Registros: registros, Palabras_clave: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("palabras_clave.id,palabras_clave.palabraclave,palabras_clave.idrubro,rubros.descrirubro,palabras_clave.estado")
	db = db.Joins("JOIN rubros ON rubros.id=palabras_clave.idrubro")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" idrubro = ? AND palabraclave LIKE ? AND estado<>'B'", c.Param("idrubro"), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idrubro = ? AND estado<>'B'", c.Param("idrubro"))
	}

	db = db.Order("palabraclave")

	// Ejecuto consulta
	var palabras []Modelos.Palabras_clave
	db.Find(&palabras)
	data := Data{Palabras_clave: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPalabraClave(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	palabras := new(Modelos.Palabras_clave)
	db.Find(&palabras, id)

	data := Data{Palabra_clave: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	palabras := new(Modelos.Palabras_clave)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(palabras); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("estado").Omit("descrirubro").Create(&palabras).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Palabra_clave: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	palabras := new(Modelos.Palabras_clave)
	if err := c.Bind(palabras); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("estado").Omit("descrirubro").Save(&palabras).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Palabra_clave: palabras}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE palabras_clave SET estado='B', fechabaja=? WHERE id = ?", time.Now(), c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE palabras_clave SET estado='' WHERE id = ?", c.Param("id")).Error; err != nil {
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
