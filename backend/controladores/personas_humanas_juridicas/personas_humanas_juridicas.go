package personas_humanas_juridicas

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
	"gopkg.in/gomail.v2"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros int                                 `json:"registros,omitempty"`
	Personas  []Propietarios                      `json:"propietarios,omitempty"`
	Persona   *Modelos.Personas_humanas_juridicas `json:"propietario,omitempty"`
	Token     string                              `json:"token,omitempty"`
}

type Propietarios struct {
	Id           uint   `json:"id" gorm:"primary_key"`
	Razonsocial  string `json:"razonsocial"`
	Nofantasia   string `json:"nofantasia"`
	Idcondiva    uint   `json:"idcondiva"`
	Descriciva   string `json:"descriciva"`
	Numerocuit   string `json:"numerocuit"`
	Direccion    string `json:"direccion"`
	Idlocalidad  uint   `json:"idlocalidad"`
	Nombrelocali string `json:"nombrelocali"`
	Idprovincia  uint   `json:"idprovincia"`
	Nombrepcia   string `json:"nombrepcia"`
	Idpais       uint   `json:"idpais"`
	Nombrepais   string `json:"nombrepais"`
	Telefono     string `json:"telefono"`
	Telefono2    string `json:"telefono2"`
	Celular1     string `json:"celular1"`
	Celular2     string `json:"celular2"`
	Celular3     string `json:"celular3"`
	Email        string `json:"email"`
	Estado       string `json:"estado"`
	Cambiarpass  uint   `json:"cambiarpass"`
}

func (Propietarios) TableName() string {
	return "personas_humanas_juridicas"
}

type Auth struct {
	Email string `json:"email,omitempty"`
	Clave string `json:"clave,omitempty"`
}

type Clave struct {
	Idpersona string `json:"id"`
	Clave     string `json:"clave"`
}

type Cambiar_clave struct {
	Idpersona string `json:"id"`
	OldClave  string `json:"oldclave"`
	Clave     string `json:"clave"`
}

type Mail struct {
	Mail string `json:"mail"`
}

type Propietario struct {
	Contrasena string `json:"mail"`
}

