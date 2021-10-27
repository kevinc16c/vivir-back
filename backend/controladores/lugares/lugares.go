package lugares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	config "../../config"
	"../../database"
	Modelos "../../modelos"
	"../../utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/jung-kurt/gofpdf"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros           int                        `json:"registros,omitempty"`
	Lugares             []Lugar                    `json:"lugares,omitempty"`
	Lugares_tipo        []Lugares_tipo             `json:"lugares_tipo,omitempty"`
	Lugar               *Lugar                     `json:"lugar,omitempty"`
	LugarAlta           *Modelos.Lugares           `json:"lugar_alta,omitempty"`
	Horarios            []Modelos.Lugares_horarios `json:"horarios,omitempty"`
	Redes               []Red                      `json:"redes,omitempty"`
	Productos_categoria []Productos_categoria      `json:"productos_categorias,omitempty"`
	Promociones         []Promociones_lugar        `json:"promociones,omitempty"`
	Valoraciones        *Modelos.Valoraciones      `json:"valoraciones,omitempty"`
}

type Lugar struct {
	Idlugar        uint                `json:"idlugar"`
	Idpropietario  uint                `json:"idpropietario,omitempty"`
	Propietario    string              `json:"propietario,omitempty"`
	Idrubro        uint                `json:"idrubro"`
	Descrirubro    string              `json:"descrirubro"`
	Idsubrubro     uint                `json:"idsubrubro"`
	Dsubrubro      string              `json:"dsubrubro"`
	Idtipolugar    uint                `json:"idtipolugar"`
	Tipolugar      string              `json:"tipolugar"`
	Urlmarker      string              `json:"urlmarker,omitempty"`
	Idtipoconv     uint                `json:"idtipoconv,omitempty"`
	Desctconv      string              `json:"desctconv,omitempty"`
	Direccion      string              `json:"direccion"`
	Idlocalidad    uint                `json:"idlocalidad"`
	Nombrelocali   string              `json:"nombrelocali"`
	Idprovincia    uint                `json:"idprovincia,omitempty"`
	Nombrepcia     string              `json:"nombrepcia,omitempty"`
	Idpais         uint                `json:"idpais,omitempty"`
	Nombrepais     string              `json:"nombrepais,omitempty"`
	Latitud        float64             `json:"latitud"`
	Longitud       float64             `json:"longitud"`
	Nombrelugar    string              `json:"nombrelugar"`
	Telefono       string              `json:"telefono"`
	Celular        string              `json:"celular"`
	E_mail         string              `json:"e_mail"`
	Sitioweb       string              `json:"sitioweb"`
	Describreve    string              `json:"describreve"`
	Rutafoto       string              `json:"rutafoto"`
	Ofertas        uint                `json:"ofertas"`
	Idtipodelivery uint                `json:"idtipodelivery"`
	Tipodelivery   string              `json:"tipodelivery"`
	Precdelivery   float64             `json:"precdelivery"`
	Conpddiferido  uint                `json:"conpddiferido"`
	Cpraminima     float64             `json:"cpraminima"`
	Fechaalta      time.Time           `json:"fechaalta,omitempty"`
	Vencimiento    time.Time           `json:"vencimiento,omitempty"`
	Porcomision    float64             `json:"porcomision"`
	Activo         uint                `json:"activo"`
	Qrasignado     uint                `json:"qrasignado"`
	Abierto        uint                `json:"abierto"`
	Deturno        uint                `json:"deturno"`
	Isfavorito     uint                `json:"isfavorito"`
	Pagoonline     uint                `json:"pagoonline"`
	Publickeymp    string              `json:"publickeymp"`
	Distancia      float64             `json:"distancia,omitempty"`
	Valoraciones   uint                `json:"valoraciones"`
	Puntuacion     float64             `json:"puntuacion"`
	Productos      []Modelos.Productos `json:"productos,omitempty"`
}

func (Lugar) TableName() string {
	return "lugares"
}

