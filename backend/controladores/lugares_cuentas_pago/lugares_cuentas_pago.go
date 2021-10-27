package lugares_cuentas_pago

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
	Registros int                            `json:"registros,omitempty"`
	Cuentas   []Modelos.Lugares_cuentas_pago `json:"cuentas,omitempty"`
	Cuenta    *Modelos.Lugares_cuentas_pago  `json:"cuenta,omitempty"`
}

func CuentasPagoLugar(c echo.Context) error {
	db := database.GetDb()
	idlugar := c.Param("id")

	db = db.Select("id,idlugar,userid,vencimiento,suspendido")
	db = db.Where(" idlugar=? ", idlugar)

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
	var cuentas []Modelos.Lugares_cuentas_pago
	db.Offset(offset).Limit(limite).Find(&cuentas)
	db.Table("lugares_cuentas_pago").Count(&registros)
	data := Data{Registros: registros, Cuentas: cuentas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetCuenta(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	cuentas := new(Modelos.Lugares_cuentas_pago)
	db.Find(&cuentas, id)

	data := Data{Cuenta: cuentas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	cuentas := new(Modelos.Lugares_cuentas_pago)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(cuentas); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Inserta registro en la tabla
	ahora := time.Now()                                   //utils.GetNow()
	cuentas.Vencimiento = ahora.Add(180 * 24 * time.Hour) // sumo 180 dias a la fecha actual para el vencimiento del token de pago
	cuentas.Suspendido = 1
	if err := db.Create(&cuentas).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Preparo mensaje de retorno
	//data := Data{Cuenta: cuentas}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		//Data:   data,
		Message: "Cuenta dada de alta con éxito",
	})
}

func Baja(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("DELETE FROM lugares_cuentas_pago WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Cuenta dada de baja con éxito",
	})
}

func Suspender(c echo.Context) error {
	db := database.GetDb()

	if err := db.Exec("UPDATE lugares_cuentas_pago SET suspendido=1 WHERE id = ?", c.Param("id")).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Cuenta suspendida con éxito",
	})
}

func Habilitar(c echo.Context) error {
	db := database.GetDb()

	var cuentas []Modelos.Lugares_cuentas_pago
	db.Raw("SELECT * FROM lugares_cuentas_pago WHERE suspendido=0 AND idlugar=?", c.Param("idlugar")).Scan(&cuentas)
	if len(cuentas) > 0 {

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Ya existe una cuenta habilitada. Desabilitela antes de habilitar otra.",
		})

	} else {

		if err := db.Exec("UPDATE lugares_cuentas_pago SET suspendido=0 WHERE id = ?", c.Param("id")).Error; err != nil {
			respuesta := Respuesta{
				Status:  "error",
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, respuesta)
		}

		return c.JSON(http.StatusOK, Respuesta{
			Status:  "success",
			Message: "Cuenta habilitada con éxito",
		})

	}
}
