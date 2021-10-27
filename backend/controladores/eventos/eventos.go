package lugares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"../../database"
	Modelos "../../modelos"
	"../../utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros   int               `json:"registros,omitempty"`
	Eventos     []Modelos.Eventos `json:"eventos,omitempty"`
	Evento      *Modelos.Eventos  `json:"evento,omitempty"`
	Evento_alta *Eventos_alta     `json:"evento_alta,omitempty"`
}

type Eventos_alta struct {
	Idtipoeven   uint      `json:"idtipoeven"`
	Idlocalidad  uint      `json:"idlocalidad"`
	Latitud      float64   `json:"latitud"`
	Longitud     float64   `json:"longitud"`
	Evento       string    `json:"evento"`
	Descripcion  string    `json:"descripcion"`
	Direccion    string    `json:"direccion"`
	Telefono     string    `json:"telefono"`
	Celular      string    `json:"celular"`
	E_mail       string    `json:"e_mail"`
	Sitioweb     string    `json:"sitioweb"`
	Fechaalta    time.Time `json:"fechaalta,omitempty"`
	Inicioevento time.Time `json:"inicioevento"`
	Finalevento  time.Time `json:"finalevento"`
}

func (Eventos_alta) TableName() string {
	return "eventos"
}

type Red struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Nombrerrss string `json:"nombrerrss"`
	Rutaimgapp string `json:"rutaimgapp"`
	Urlrrss    string `json:"urlrrss"`
}

func (Red) TableName() string {
	return "lugares_rrss"
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" id like ? OR evento like ? OR tipo like ? OR descripcion like ?", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("evento")
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
	var eventos []Modelos.Eventos
	db.Preload("Tipo").Preload("Localidad.Provincia.Pais").Offset(offset).Limit(limite).Find(&eventos)
	db.Table("eventos").Count(&registros)
	data := Data{Registros: registros, Eventos: eventos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Eventos(c echo.Context) error {
	db := database.GetDb()
	fecha := time.Now() //utils.GetNow()

	// Controlo si vienen las coordenas del usuario en la url
	var lat string = "0.0"
	var long string = "0.0"
	var orden string = ""
	if c.QueryParam("location") != "" {
		coordenadas := strings.Split(c.QueryParam("location"), ",")
		lat = coordenadas[0]
		long = coordenadas[1]
		orden = "distancia"
	} else {
		lat = "0.0"
		long = "0.0"
		orden = "inicioevento"
	}

	// SELECT *******************************************************************
	db = db.Select("eventos.id,idtipoeven,idlocalidad,latitud,longitud,evento,descripcion,direccion,telefono,celular,e_mail,sitioweb,inicioevento,finalevento,(SELECT rutaimg FROM eventos_img WHERE idevento=eventos.id LIMIT 1) AS rutafoto,ST_Distance_Sphere(POINT(eventos.longitud, eventos.latitud), POINT(?, ?)) AS distancia", long, lat)
	db = db.Joins("JOIN tipos_eventos ON tipos_eventos.id=eventos.idtipoeven")

	// CONDICIONES **************************************************************
	strFecha := fmt.Sprintf("%v", fecha)
	var condicion string = "eventos.finalevento >= '" + strFecha + "'"

	// query
	if c.QueryParam("query") != "" {
		condicion = condicion + " AND (eventos.evento like " + "'%" + c.QueryParam("query") + "%'" + " OR tipoos_eventos.tipo like " + "'%" + c.QueryParam("query") + "%'" + " OR descripcion like " + "'%" + c.QueryParam("query") + "%'" + ")"
	}
	db = db.Where(condicion)

	// ORDENES ******************************************************************
	db = db.Order(orden)

	// PAGINACION ***************************************************************
	var pagina uint = 1
	var limite uint = 20
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
	var eventos []Modelos.Eventos
	db.Preload("Tipo").Preload("Localidad.Provincia.Pais").Offset(offset).Limit(limite).Find(&eventos)
	db.Table("eventos").Count(&registros)
	data := Data{Registros: registros, Eventos: eventos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetEvento(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	eventos := new(Modelos.Eventos)
	db.Preload("Tipo").Preload("Localidad.Provincia.Pais").Preload("Redes", func(db *gorm.DB) *gorm.DB {
		db = db.Select("eventos_rrss.id,eventos_rrss.idevento,eventos_rrss.idrrss,eventos_rrss.descriprrss,redes_sociales.rutaimgapp,eventos_rrss.urlrrss")
		db = db.Joins("LEFT JOIN redes_sociales ON redes_sociales.id=eventos_rrss.idrrss")
		return db
	}).Preload("Imagenes").Find(&eventos, id)

	data := Data{Evento: eventos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	eventos := new(Eventos_alta)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(eventos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	eventos.Fechaalta = time.Now() //utils.GetNow()
	if err := db.Create(&eventos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Evento_alta: eventos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	eventos := new(Modelos.Eventos)
	if err := c.Bind(eventos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("fechaalta", "rutafoto", "distancia").Save(&eventos).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Evento: eventos}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM eventos WHERE id = ?", c.Param("id")).Error; err != nil {
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
