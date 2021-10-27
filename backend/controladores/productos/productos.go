package productos

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
	Registros    int                         `json:"registros,omitempty"`
	Productos    []ProductosLista            `json:"productos,omitempty"`
	Producto     *ProductosLista             `json:"producto,omitempty"`
	ProductoAlta *Modelos.Productos          `json:"producto_alta,omitempty"`
	Imagenes     []Modelos.Productos_img     `json:"imagenes,omitempty"`
	Sabores      []Modelos.Productos_insumos `json:"sabores,omitempty"`
}

type ProductosLista struct {
	Id            uint       `json:"id" gorm:"primary_key"`
	Idlugar       uint       `json:"idlugar"`
	Nombrelugar   string     `json:"nombrelugar"`
	Codintprod    string     `json:"codintprod"`
	Idcategprod   uint       `json:"idcategprod"`
	Descricatprod string     `json:"descricatprod"`
	Descriprod    string     `json:"descriprod"`
	Desextprod    string     `json:"desextprod"`
	Prunitprod    float64    `json:"prunitprod"`
	Aliivaprod    float64    `json:"Aliivaprod"`
	Existencia    uint       `json:"existencia"`
	Agregados     uint       `json:"agregados"`
	Ocupdesde     *time.Time `json:"ocupdesde"`
	Ocuphasta     *time.Time `json:"ocuphasta"`
	Suspendido    uint       `json:"suspendido"`
	Controlstock  uint       `json:"controlstock"`
	Estado        string     `json:"estado"`
	Fechabaja     *time.Time `json:"fechabaja"`
}

func (ProductosLista) TableName() string {
	return "productos"
}

func PaginacionProductosPropietario(c echo.Context) error {
	db := database.GetDb()
	idpropietario := c.Param("id")

	db = db.Select("productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja")
	db = db.Joins("JOIN lugares ON lugares.idlugar=productos.idlugar").
		Joins("JOIN productos_categorias ON productos_categorias.id=productos.idcategprod")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" lugares.idpropietario = ? AND (productos.id like ? OR descriprod like ? OR desextprod like ? OR descricatprod like ? OR nombrelugar like ?)", idpropietario, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" lugares.idpropietario = ?", idpropietario)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("descriprod")
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
	var productos []ProductosLista
	db.Offset(offset).Limit(limite).Find(&productos)
	db.Table("productos").Count(&registros)
	data := Data{Registros: registros, Productos: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PaginacionProductosLugar(c echo.Context) error {
	db := database.GetDb()
	idlugar := c.Param("id")

	db = db.Select("productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.agregados,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja")
	db = db.Joins("JOIN lugares ON lugares.idlugar=productos.idlugar").
		Joins("LEFT JOIN productos_categorias ON productos_categorias.id=productos.idcategprod")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" productos.idlugar = ? AND (productos.id like ? OR descriprod like ? OR desextprod like ? OR descricatprod like ?)", idlugar, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" productos.idlugar = ?", idlugar)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("descriprod")
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
	var productos []ProductosLista
	db.Offset(offset).Limit(limite).Find(&productos)
	db.Table("productos").Count(&registros)
	data := Data{Registros: registros, Productos: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func ListaProductosPropietario(c echo.Context) error {
	db := database.GetDb()
	idpropietario := c.Param("id")

	// Armo select
	db = db.Select("productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.agregados,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja")
	db = db.Joins("JOIN lugares ON lugares.idlugar=productos.idlugar").
		Joins("JOIN productos_categorias ON productos_categorias.id=productos.idcategprod")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" lugares.idpropietario = ? AND idescriprod LIKE ? AND productos.estado<>'B'", idpropietario, "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" lugares.idpropietario = ? AND productos.estado<>'B'", idpropietario)
	}

	db = db.Order("descriprod")

	// Ejecuto consulta
	var productos []ProductosLista
	db.Find(&productos)
	data := Data{Productos: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func ListaProductosLugar(c echo.Context) error {
	db := database.GetDb()
	idlugar := c.Param("id")

	// Armo select
	db = db.Select("productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.agregados,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja")
	db = db.Joins("JOIN lugares ON lugares.idlugar=productos.idlugar").
		Joins("LEFT JOIN productos_categorias ON productos_categorias.id=productos.idcategprod")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" productos.idlugar = ? AND descriprod LIKE ? AND productos.suspendido=0 AND productos.estado<>'B'", idlugar, "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" productos.idlugar = ? AND productos.suspendido=0 AND productos.estado<>'B'", idlugar)
	}

	db = db.Order("descriprod")

	// Ejecuto consulta
	var productos []ProductosLista
	db.Find(&productos)
	data := Data{Productos: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetProducto(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	productos := new(ProductosLista)
	//db.Preload("Lugares.Personas_humanas_juridicas.Localidad").Preload("Categoria").Find(&productos, id)

	db = db.Select("productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.agregados,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja")
	db = db.Joins("JOIN lugares ON lugares.idlugar=productos.idlugar").
		Joins("LEFT JOIN productos_categorias ON productos_categorias.id=productos.idcategprod")
	db = db.Where(" productos.id = ?", id)

	db.Find(&productos)

	data := Data{Producto: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetProductoDetalle(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	// Producto
	producto := new(ProductosLista)
	db.Raw("SELECT productos.id,productos.idlugar,lugares.nombrelugar,productos.codintprod,productos.idcategprod,productos_categorias.descricatprod,productos.descriprod,productos.desextprod,productos.prunitprod,productos.aliivaprod,productos.existencia,productos.agregados,productos.suspendido,productos.controlstock,productos.estado,productos.fechabaja "+
		"FROM productos "+
		"JOIN lugares ON lugares.idlugar=productos.idlugar "+
		"LEFT JOIN productos_categorias ON productos_categorias.id=productos.idcategprod "+
		"WHERE productos.id = ?", id).Scan(&producto)

	// Imagenes
	var imagenes []Modelos.Productos_img
	db.Raw("SELECT id,idproducto,rutaimgprod FROM productos_img WHERE idproducto = ?", id).Scan(&imagenes)

	// insumos
	var sabores []Modelos.Productos_insumos
	db.Raw("SELECT id,idproducto,idinsumo,dinsuprodu FROM productos_insumos WHERE idproducto = ?", id).Scan(&sabores)

	data := Data{Producto: producto, Imagenes: imagenes, Sabores: sabores}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	productos := new(Modelos.Productos)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(productos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&productos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{ProductoAlta: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	productos := new(Modelos.Productos)
	if err := c.Bind(productos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&productos).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{ProductoAlta: productos}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE productos SET estado='B', fechabaja=? WHERE id = ?", time.Now(), c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE productos SET estado='', fechabaja=null WHERE id = ?", c.Param("id")).Error; err != nil {
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

func Agotado(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE productos SET suspendido=1 WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro marcado como agotado con éxito",
	})
}

func EnStock(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE productos SET suspendido=0 WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Registro marcado en stock con éxito",
	})
}
