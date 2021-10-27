package usuarios_sesiones

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
	Registros int                                  `json:"registros,omitempty"`
	Favoritos []Modelos.Usuarios_lugares_favoritos `json:"favoritos,omitempty"`
}

type Favorito struct {
	Idusuario uint `json:"idusuario"`
	Idlugar   uint `json:"idlugar"`
}

func (Favorito) TableName() string {
	return "usuarios_lugares_favoritos"
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	db = db.Select("usuarios_lugares_favoritos.id,usuarios_lugares_favoritos.idusuario,usuarios_lugares_favoritos.idlugar,lugares.nombrelugar,lugares.idrubro,rubros.descrirubro,lugares.idsubrubro,subrubros.dsubrubro,lugares.direccion,lugares.idlocalidad,localidades.nombrelocali,provincias.idprovincia,provincias.nombrepcia,(SELECT rutaimg FROM lugares_img WHERE idlugar=lugares.idlugar LIMIT 1) AS rutafoto,lugares.precdelivery,lugares.cpraminima")
	db = db.Joins("JOIN lugares ON lugares.idlugar=usuarios_lugares_favoritos.idlugar").
		Joins("JOIN rubros ON rubros.id=lugares.idrubro").
		Joins("JOIN subrubros ON subrubros.id=lugares.idsubrubro").
		Joins("JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia")

	// Controlo valores para filtro y paginacion que llegan de la url
	db = db.Where(" usuarios_lugares_favoritos.idusuario = ? ", c.Param("id"))

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
	var favoritos []Modelos.Usuarios_lugares_favoritos
	db.Offset(offset).Limit(limite).Find(&favoritos)
	db.Table("usuarios_lugares_favoritos").Count(&registros)
	data := Data{Registros: registros, Favoritos: favoritos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	favorito := new(Favorito)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(favorito); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia del lugar en favoritos
	favorito2 := new(Modelos.Usuarios_lugares_favoritos)
	db.Where("idusuario = ? AND idlugar = ?", favorito.Idusuario, favorito.Idlugar).First(&favorito2)
	if favorito2.Id > 0 {

		if err := db.Exec("DELETE FROM usuarios_lugares_favoritos WHERE idusuario = ? AND idlugar = ?", favorito.Idusuario, favorito.Idlugar).Error; err != nil {
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Eliminado de favoritos con éxito",
		})

	} else {

		// Inserta registro en la tabla
		if err := db.Create(&favorito).Error; err != nil {
			response := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Agregado a favoritos con éxito",
		})

	}
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM usuarios_lugares_favoritos WHERE id=?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Eliminado de favoritos con éxito",
	})
}
