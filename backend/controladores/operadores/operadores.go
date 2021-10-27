package operadores

import (
	"fmt"
	"net/http"
	"time"

	"../../config"
	"../../database"
	Modelos "../../modelos"
	"../../utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros  int                  `json:"registros,omitempty"`
	Operadores []Modelos.Operadores `json:"operadores,omitempty"`
	Operador   *Modelos.Operadores  `json:"operador,omitempty"`
	Token      string               `json:"token,omitempty"`
}

type Auth struct {
	Usuario string `json:"usuario,omitempty"`
	Clave   string `json:"clave,omitempty"`
}

type Clave struct {
	IDusuario string `json:"idusuario"`
	Clave     string `json:"clave"`
}

type Cambiar_clave struct {
	IDusuario string `json:"idusuario"`
	OldClave  string `json:"oldclave"`
	Clave     string `json:"clave"`
}

func Login(c echo.Context) error {
	db := database.GetDb()
	auth := new(Auth)
	if err := c.Bind(auth); err != nil {
		response := Respuesta{
			Status:  "error",
			Message: "invalid request body",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	operador := new(Modelos.Operadores)
	db.Where("nickoperador = BINARY ? and contrasena = BINARY ? and estado <> 'B'", auth.Usuario, auth.Clave).Preload("Nivel").First(&operador)

	if operador.ID > 0 {

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = fmt.Sprintf("%d", operador.ID)
		claims["contrasena"] = fmt.Sprintf("%s", auth.Clave)
		claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.SecretJwt))
		if err != nil {
			return err
		}
		data := Data{Operador: operador, Token: t}
		return c.JSON(http.StatusOK, Respuesta{Status: "success", Message: "Login correcto.", Data: data})

	} else {
		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "Datos de Login incorrectos."})
	}
}

func GetAutenticacionOperador(c echo.Context) error {
	db := database.GetDb()

	opr := c.Get("user").(*jwt.Token)
	claims := opr.Claims.(jwt.MapClaims)
	idoperador := claims["id"].(string)

	operador := new(Modelos.Operadores)
	db.Where("id = ? and estado <> 'B'", idoperador).Preload("Nivel").First(&operador)

	data := Data{Operador: operador}
	return c.JSON(http.StatusOK, Respuesta{Status: "success", Data: data})
}

func Contrasena(c echo.Context) error {
	db := database.GetDb()

	clave := new(Clave)
	if err := c.Bind(clave); err != nil {
		response := Respuesta{
			Status:  "error",
			Message: "invalid request body",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := db.Exec("UPDATE operadores_sistema SET contrasena = ? WHERE id = ?", clave.Clave, clave.IDusuario).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Contraseña asignada con éxito",
	})
}

func CambiarContrasena(c echo.Context) error {
	db := database.GetDb()

	cambiar_clave := new(Cambiar_clave)
	if err := c.Bind(cambiar_clave); err != nil {
		response := Respuesta{
			Status:  "error",
			Message: "invalid request body",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	operador := new(Modelos.Operadores)
	db.Where("id = ? AND contrasena = BINARY ?", cambiar_clave.IDusuario, cambiar_clave.OldClave).First(&operador)

	if operador.ID > 0 {

		if err := db.Exec("UPDATE operadores_sistema SET contrasena = ? WHERE id = ?", cambiar_clave.Clave, cambiar_clave.IDusuario).Error; err != nil {
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Contraseña asignada con éxito",
		})

	} else {

		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "Contraseña anterior incorrecta."})

	}
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" id like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nickoperador like ? ", "%"+c.QueryParam("query")+"%").
			Or(" apynombres like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("apynombres")
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
	var operadores []Modelos.Operadores
	db.Preload("Nivel").Offset(offset).Limit(limite).Find(&operadores)
	db.Table("operadores_sistema").Count(&registros)
	data := Data{Registros: registros, Operadores: operadores}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" apynombres LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("apynombres")

	// Ejecuto consulta
	var operadores []Modelos.Operadores
	db.Preload("Nivel").Find(&operadores)
	data := Data{Operadores: operadores}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetOperador(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	operadores := new(Modelos.Operadores)
	db.Preload("Nivel").Find(&operadores, id)

	data := Data{Operador: operadores}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	operadores := new(Modelos.Operadores)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(operadores); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Create(&operadores).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Operador: operadores}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	operadores := new(Modelos.Operadores)
	if err := c.Bind(operadores); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Save(&operadores).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Operador: operadores}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE operadores_sistema SET estado='B' WHERE id = ?", c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE operadores_sistema SET estado='' WHERE id = ?", c.Param("id")).Error; err != nil {
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