func (Propietario) TableName() string {
	return "personas_humanas_juridicas"
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

	persona := new(Modelos.Personas_humanas_juridicas)
	db.Where("email = BINARY ? and contrasena = BINARY ? and estado <> 'B'", auth.Email, auth.Clave).First(&persona)

	if persona.Id > 0 {

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = fmt.Sprintf("%d", persona.Id)
		claims["contrasena"] = fmt.Sprintf("%s", auth.Clave)
		claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.SecretPanelJwt))
		if err != nil {
			return err
		}
		data := Data{Persona: persona, Token: t}
		return c.JSON(http.StatusOK, Respuesta{Status: "success", Message: "Login correcto.", Data: data})

	} else {
		return c.JSON(http.StatusOK, Respuesta{Status: "error", Message: "Datos de Login incorrectos."})
	}
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

	if err := db.Exec("UPDATE personas_humanas_juridicas SET contrasena = ?, cambiarpass=1 WHERE id = ?", clave.Clave, clave.Idpersona).Error; err != nil {
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

	persona := new(Modelos.Personas_humanas_juridicas)
	db.Where("id = ? AND contrasena = BINARY ?", cambiar_clave.Idpersona, cambiar_clave.OldClave).First(&persona)

	if persona.Id > 0 {

		if err := db.Exec("UPDATE personas_humanas_juridicas SET contrasena = ?, cambiarpass=0 WHERE id = ?", cambiar_clave.Clave, cambiar_clave.Idpersona).Error; err != nil {
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

func GetAutenticacionPropietario(c echo.Context) error {
	db := database.GetDb()

	opr := c.Get("user").(*jwt.Token)
	claims := opr.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	propietarios := new(Modelos.Personas_humanas_juridicas)
	db.Where("id = ? and estado <> 'B'", id).Preload("Condicioniva").Preload("Localidad").First(&propietarios)

	data := Data{Persona: propietarios}
	return c.JSON(http.StatusOK, Respuesta{Status: "success", Data: data})
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("personas_humanas_juridicas.id,personas_humanas_juridicas.razonsocial,personas_humanas_juridicas.nofantasia,personas_humanas_juridicas.idcondiva,condicion_iva.descriciva,personas_humanas_juridicas.numerocuit,personas_humanas_juridicas.direccion,personas_humanas_juridicas.idlocalidad,localidades.nombrelocali,localidades.idprovincia,provincias.nombrepcia,provincias.idpais,paises.nombrepais,personas_humanas_juridicas.telefono,personas_humanas_juridicas.telefono2,personas_humanas_juridicas.celular1,personas_humanas_juridicas.celular2,personas_humanas_juridicas.celular3,personas_humanas_juridicas.email,personas_humanas_juridicas.estado")
	db = db.Joins("JOIN condicion_iva ON condicion_iva.codigociva=personas_humanas_juridicas.idcondiva").
		Joins("JOIN localidades ON localidades.idlocalidad=personas_humanas_juridicas.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais")

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" id like ? ", "%"+c.QueryParam("query")+"%").
			Or(" razonsocial like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nofantasia like ? ", "%"+c.QueryParam("query")+"%").
			Or(" usuario like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("razonsocial")
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
	var personas []Propietarios
	db.Offset(offset).Limit(limite).Find(&personas)
	db.Table("personas_humanas_juridicas").Count(&registros)
	data := Data{Registros: registros, Personas: personas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	// Armo select
	db = db.Select("personas_humanas_juridicas.id,personas_humanas_juridicas.razonsocial,personas_humanas_juridicas.nofantasia,personas_humanas_juridicas.idcondiva,condicion_iva.descriciva,personas_humanas_juridicas.numerocuit,personas_humanas_juridicas.direccion,personas_humanas_juridicas.idlocalidad,localidades.nombrelocali,localidades.idprovincia,provincias.nombrepcia,provincias.idpais,paises.nombrepais,personas_humanas_juridicas.telefono,personas_humanas_juridicas.telefono2,personas_humanas_juridicas.celular1,personas_humanas_juridicas.celular2,personas_humanas_juridicas.celular3,personas_humanas_juridicas.email,personas_humanas_juridicas.estado")
	db = db.Joins("JOIN condicion_iva ON condicion_iva.codigociva=personas_humanas_juridicas.idcondiva").
		Joins("JOIN localidades ON localidades.idlocalidad=personas_humanas_juridicas.idlocalidad").
		Joins("JOIN provincias ON provincias.idprovincia=localidades.idprovincia").
		Joins("JOIN paises ON paises.idpais=provincias.idpais")

	// Controlo parametro query recibido
	if c.QueryParam("query") != "" {
		db = db.Where(" razonsocial LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	db = db.Order("razonsocial")

	// Ejecuto consulta
	var personas []Propietarios
	db.Find(&personas)
	data := Data{Personas: personas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPersona(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	personas := new(Modelos.Personas_humanas_juridicas)
	db.Preload("Condicioniva").Preload("Localidad.Provincia.Pais").Find(&personas, id)

	data := Data{Persona: personas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	personas := new(Modelos.Personas_humanas_juridicas)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(personas); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	if err := db.Omit("fechaestado").Create(&personas).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	data := Data{Persona: personas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Modificar(c echo.Context) error {
	db := database.GetDb()

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	personas := new(Modelos.Personas_humanas_juridicas)
	if err := c.Bind(personas); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body ",
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza el registro
	if err := db.Omit("fechaestado").Save(&personas).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Preparo mensaje de retorno
	data := Data{Persona: personas}
	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Data:    data,
		Message: "Los datos se actualizaron correctamente. ",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE personas_humanas_juridicas SET estado='B', fechaestado=? WHERE id = ?", time.Now(), c.Param("id")).Error; err != nil {
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

	if err := db.Exec("UPDATE personas_humanas_juridicas SET estado='', fechaestado=? WHERE id = ?", time.Now(), c.Param("id")).Error; err != nil {
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

// func EnviarClaveEmail(c echo.Context) error {
// 	db := database.GetDb()
// 	registro := new(Mail)
// 	id := c.Param("id")

// 	if err := c.Bind(registro); err != nil {
// 		response := Respuesta{
// 			Status:  "error",
// 			Message: err.Error(),
// 		}
// 		return c.JSON(http.StatusBadRequest, response)
// 	}

// 	propietario := new(Propietario)
// 	db.Where("id = ?", id).First(&propietario)

// 	if propietario.Contrasena != "" {

// 		from := "girqta@gmail.com"
// 		pass := "ytumerecibesasi"
// 		to := registro.Mail
// 		msg := ""
// 		msg = "From: " + from + "\n" +
// 			"To: " + to + "\n" +
// 			"Estimado propietario, informamos que los datos para ingresar al panel de Vivir Carlos Paz es los siguientes:" + "\n" +
// 			"\n" +
// 			"Usuario: " + registro.Mail + "\n" +
// 			"Contraseña: " + propietario.Contrasena + "\n" +
// 			"\n" +
// 			"Gracias." + "\n"

// 		err := smtp.SendMail("smtp.gmail.com:465",
// 			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
// 			from, []string{to}, []byte(msg))

// 		if err != nil {
// 			return c.JSON(http.StatusOK, Respuesta{
// 				Status:  "errord",
// 				Message: err.Error(),
// 			})
// 		}

// 		return c.JSON(http.StatusOK, Respuesta{
// 			Status:  "success",
// 			Message: "Ok",
// 		})

// 	} else {

// 		return c.JSON(http.StatusOK, Respuesta{
// 			Status:  "error",
// 			Message: "No se encontró propietario",
// 		})

// 	}
// }

func EnviarClaveEmail1(c echo.Context) error {
	db := database.GetDb()
	registro := new(Mail)
	id := c.Param("id")

	if err := c.Bind(registro); err != nil {
		response := Respuesta{
			Status:  "error1",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	fmt.Println(registro)
	propietario := new(Propietario)
	db.Where("id = ?", id).First(&propietario)

	if propietario.Contrasena != "" {
		m := gomail.NewMessage()
		m.SetHeader("From", "girqta@gmail.com")
		m.SetHeader("To", "kevinc16c@gmail.com")
		m.SetHeader("Subject", "Datos de acceso Vivir Carlos Paz")
		m.SetBody("text/html", `
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html xmlns="http://www.w3.org/1999/xhtml">
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<title>Demystifying Email Design</title>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
			<style>
				.city {
				background-color: #ebebeb;
				color: white;
				padding: 10px;
				};
				.fonde {
				padding: 10px;
				background: black;
				}
			</style>
		</head>
		<body>
			<div class="container" style="font-size:15px; background:#f1f1f1">
				<div class="row">
					<p class="text-center" style="font-size:20px">Estimado propietario, informamos que los datos para ingresar al panel de Vivir Carlos Paz son los siguientes:</p>
					<p>Usuario: <b>`+fmt.Sprint(registro.Mail)+`</b></p> 
				</div>
                <div class="row">
					<p>Contraseña: <b>`+propietario.Contrasena+`</b></p> 
				</div>
                <div class="row">
					<p>Gracias.</p>
				</div>
			</div>
			<br/>
			<img src="https://www.vivircarlospaz.com/img/static/vcp-logogv.png"/>
		</body>
		</html>`)

		// Send the email to Bob
		d := gomail.NewPlainDialer("smtp.gmail.com", 587, "girqta@gmail.com", "ytumerecibesasi")
		if err := d.DialAndSend(m); err != nil {
			return c.JSON(http.StatusOK, Respuesta{
				Status:  "errord",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Ok",
		})
	} else {

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "error",
			Message: "No se encontró propietario",
		})

	}
}

func EnviarClaveEmail(c echo.Context) error {
	db := database.GetDb()
	registro := new(Mail)
	id := c.Param("id")

	if err := c.Bind(registro); err != nil {
		response := Respuesta{
			Status:  "error1",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	fmt.Println(registro)
	propietario := new(Propietario)
	db.Where("id = ?", id).First(&propietario)

	if propietario.Contrasena != "" {
		m := gomail.NewMessage()
		m.SetHeader("From", "girqta@gmail.com")
		m.SetHeader("To", "kevinc16c@gmail.com")
		m.SetHeader("Subject", "Datos de acceso Vivir Carlos Paz")
		m.SetBody("text/html", `
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html xmlns="http://www.w3.org/1999/xhtml">
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<title>Demystifying Email Design</title>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
			<style>
				.city {
				background-color: #ebebeb;
				color: white;
				padding: 10px;
				};
				.fonde {
				padding: 10px;
				background: black;
				}
			</style>
		</head>
		<body>
			<div class="container" style="padding: 8%">
			<img src="https://www.vivircarlospaz.com/img/static/vcp-logogv.png" style="width:250px;display: block;margin-left: auto;margin-right: auto;"/>
			<div class="container">
				<div class="row">
					<p class="text-center" style="font-size:20px">Tu pedido se está preparando</p>
					<p class="text-center" style="font-size:20px">Hora estimada de entrega: HH:MM y hh:mm hs.</p>
					<p class="text-center" style="font-size:15px">Comunicate con nosotros al <a href="https://www.vivircarlospaz.com/">+54 9 3482 58-1576</a></p>
				</div>
			</div>
				<div class="container" style="background-color:#d9d9d9; border-radius: 12px;">
					<div class="row">
						<table style="width:100%; padding-left:4%; padding-right:4%">
						<tr class="blank_row">
							<td colspan="3">&nbsp;</td>
						</tr>
						<tr>
							<td>1x Pizza de piña</td>
							<td style="text-align:right">$450</td>
						</tr>
						<tr class="blank_row">
							<td colspan="3">&nbsp;</td>
						</tr>
						<tr>
						<td>Sub-total</td>
							<td style="text-align:right">$450</td>
						</tr>
						<tr>
							<td>Costo de envio</td>
							<td style="text-align:right">$450</td>
						</tr>
						<tr style="color:#6a8700">
							<td>Descuento por cupón</td>
							<td style="text-align:right">$450</td>
						</tr>
						<tr style="font-weight:bold; font-size:20px">
							<td>Total</td>
							<td style="text-align:right">$450</td>
						</tr>
						<tr class="blank_row">
							<td colspan="3">&nbsp;</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	</body>
	</html>`)

		// Send the email to Bob
		d := gomail.NewPlainDialer("smtp.gmail.com", 587, "girqta@gmail.com", "ytumerecibesasi")
		if err := d.DialAndSend(m); err != nil {
			return c.JSON(http.StatusOK, Respuesta{
				Status:  "errord",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Ok",
		})
	} else {

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "error",
			Message: "No se encontró propietario",
		})

	}
}
