package promociones

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"../../config"
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
	Registros   int                        `json:"registros,omitempty"`
	Promociones []Modelos.Promociones      `json:"promociones,omitempty"`
	Promocion   *Modelos.Promociones       `json:"promocion,omitempty"`
	Canje       *Modelos.Promociones_canje `json:"canje,omitempty"`
}

type PromocionesAlta struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	Idlugar     uint      `json:"idlugar"`
	Vencimiento time.Time `json:"vencimiento"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Terminos    string    `json:"terminos"`
	Canticupos  uint      `json:"canticupos"`
	Cuposdispon uint      `json:"cuposdispon"`
	Imagen      string    `json:"imagen,omitempty" gorm:"-"`
	Idtipopromo uint      `json:"idtipopromo"`
	Rutaimg     string    `json:"rutaimg,omitempty"`
	Fechaalta   time.Time `json:"fechaalta,omitempty"`
	Fechamodif  time.Time `json:"fechamodif,omitempty"`
	Estado      string    `json:"estado"`
}

func (PromocionesAlta) TableName() string {
	return "promociones"
}

type Imagen struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Rutaimgold string `json:"rutaimagenold"`
	Imagen     string `json:"imagen,omitempty"`
}

type Canje_promocion struct {
	Id          uint      `json:"id"`
	Idpromocion uint      `json:"idpromocion"`
	Idusuario   uint      `json:"idusuario"`
	Tokenlugar  string    `json:"tokenlugar"`
	Fechacarga  time.Time `json:"fechacarga"`
}

func (Canje_promocion) TableName() string {
	return "promociones_canje"
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("promociones.id,promociones.idlugar,lugares.nombrelugar,lugares.direccion,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,promociones.vencimiento,promociones.titulo,promociones.descripcion,promociones.terminos,promociones.fechamodif,promociones.canticupos,promociones.cuposdispon,promociones.rutaimg,promociones.idtipopromo,promociones.estado")
	db = db.Joins("JOIN lugares ON lugares.idlugar=promociones.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idtipopromo<>3 AND promociones.id like ? OR titulo like ? OR descripcion like ? OR nombrelugar like ? OR descrirubro like ? OR dsubrubro like ?", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idtipopromo<>3")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("id")
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
	var promociones []Modelos.Promociones
	db.Offset(offset).Limit(limite).Find(&promociones)
	db.Table("promociones").Count(&registros)
	data := Data{Registros: registros, Promociones: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PaginacionPromocionesLugar(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("promociones.id,promociones.idlugar,lugares.nombrelugar,lugares.direccion,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,promociones.vencimiento,promociones.titulo,promociones.descripcion,promociones.terminos,promociones.fechamodif,promociones.canticupos,promociones.cuposdispon,promociones.rutaimg,promociones.idtipopromo,promociones.estado")
	db = db.Joins("JOIN lugares ON lugares.idlugar=promociones.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idtipopromo=3 AND promociones.idlugar=? AND (promociones.id like ? OR titulo like ? OR descripcion like ? OR nombrelugar like ? OR descrirubro like ? OR dsubrubro like ?)", c.Param("id"), "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idtipopromo=3 AND promociones.idlugar=?", c.Param("id"))
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("id")
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
	var promociones []Modelos.Promociones
	db.Offset(offset).Limit(limite).Find(&promociones)
	db.Table("promociones").Count(&registros)
	data := Data{Registros: registros, Promociones: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPromocionesExclusivas(c echo.Context) error {
	db := database.GetDb()
	fecha := time.Now() //utils.GetNow()

	db = db.Select("promociones.id,promociones.idlugar,lugares.nombrelugar,lugares.direccion,lugares.latitud,lugares.longitud,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,promociones.vencimiento,promociones.titulo,promociones.descripcion,promociones.terminos,promociones.fechamodif,promociones.canticupos,promociones.cuposdispon,promociones.rutaimg,promociones.visitas,promociones.estado")
	db = db.Joins("JOIN lugares ON lugares.idlugar=promociones.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro")

	//db = db.Where(" idtipopromo=1 AND cuposdispon>0 AND promociones.vencimiento>=? AND promociones.estado<>'B'", fecha)

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idtipopromo=1 AND cuposdispon>0 AND promociones.vencimiento>=? AND promociones.estado<>'B' AND (nombrelugar LIKE ? OR descrirubro LIKE ? OR dsubrubro LIKE ? OR descripcion LIKE ?)", fecha, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idtipopromo=1 AND cuposdispon>0 AND promociones.vencimiento>=? AND promociones.estado<>'B'", fecha)
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
	var promociones []Modelos.Promociones
	db.Offset(offset).Limit(limite).Find(&promociones)
	db.Table("promociones").Count(&registros)
	data := Data{Registros: registros, Promociones: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPromocion(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	promociones := new(Modelos.Promociones)
	db = db.Select("promociones.id,promociones.idlugar,lugares.nombrelugar,lugares.direccion,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,promociones.vencimiento,promociones.titulo,promociones.descripcion,promociones.terminos,promociones.fechamodif,promociones.canticupos,promociones.cuposdispon,promociones.rutaimg,promociones.idtipopromo,promociones.estado")
	db = db.Joins("JOIN lugares ON lugares.idlugar=promociones.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Find(&promociones, id)

	data := Data{Promocion: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	promociones := new(PromocionesAlta)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(promociones); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Subo la imagen
	fecha := time.Now() //utils.GetNow()
	strFecha := fmt.Sprintf("%v", fecha.Format("20060102150405"))
	nombreImagen := strFecha
	urlImagen := ""

	formato := utils.GetFormatoImagen(promociones.Imagen)
	if formato == "png" {
		nombreImagen = nombreImagen + ".png"
		urlImagen = config.DirImgPromociones + nombreImagen

		err := utils.SavePng(promociones.Imagen, nombreImagen, urlImagen)
		if err != "ok" {
			response := Respuesta{
				Status:  "error",
				Message: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	} else {
		nombreImagen = nombreImagen + ".jpg"
		urlImagen = config.DirImgPromociones + nombreImagen

		err := utils.SaveJpg(promociones.Imagen, nombreImagen, urlImagen)
		if err != "ok" {
			response := Respuesta{
				Status:  "error",
				Message: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	// Inserta registro en la tabla
	promociones.Rutaimg = config.UrlImgPromociones + nombreImagen
	promociones.Estado = ""
	promociones.Fechaalta = time.Now() //utils.GetNow()
	if err := db.Omit("fechamodif").Create(&promociones).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	//data := Data{Promocion: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		//Data:   data,
		Message: "Promoción guardada con éxito",
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	promociones := new(PromocionesAlta)
	if err := c.Bind(promociones); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	promociones.Fechamodif = time.Now() //utils.GetNow()
	if err := db.Omit("fechaalta", "rutaimg", "estado").Save(&promociones).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	//data := Data{Pro: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		//Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE promociones SET estado='B' WHERE id = ?", c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE promociones SET estado='' WHERE id = ?", c.Param("id")).Error; err != nil {
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

func CambiarImagen(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	imagen := new(Imagen)
	if err := c.Bind(imagen); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Subo la imagen nueva
	fecha := time.Now() //utils.GetNow()
	strFecha := fmt.Sprintf("%v", fecha.Format("20060102150405"))
	nombreImagen := strFecha + ".png"
	urlImagen := config.DirImgPromociones + nombreImagen

	formato := utils.GetFormatoImagen(imagen.Imagen)
	if formato == "png" {
		nombreImagen = nombreImagen + ".png"
		urlImagen = config.DirImgPromociones + nombreImagen

		err := utils.SavePng(imagen.Imagen, nombreImagen, urlImagen)
		if err != "ok" {
			response := Respuesta{
				Status:  "error",
				Message: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	} else {
		nombreImagen = nombreImagen + ".jpg"
		urlImagen = config.DirImgPromociones + nombreImagen

		err := utils.SaveJpg(imagen.Imagen, nombreImagen, urlImagen)
		if err != "ok" {
			response := Respuesta{
				Status:  "error",
				Message: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	// Actualiza el registro
	rutaImagen := config.UrlImgPromociones + nombreImagen
	if err := db.Exec("UPDATE promociones SET rutaimg = ? WHERE id = ?", rutaImagen, imagen.Id).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// borro imagen anterior
	urlImagenOld := config.DirProyecto + imagen.Rutaimgold
	os.Remove(urlImagenOld)

	// Preparo mensaje de retorno
	//data := Data{Pro: promociones}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		//Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Canje(c echo.Context) error {
	db := database.GetDb()

	canje := new(Canje_promocion)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(canje); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo si existe una promo asociada al token
	promocion := new(Modelos.Promociones)
	db.Raw("SELECT promociones.*,lugares.tokenlugar FROM promociones JOIN lugares ON lugares.idlugar=promociones.idlugar WHERE promociones.id = ? and lugares.tokenlugar = BINARY ?", canje.Idpromocion, canje.Tokenlugar).
		Scan(&promocion)
	if promocion.Id > 0 {

		// transaccion
		tr := db.Begin()

		if err := db.Exec("UPDATE promociones SET cuposdispon=cuposdispon-1 WHERE id=?", canje.Idpromocion).Error; err != nil {
			tr.Rollback()
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		canje.Fechacarga = time.Now() //utils.GetNow()
		if err := db.Omit("Tokenlugar").Create(&canje).Error; err != nil {
			tr.Rollback()
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		// Commit
		tr.Commit()

		// Consulto canje
		promociones_canje := new(Modelos.Promociones_canje)
		db = db.Select("promociones_canje.id,promociones_canje.idpromocion,promociones.titulo,promociones.descripcion,promociones.rutaimg AS imagen,promociones_canje.fechacarga AS fecha")
		db = db.Joins("JOIN promociones ON promociones.id=promociones_canje.idpromocion")
		db = db.Where(" promociones_canje.id=?", canje.Id)
		db.Find(&promociones_canje)

		data := Data{Canje: promociones_canje}
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Data:    data,
			Message: "Canje realizado con éxito",
		})

	} else {

		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "El QR escaneado no corresponde a la promoción seleccionada."})

	}
}
