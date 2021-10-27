package usuarios

import (
	"fmt"
	"net/http"

	"../../config"
	"../../database"
	Modelos "../../modelos"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros int                `json:"registros,omitempty"`
	Usuarios  []Modelos.Usuarios `json:"usuarios,omitempty"`
	Usuario   *Modelos.Usuarios  `json:"usuario,omitempty"`
	Token     string             `json:"token,omitempty"`
}

type Auth struct {
	Email string `json:"email,omitempty"`
	Clave string `json:"clave,omitempty"`
}

type Crear_clave struct {
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

	usuario := new(Modelos.Usuarios)
	db.Where("correoelec = BINARY ? and contrasena = BINARY ? and asignarpass=0", auth.Email, auth.Clave).First(&usuario)

	if usuario.Id > 0 {

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = fmt.Sprintf("%d", usuario.Id)
		claims["correoelec"] = fmt.Sprintf("%s", usuario.Correoelec)
		claims["exp"] = config.ExpJWTUser

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.SecretJwt))
		if err != nil {
			return err
		}
		usuario.Contrasena = "" // blanqueo contraseña antes de devolver datos
		data := Data{Usuario: usuario, Token: t}
		return c.JSON(http.StatusOK, Respuesta{Status: "success", Message: "Login correcto.", Data: data})

	} else {
		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "Datos de Login incorrectos."})
	}
}

func CrearContrasena(c echo.Context) error {
	db := database.GetDb()

	crear_clave := new(Crear_clave)
	if err := c.Bind(crear_clave); err != nil {
		response := Respuesta{
			Status:  "error",
			Message: "invalid request body",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := db.Exec("UPDATE usuarios_app SET contrasena = ?, asignarpass=0 WHERE id = ?", crear_clave.Clave, crear_clave.IDusuario).Error; err != nil {
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

	usuario := new(Modelos.Usuarios)
	db.Where("id = ? AND contrasena = BINARY ?", cambiar_clave.IDusuario, cambiar_clave.OldClave).First(&usuario)

	if usuario.Id > 0 {

		if err := db.Exec("UPDATE usuarios_app SET contrasena = ? WHERE id = ?", cambiar_clave.Clave, cambiar_clave.IDusuario).Error; err != nil {
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

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("id,apellido,nombres,numedocume,celular,correoelec")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" apellido LIKE ? ", "%"+c.QueryParam("query")+"%").
			Or(" nombres like ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("apellido, nombres ASC")

	// Ejecuto consulta
	var usuarios []Modelos.Usuarios
	db.Find(&usuarios)
	data := Data{Usuarios: usuarios}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	usuario := new(Modelos.Usuarios)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(usuario); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia del email
	usuario2 := new(Modelos.Usuarios)
	db.Where("correoelec = BINARY ?", usuario.Correoelec).First(&usuario2)

	if usuario2.Id > 0 {
		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "El e-mail ingresado ya está en uso."})
	}

	// Inserta registro en la tabla
	usuario.Asignarpass = 0
	if err := db.Create(&usuario).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo JWT para incluir en la respuesta
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = fmt.Sprintf("%d", usuario.Id)
	claims["correoelec"] = fmt.Sprintf("%s", usuario.Correoelec)
	claims["exp"] = config.ExpJWTUser

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.SecretJwt))
	if err != nil {
		return err
	}

	// Preparo mensaje de retorno
	usuario.Contrasena = "" // blanqueo contraseña antes de devolver datos
	data := Data{Usuario: usuario, Token: t}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	usuario := new(Modelos.Usuarios)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(usuario); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	if err := db.Exec("UPDATE usuarios_app SET numedocume=?, celular=? WHERE id = ?", usuario.Numedocume, usuario.Celular, usuario.Id).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	usuario.Contrasena = "" // blanqueo contraseña antes de devolver datos
	data := Data{Usuario: usuario}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func CuentaRedesSociales(c echo.Context) error {
	db := database.GetDb()

	usuario := new(Modelos.Usuarios)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(usuario); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Controlo la existencia del email
	usuario2 := new(Modelos.Usuarios)
	db.Where("correoelec = BINARY ?", usuario.Correoelec).First(&usuario2)

	// Si existe el e-mail, hago un login
	if usuario2.Id > 0 {

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = fmt.Sprintf("%d", usuario2.Id)
		claims["correoelec"] = fmt.Sprintf("%s", usuario2.Correoelec)
		claims["exp"] = config.ExpJWTUser

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.SecretJwt))
		if err != nil {
			return err
		}
		usuario2.Contrasena = "" // blanqueo contraseña antes de devolver datos
		data := Data{Usuario: usuario2, Token: t}
		return c.JSON(http.StatusOK, Respuesta{Status: "success", Message: "Login correcto.", Data: data})
	}

	// Si no existe el e-mail, Inserta registro en la tabla
	usuario.Asignarpass = 1
	if err := db.Create(&usuario).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo JWT para incluir en la respuesta
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = fmt.Sprintf("%d", usuario.Id)
	claims["correoelec"] = fmt.Sprintf("%s", usuario.Correoelec)
	claims["exp"] = config.ExpJWTUser

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.SecretJwt))
	if err != nil {
		return err
	}

	// Preparo mensaje de retorno
	usuario.Contrasena = "" // blanqueo contraseña antes de devolver datos
	data := Data{Usuario: usuario, Token: t}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
