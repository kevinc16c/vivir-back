package registro_usuarios_lugares

import (
	"fmt"
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
	Registros          int                                 `json:"registros,omitempty"`
	Registros_usuarios []Modelos.Registro_usuarios_lugares `json:"registros_usuarios,omitempty"`
	Registro_usuario   *Modelos.Registro_usuarios_lugares  `json:"registros_persona,omitempty"`
	Registro_alta      *Registro                           `json:"registros_alta,omitempty"`
}

type Registro struct {
	Idusuario  uint      `json:"idusuario"`
	Idlugar    uint      `json:"idlugar"`
	Tokenlugar string    `json:"tokenlugar"`
	Fechayhora time.Time `json:"fechayhora"`
}

func (Registro) TableName() string {
	return "registros_usuarios_lugares"
}

func Consultar(c echo.Context) error {
	db := database.GetDb()
	tipo := c.QueryParam("tipo")

	db = db.Select("registros_usuarios_lugares.id,registros_usuarios_lugares.idusuario,registros_usuarios_lugares.idlugar,lugares.nombrelugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.direccion,lugares.telefono,lugares.celular,lugares.e_mail,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,registros_usuarios_lugares.fechayhora")
	db = db.Joins("JOIN lugares ON lugares.idlugar=registros_usuarios_lugares.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia")

	// control de filtro
	var condicion string
	if tipo == "usuarios" {
		condicion = "registros_usuarios_lugares.idusuario = " + fmt.Sprintf("%v", c.Param("id")) + " AND fechayhora BETWEEN '" + c.QueryParam("desde") + "' AND '" + c.QueryParam("hasta") + "'"
	} else {
		condicion = "registros_usuarios_lugares.idlugar = " + fmt.Sprintf("%v", c.Param("id")) + " AND fechayhora BETWEEN '" + c.QueryParam("desde") + "' AND '" + c.QueryParam("hasta") + "'"
	}

	db = db.Where(condicion)

	// ORDENES ******************************************************************
	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fechayhora")
	}

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
	var registros_usuarios []Modelos.Registro_usuarios_lugares
	db.Preload("Usuario").Offset(offset).Limit(limite).Find(&registros_usuarios)
	db.Table("registros_usuarios_lugares").Count(&registros)
	data := Data{Registros: registros, Registros_usuarios: registros_usuarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	registro := new(Registro)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(registro); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	lugar := new(Modelos.Lugares)
	db.Raw("SELECT * FROM lugares WHERE lugares.tokenlugar = BINARY ?", registro.Tokenlugar).Scan(&lugar)
	if lugar.Idlugar > 0 {

		// Inserta registro en la tabla
		registro.Idlugar = lugar.Idlugar
		if err := db.Omit("tokenlugar").Create(&registro).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		// Preparo mensaje de retorno
		data := Data{Registro_alta: registro}
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "¡Registro realizado con éxito!",
			Data:    data,
		})
	} else {
		response := Respuesta{
			Status:  "error",
			Message: "No se pudo crear el registro. Intente nuevamente.",
		}
		return c.JSON(http.StatusBadRequest, response)
	}
}
