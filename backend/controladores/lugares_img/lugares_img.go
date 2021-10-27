package lugares_img

import (
	"fmt"
	"net/http"
	"time"

	config "../../config"
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
	Registros int                   `json:"registros,omitempty"`
	Imgenes   []Modelos.Lugares_img `json:"imagenes,omitempty"`
	Imagen    *Modelos.Lugares_img  `json:"imagen,omitempty"`
}

type Imagen struct {
	Idlugar   uint   `json:"idlugar"`
	Tituloimg string `json:"tituloimg"`
	Descriimg string `json:"descriimg"`
	Imagen    string `json:"imagen"`
	Rutaimg   string `json:"rutaimg,omitempty"`
}

func (Imagen) TableName() string {
	return "lugares_img"
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" idlugar = ? and (tituloimg like ? or descriimg like ?)", (c.Param("id")), "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" idlugar = ? ", c.Param("id"))
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("id")
	}

	// Preparo paginacion
	var pagina uint = 1
	var limite uint = 100
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
	var imagenes []Modelos.Lugares_img
	db.Offset(offset).Limit(limite).Find(&imagenes)
	db.Table("lugares_img").Count(&registros)
	data := Data{Registros: registros, Imgenes: imagenes}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetImagen(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	imagenes := new(Modelos.Lugares_img)
	db.Find(&imagenes, id)

	data := Data{Imagen: imagenes}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	imagen := new(Imagen)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(imagen); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Subo la imagen
	fecha := time.Now() //utils.GetNow()
	strFecha := fmt.Sprintf("%v", fecha.Format("20060102150405"))
	nombreImagen := fmt.Sprintf("%v", imagen.Idlugar) + "_" + strFecha
	urlImagen := ""

	formato := utils.GetFormatoImagen(imagen.Imagen)
	if formato == "png" {
		nombreImagen = nombreImagen + ".png"
		urlImagen = config.DirImgLugares + nombreImagen

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
		urlImagen = config.DirImgLugares + nombreImagen

		err := utils.SaveJpg(imagen.Imagen, nombreImagen, urlImagen)
		if err != "ok" {
			response := Respuesta{
				Status:  "error",
				Message: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	// Inserta registro en la tabla
	imagen.Rutaimg = config.UrlImgLugares + nombreImagen
	if err := db.Omit("imagen").Create(&imagen).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	//data := Data{Imagen: imagen}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		//Data:   data,
		Message: "Imágen guardada con éxito",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM lugares_img WHERE id = ?", c.Param("id")).Error; err != nil {
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