type Lugares_tipo struct {
	Id      uint    `json:"id" gorm:"primary_key"`
	Tipo    string  `json:"tipo"`
	Lugares []Lugar `json:"lugares"`
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

type Productos_lugar struct {
	Id          uint    `json:"id" gorm:"primary_key"`
	Codintprod  string  `json:"codintprod"`
	Idcategprod uint    `json:"idcategprod"`
	Descriprod  string  `json:"descriprod"`
	Desextprod  string  `json:"desextprod"`
	Prunitprod  float64 `json:"prunitprod"`
	Aliivaprod  float64 `json:"Aliivaprod"`
	Suspendido  uint    `json:"suspendido"`
	Rutaimagen  string  `json:"imagen"`
}

func (Productos_lugar) TableName() string {
	return "productos"
}

type Categorias_productos_lugar struct {
	Id            uint   `json:"id" gorm:"primary_key"`
	Descricatprod string `json:"descricatprod"`
}

func (Categorias_productos_lugar) TableName() string {
	return "productos_categorias"
}

type Promociones_lugar struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Vencimiento time.Time `json:"vencimiento"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Terminos    string    `json:"terminos"`
	Cuposdispon uint      `json:"cuposdispon"`
	Rutaimg     string    `json:"rutaimg"`
}

func (Promociones_lugar) TableName() string {
	return "promociones"
}

type Productos_categoria struct {
	Id        uint              `json:"id" gorm:"primary_key"`
	Categoria string            `json:"categoria"`
	Productos []Productos_lugar `json:"productos"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("lugares.idlugar,lugares.idpropietario,propietarios.razonsocial AS propietario,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN personas_humanas_juridicas AS propietarios ON propietarios.id=lugares.idpropietario").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idlugar like ? OR nombrelugar like ? OR descrirubro like ? OR dsubrubro like ? OR tipolugar like ? ", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("nombrelugar")
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
	var lugares []Lugar
	db.Offset(offset).Limit(limite).Find(&lugares)
	db.Table("lugares").Count(&registros)
	data := Data{Registros: registros, Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PaginacionLugaresPropietario(c echo.Context) error {
	db := database.GetDb()
	idpropietario := c.Param("id")

	db = db.Select("lugares.idlugar,lugares.idpropietario,propietarios.razonsocial AS propietario,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN personas_humanas_juridicas AS propietarios ON propietarios.id=lugares.idpropietario").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idpropietario = ? AND (idlugar like ? OR nombrelugar like ? OR descrirubro like ? OR dsubrubro like ? OR tipolugar like ?)", idpropietario, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idpropietario = ?", idpropietario)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("idlugar")
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
	var lugares []Lugar
	db.Offset(offset).Limit(limite).Find(&lugares)
	db.Table("lugares").Count(&registros)
	data := Data{Registros: registros, Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" nombrelugar LIKE ? AND estado<>'B'", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" estado<>'B'")
	}

	db = db.Order("nombrelugar")

	// Ejecuto consulta
	var lugares []Lugar
	db.Find(&lugares)
	data := Data{Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func ListaLugaresPropietario(c echo.Context) error {
	db := database.GetDb()
	idpropietario := c.Param("id")

	// Armo select
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" idpropietario = ? AND nombrelugar LIKE ? AND estado<>'B'", idpropietario, "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idpropietario = ? AND estado<>'B'", idpropietario)
	}

	db = db.Order("nombrelugar")

	// Ejecuto consulta
	var lugares []Lugar
	db.Find(&lugares)
	data := Data{Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func ListaLugaresSubrubros(c echo.Context) error {
	db := database.GetDb()
	idsubrubro := c.Param("id")

	// Armo select
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,tipos_lugares.rutaicono AS urlmarker,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" idsubrubro = ? AND nombrelugar LIKE ? AND estado<>'B'", idsubrubro, "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idsubrubro = ? AND estado<>'B'", idsubrubro)
	}

	db = db.Order("nombrelugar")

	// Ejecuto consulta
	var lugares []Lugar
	db.Find(&lugares)
	data := Data{Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func LugaresBusqueda(c echo.Context) error {
	db := database.GetDb()
	fecha := time.Now() //utils.GetNow()

	// Controlo si vienen las coordenas del usuario en la url
	var lat string = "0.0"
	var long string = "0.0"
	if c.QueryParam("location") != "" {
		coordenadas := strings.Split(c.QueryParam("location"), ",")
		lat = coordenadas[0]
		long = coordenadas[1]
	} else {
		lat = "0.0"
		long = "0.0"
	}

	// busco tipos de lugar
	var tipos_lugares []Modelos.Tipos_lugares
	db.Raw("SELECT id,tipolugar FROM tipos_lugares").Scan(&tipos_lugares)

	// SELECT *******************************************************************
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,tipos_lugares.rutaicono AS urlmarker,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar AND vencimiento>=? AND idtipopromo=3 AND estado<>'B') AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado,(SELECT COUNT(id) FROM valoraciones WHERE idlugar=lugares.idlugar) AS valoraciones,(SELECT SUM(puntuacion)/valoraciones FROM valoraciones WHERE idlugar=lugares.idlugar) AS puntuacion,(SELECT COUNT(id) FROM turnos_farmacias WHERE idlugar=lugares.idlugar AND '"+fmt.Sprintf("%v", fecha)+"' BETWEEN inicioturno AND finalturno) AS deturno,ST_Distance_Sphere(POINT(lugares.longitud, lugares.latitud), POINT(?, ?)) AS distancia", fecha, long, lat)
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery").
		Joins("LEFT JOIN (SELECT lugares_plc.idlugar,palabras_clave.palabraclave FROM lugares_plc JOIN palabras_clave ON palabras_clave.id=lugares_plc.idpalabraclave) palabras ON palabras.idlugar=lugares.idlugar")

	// CONDICIONES **************************************************************
	strFecha := fmt.Sprintf("%v", fecha.Format("2006-01-02"))
	var condicion string = "lugares.estado<>'B' AND lugares.vencimiento >= '" + strFecha + "'"

	// query
	condicion = condicion + " AND (lugares.nombrelugar like " + "'%" + c.QueryParam("query") + "%'" + " OR descrirubro like " + "'%" + c.QueryParam("query") + "%'" + " OR dsubrubro like " + "'%" + c.QueryParam("query") + "%'" + " OR palabras.palabraclave like " + "'%" + c.QueryParam("query") + "%')"

	// WHERE ********************************************************************
	db = db.Where(condicion)

	// GROUP BY *****************************************************************
	db = db.Group("lugares.idlugar")

	// ORDENES ******************************************************************
	var orden string = ""
	if c.QueryParam("order") == "distancia" {
		orden = "distancia"
	}

	if c.QueryParam("order") == "nombre" {
		orden = "nombrelugar"
	}

	db = db.Order(orden)

	// Ejecuto consultas
	var lugares []Lugar
	db.Find(&lugares)

	// Corte de control de lugares segun tipo de lugar
	var lugares_tipo []Lugares_tipo
	for i := 0; i < len(tipos_lugares); i++ {
		var lugaresTipo []Lugar
		for j := 0; j < len(lugares); j++ {
			if lugares[j].Idtipolugar == tipos_lugares[i].Id {
				lugaresTipo = append(lugaresTipo, lugares[j])
			}
		}

		// Lugares que corresponden a un tipo de lugar
		if len(lugaresTipo) > 0 {
			lug_tipo := new(Lugares_tipo)
			lug_tipo.Id = tipos_lugares[i].Id
			lug_tipo.Tipo = tipos_lugares[i].Tipolugar
			lug_tipo.Lugares = lugaresTipo
			lugares_tipo = append(lugares_tipo, *lug_tipo)
		}
	}

	data := Data{Lugares_tipo: lugares_tipo}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func LugaresTipo(c echo.Context) error {
	db := database.GetDb()
	tipo := c.QueryParam("tipo")
	fecha := time.Now() //utils.GetNow()

	// control de filtro
	idtipo := 0
	if tipo == "alojamiento" {
		idtipo = 1
	} else if tipo == "gastronomia" {
		idtipo = 2
	} else if tipo == "entretenimiento" {
		idtipo = 3
	} else if tipo == "atractivos" {
		idtipo = 4
	} else if tipo == "comercios" {
		idtipo = 5
	} else if tipo == "utilidades" {
		idtipo = 6
	}

	// Controlo si vienen las coordenas del usuario en la url
	var lat string = "0.0"
	var long string = "0.0"
	if c.QueryParam("location") != "" {
		coordenadas := strings.Split(c.QueryParam("location"), ",")
		lat = coordenadas[0]
		long = coordenadas[1]
	} else {
		lat = "0.0"
		long = "0.0"
	}

	// SELECT *******************************************************************
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,tipos_lugares.rutaicono AS urlmarker,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar AND vencimiento>=? AND idtipopromo=3 AND estado<>'B') AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado,(SELECT COUNT(id) FROM valoraciones WHERE idlugar=lugares.idlugar) AS valoraciones,(SELECT SUM(puntuacion)/valoraciones FROM valoraciones WHERE idlugar=lugares.idlugar) AS puntuacion,(SELECT COUNT(id) FROM turnos_farmacias WHERE idlugar=lugares.idlugar AND '"+fmt.Sprintf("%v", fecha)+"' BETWEEN inicioturno AND finalturno) AS deturno,ST_Distance_Sphere(POINT(lugares.longitud, lugares.latitud), POINT(?, ?)) AS distancia", fecha, long, lat)
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery").
		Joins("LEFT JOIN (SELECT lugares_plc.idlugar,palabras_clave.palabraclave FROM lugares_plc JOIN palabras_clave ON palabras_clave.id=lugares_plc.idpalabraclave) palabras ON palabras.idlugar=lugares.idlugar")

	// CONDICIONES **************************************************************
	strFecha := fmt.Sprintf("%v", fecha.Format("2006-01-02"))
	var condicion string = "lugares.idtipolugar = " + fmt.Sprintf("%v", idtipo) + " AND lugares.estado<>'B' AND lugares.vencimiento >= '" + strFecha + "'"

	// query
	if c.QueryParam("query") != "" {
		condicion = condicion + " AND (lugares.nombrelugar like " + "'%" + c.QueryParam("query") + "%'" + " OR dsubrubro like " + "'%" + c.QueryParam("query") + "%'" + " OR palabras.palabraclave like " + "'%" + c.QueryParam("query") + "%')"
	}

	// localidad
	if c.QueryParam("idloc") != "" {
		condicion = condicion + " AND lugares.idlocalidad = " + c.QueryParam("idloc")
	}

	// subrubro
	if c.QueryParam("idsubr") != "" {
		condicion = condicion + " AND lugares.idsubrubro = " + c.QueryParam("idsubr")
	}

	// lugares con promociones/ofertas
	var having string = ""
	if c.QueryParam("isoffer") == "true" {
		having = having + "ofertas > 0"
		db = db.Having(having)
	}

	db = db.Where(condicion)

	// GROUP BY *****************************************************************
	db = db.Group("lugares.idlugar")

	// ORDENES ******************************************************************
	var orden string = ""
	if c.QueryParam("order") == "distancia" {
		orden = "distancia"
	}

	if c.QueryParam("order") == "puntuacion" {
		orden = "valoraciones DESC"
	}

	if c.QueryParam("order") == "ofertas" {
		orden = "ofertas DESC"
	}

	if c.QueryParam("order") == "nombre" {
		orden = "nombrelugar"
	}

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
	var lugares []Lugar
	db.Offset(offset).Limit(limite).Find(&lugares)
	db.Table("lugares").Count(&registros)
	data := Data{Registros: registros, Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func LugaresRubro(c echo.Context) error {
	db := database.GetDb()
	idrubro := c.QueryParam("idrubro")
	fecha := time.Now() //utils.GetNow()

	// Controlo si vienen las coordenas del usuario en la url
	var lat string = "0.0"
	var long string = "0.0"
	if c.QueryParam("location") != "" {
		coordenadas := strings.Split(c.QueryParam("location"), ",")
		lat = coordenadas[0]
		long = coordenadas[1]
	} else {
		lat = "0.0"
		long = "0.0"
	}

	// SELECT *******************************************************************
	//strFechaTurno := fmt.Sprintf("%v", fecha.Format("2006-01-02 15:04:05"))
	db = db.Select("lugares.idlugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,tipos_lugares.rutaicono AS urlmarker,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar AND vencimiento>=? AND idtipopromo=3 AND estado<>'B') AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado,(SELECT COUNT(id) FROM valoraciones WHERE idlugar=lugares.idlugar) AS valoraciones,(SELECT SUM(puntuacion)/valoraciones FROM valoraciones WHERE idlugar=lugares.idlugar) AS puntuacion,(SELECT COUNT(id) FROM turnos_farmacias WHERE idlugar=lugares.idlugar AND '"+fmt.Sprintf("%v", fecha)+"' BETWEEN inicioturno AND finalturno) AS deturno,ST_Distance_Sphere(POINT(lugares.longitud, lugares.latitud), POINT(?, ?)) AS distancia", fecha, long, lat)
	db = db.Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery").
		Joins("LEFT JOIN (SELECT lugares_plc.idlugar,palabras_clave.palabraclave FROM lugares_plc JOIN palabras_clave ON palabras_clave.id=lugares_plc.idpalabraclave) palabras ON palabras.idlugar=lugares.idlugar")

	// CONDICIONES **************************************************************
	strFecha := fmt.Sprintf("%v", fecha.Format("2006-01-02"))
	var condicion string = "lugares.idrubro = " + fmt.Sprintf("%v", idrubro) + " AND lugares.estado<>'B' AND lugares.vencimiento >= '" + strFecha + "'"

	// query
	if c.QueryParam("query") != "" {
		condicion = condicion + " AND (lugares.nombrelugar like " + "'%" + c.QueryParam("query") + "%'" + " OR dsubrubro like " + "'%" + c.QueryParam("query") + "%'" + " OR palabras.palabraclave like " + "'%" + c.QueryParam("query") + "%')"
	}

	db = db.Where(condicion)

	// GROUP BY *****************************************************************
	db = db.Group("lugares.idlugar")

	// ORDENES ******************************************************************
	var orden string = ""
	if c.QueryParam("order") == "distancia" {
		orden = "deturno DESC, distancia"
	}

	if c.QueryParam("order") == "nombre" {
		orden = "deturno DESC, nombrelugar"
	}

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
	var lugares []Lugar
	db.Offset(offset).Limit(limite).Find(&lugares)
	db.Table("lugares").Count(&registros)
	data := Data{Registros: registros, Lugares: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetLugar(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	lugares := new(Lugar)
	// db.Preload("Localidad.Provincia.Pais").
	// 	Preload("Personas_humanas_juridicas.Condicioniva").
	// 	Preload("Personas_humanas_juridicas.Localidad.Provincia.Pais").
	// 	Preload("Horarios", func(db *gorm.DB) *gorm.DB {
	// 		db = db.Select("lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas")
	// 		db = db.Joins("JOIN dias ON dias.id=lugares_horarios.iddia")
	// 		return db
	// 	}).
	// 	Preload("Redes", func(db *gorm.DB) *gorm.DB {
	// 		db = db.Select("lugares_rrss.id,lugares_rrss.idlugar,lugares_rrss.idrrss,redes_sociales.nombrerrss,lugares_rrss.descriprrss,lugares_rrss.urlrrss")
	// 		db = db.Joins("JOIN redes_sociales ON redes_sociales.id=lugares_rrss.idrrss")
	// 		return db
	// 	}).
	// 	Preload("Imagenes").
	// 	Find(&lugares, id)

	db = db.Select("lugares.idlugar,lugares.idpropietario,propietarios.razonsocial AS propietario,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.idtipoconv,tipos_convenio.desctconv,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar LIMIT 1) AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,fechaalta,vencimiento,lugares.porcomision,lugares.activo,lugares.qrasignado")
	db = db.Joins("JOIN personas_humanas_juridicas AS propietarios ON propietarios.id=lugares.idpropietario").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar").
		Joins("JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery")
	db = db.Preload("Horarios", func(db *gorm.DB) *gorm.DB {
		db = db.Select("lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas")
		db = db.Joins("JOIN dias ON dias.id=lugares_horarios.iddia")
		return db
	})
	db = db.Preload("Productos", func(db *gorm.DB) *gorm.DB {
		db = db.Select("lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas")
		db = db.Joins("JOIN dias ON dias.id=lugares_horarios.iddia")
		return db
	})
	db = db.Where("lugares.idlugar=?", id).Find(&lugares)

	data := Data{Lugar: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetDetalleLugar(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")
	fecha := time.Now() //utils.GetNow()
	//strFecha := fmt.Sprintf("%v", fecha.Format("2006-01-02"))

	idUsuario := "0"
	if c.QueryParam("usuario") != "" {
		idUsuario = c.QueryParam("usuario")
	}

	// Lugar
	var lugar = new(Lugar)
	db.Raw("SELECT lugares.idlugar,lugares.idpropietario,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.idtipolugar,tipos_lugares.tipolugar,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,paises.idpais,paises.nombrepais,lugares.latitud,lugares.longitud,lugares.nombrelugar,lugares.telefono,lugares.celular,lugares.e_mail,lugares.sitioweb,lugares.describreve,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,(SELECT COUNT(id) FROM promociones WHERE idlugar=lugares.idlugar AND vencimiento>=? AND idtipopromo=3 AND estado<>'B') AS ofertas,lugares.idtipodelivery,tipos_delivery.tipodelivery,lugares.precdelivery,lugares.conpddiferido,lugares.cpraminima,lugares.porcomision,lugares.activo,lugares.qrasignado,(SELECT COUNT(id) FROM valoraciones WHERE idlugar=lugares.idlugar) AS valoraciones,(SELECT SUM(puntuacion)/valoraciones FROM valoraciones WHERE idlugar=lugares.idlugar) AS puntuacion,(SELECT IF(COUNT(id)>0, 1, 0) FROM lugares_horarios WHERE idlugar=lugares.idlugar AND iddia="+fmt.Sprintf("%v", c.QueryParam("dia"))+" AND '"+fmt.Sprintf("%v", c.QueryParam("hora"))+"' BETWEEN lughorades AND lughorahas) AS abierto,(SELECT IF(COUNT(id)>0, 1, 0) FROM lugares_cuentas_pago WHERE idlugar=lugares.idlugar AND suspendido=0 AND vencimiento>=?) AS pagoonline, (SELECT publickey FROM lugares_cuentas_pago WHERE idlugar=lugares.idlugar AND suspendido=0 AND vencimiento>=?) AS publickeymp,(SELECT IF(COUNT(id)>0, 1, 0) FROM usuarios_lugares_favoritos WHERE idlugar=lugares.idlugar AND idusuario="+fmt.Sprintf("%v", idUsuario)+") AS isfavorito "+
		"FROM lugares "+
		"JOIN rubros ON rubros.id=lugares.idrubro "+
		"JOIN subrubros ON subrubros.id=lugares.idsubrubro "+
		"LEFT JOIN tipos_lugares ON tipos_lugares.id=lugares.idtipolugar "+
		"JOIN tipos_convenio ON tipos_convenio.id=lugares.idtipoconv "+
		"JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad "+
		"JOIN provincias ON provincias.idprovincia=localidades.idprovincia "+
		"JOIN paises ON paises.idpais=provincias.idpais "+
		"LEFT JOIN tipos_delivery ON tipos_delivery.id=lugares.idtipodelivery "+
		"WHERE lugares.idlugar=?", fecha, fecha, fecha, id).Scan(&lugar)

	// busco horarios
	var horarios []Modelos.Lugares_horarios
	db.Raw("SELECT lugares_horarios.id,lugares_horarios.idlugar,lugares_horarios.iddia,dias.dia,lugares_horarios.lughorades,lugares_horarios.lughorahas "+
		"FROM lugares_horarios "+
		"JOIN dias ON dias.id=lugares_horarios.iddia "+
		"WHERE lugares_horarios.idlugar=?", id).Scan(&horarios)

	// busco redes sociales
	var redes []Red
	db.Raw("SELECT lugares_rrss.id,redes_sociales.nombrerrss,redes_sociales.rutaimgapp,lugares_rrss.urlrrss "+
		"FROM lugares_rrss "+
		"JOIN redes_sociales ON redes_sociales.id=lugares_rrss.idrrss "+
		"WHERE lugares_rrss.idlugar=?", id).Scan(&redes)

	// busco productos del lugar
	var busqueda string = ""
	if c.QueryParam("query") != "" {
		busqueda = " AND (productos.descriprod like " + "'%" + c.QueryParam("query") + "%')"
	}

	var productos []Productos_lugar
	db.Raw("SELECT productos.id,productos.codintprod,productos.idcategprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.suspendido,(SELECT rutaimgprod FROM productos_img WHERE idproducto=productos.id LIMIT 1) AS rutaimagen "+
		"FROM productos "+
		"WHERE productos.idlugar=? AND suspendido=0 AND estado<>'B'"+busqueda, id).Scan(&productos)

	// busco categorias de productos del lugar
	var categorias []Modelos.Productos_categorias
	db.Raw("SELECT id,descricatprod,(SELECT IFNULL(COUNT(idcategprod), 0) FROM productos WHERE idcategprod=productos_categorias.id AND idlugar=? AND suspendido=0 AND estado<>'B') AS nroProductos FROM productos_categorias HAVING nroProductos > 0", id).
		Scan(&categorias)

	// Busco promociones del lugar
	var promociones []Promociones_lugar
	db.Raw("SELECT id,vencimiento,titulo,descripcion,terminos,cuposdispon,rutaimg FROM promociones WHERE idlugar=? AND vencimiento>=? AND idtipopromo=3 AND cuposdispon>0 AND estado<>'B'", id, fecha).
		Scan(&promociones)

	// Busco valoracion del lugar segun cliente que consulta
	var valoracion = new(Modelos.Valoraciones)
	db.Raw("SELECT valoraciones.id,valoraciones.idlugar,lugares.nombrelugar,valoraciones.idusuario,usuarios_app.apellido,usuarios_app.nombres,usuarios_app.rutaimgusu,valoraciones.puntuacion,valoraciones.fecha,valoraciones.titulo,valoraciones.descripcion,valoraciones.fechamodif "+
		"FROM valoraciones "+
		"JOIN lugares ON lugares.idlugar=valoraciones.idlugar "+
		"JOIN usuarios_app ON usuarios_app.id=valoraciones.idusuario "+
		"WHERE valoraciones.idlugar=? AND valoraciones.idusuario=?", id, idUsuario).Scan(&valoracion)

	// Corte de control de productos segun categoria
	var productos_categoria []Productos_categoria
	if len(categorias) > 0 {
		for i := 0; i < len(categorias); i++ {
			var productosCategoria []Productos_lugar
			for j := 0; j < len(productos); j++ {
				if productos[j].Idcategprod == categorias[i].Id {
					productosCategoria = append(productosCategoria, productos[j])
				}
			}

			// productos que corresponden a una categoria
			if len(productosCategoria) > 0 {
				prod_cat := new(Productos_categoria)
				prod_cat.Id = categorias[i].Id
				prod_cat.Categoria = categorias[i].Descricatprod
				prod_cat.Productos = productosCategoria
				productos_categoria = append(productos_categoria, *prod_cat)
			}
		}

		// productos sin categorias
		var productosSinCategoria []Productos_lugar
		for i := 0; i < len(productos); i++ {
			if productos[i].Idcategprod == 0 {
				productosSinCategoria = append(productosSinCategoria, productos[i])
			}
		}
		if len(productosSinCategoria) > 0 {
			prod_sin_cat := new(Productos_categoria)
			prod_sin_cat.Id = 0
			prod_sin_cat.Categoria = ""
			prod_sin_cat.Productos = productosSinCategoria
			productos_categoria = append(productos_categoria, *prod_sin_cat)
		}

	} else { // Si no hay categorias cargadas creo una unica categoria con todos los productos

		if len(productos) > 0 {
			prod_cat := new(Productos_categoria)
			prod_cat.Id = 0
			prod_cat.Categoria = "Productos"
			prod_cat.Productos = productos
			productos_categoria = append(productos_categoria, *prod_cat)
		}

	}

	data := Data{Lugar: lugar, Horarios: horarios, Redes: redes, Productos_categoria: productos_categoria, Promociones: promociones, Valoraciones: valoracion}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	lugares := new(Modelos.Lugares)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(lugares); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	lugares.Fechaalta = time.Now() //utils.GetNow()
	lugares.Activo = 1
	if err := db.Omit("fechamodif", "fechaestado").Create(&lugares).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Genero qr para el lugar
	//contenido := fmt.Sprintf("%v", lugares.Idlugar)
	contenido := utils.Encrypt(fmt.Sprintf("%v", lugares.Idlugar), config.SecretKeyEncrypt)
	url := config.DirQrLugares + fmt.Sprintf("%v", lugares.Idlugar) + `.png`
	utils.CreateQr(contenido, 700, url)

	// Actualizo registro de lugar con el token
	if err := db.Exec("UPDATE lugares SET tokenlugar=?, qrasignado=1 WHERE idlugar=?", contenido, lugares.Idlugar).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{LugarAlta: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	lugares := new(Modelos.Lugares)
	if err := c.Bind(lugares); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("fechaalta").Save(&lugares).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{LugarAlta: lugares}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE lugares SET estado='B', fechaestado=? WHERE idlugar = ?", time.Now(), c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE lugares SET estado='', fechaestado=? WHERE idlugar = ?", time.Now(), c.Param("id")).Error; err != nil {
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

func MarcarCerrado(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE lugares SET activo=0 WHERE idlugar = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Lugar marcado como cerrado con éxito",
	})
}

func MarcarAbierto(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE lugares SET activo=1 WHERE idlugar = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Lugar marcado como abierto con éxito",
	})
}

func GenerarQr(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	// Genero token y qr
	//contenido := fmt.Sprintf("%v", id)
	contenido := utils.Encrypt(fmt.Sprintf("%v", id), config.SecretKeyEncrypt)
	url := config.DirQrLugares + fmt.Sprintf("%v", id) + `.png`
	utils.CreateQr(contenido, 700, url)

	// Actualizo registro de lugar con el token
	if err := db.Exec("UPDATE lugares SET tokenlugar=?, qrasignado=1 WHERE idlugar=?", contenido, id).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "QR generado con éxito",
	})
}

func ImprimirQr(c echo.Context) error {

	pdf := gofpdf.New("P", "mm", "A5", "")
	id := c.Param("id")

	fmt.Println("'id:'", id)

	x := 0.00
	y := 0.00

	pdf.SetHeaderFunc(func() {
	})

	pdf.AddPage()

	// Plantilla
	pdf.SetXY(x-20, y-20)
	var opt gofpdf.ImageOptions
	opt.ImageType = "jpg"
	pdf.ImageOptions(config.DirPlantillaQrLugares, 0, 0, 149, 210, false, opt, 0, "") // x, y, ancho, alto

	// QR
	//url := "https://www.vivircarlospaz.com/img/lugares/qr/" + id + ".png"
	//httpimg.Register(pdf, url, "")
	//pdf.Image(url, 53.5, 129, 48, 0, false, "", 0, "") // x, y, size
	url := config.DirQrLugares + id + ".png"
	var opt2 gofpdf.ImageOptions
	opt2.ImageType = "png"
	pdf.ImageOptions(url, 53.5, 129, 48, 0, false, opt2, 0, "") // x, y, ancho, alto

	path := fmt.Sprintf("%v", "/home/ubuntu/vivircarlospaz/backend/public/pdfqr/qr-"+c.Param("id")+".pdf")
	pdf.OutputFileAndClose(path)

	r := c.File(path)
	_ = os.Remove(path)

	return r
}
