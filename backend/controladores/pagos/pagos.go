package pagos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../../config"
	"../../database"
	Modelos "../../modelos"

	"github.com/labstack/echo"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Preferencia Preferencia `json:"preferencia,omitempty"`
}

type Preferencia struct {
	Id               string `json:"id"`
	Initpoint        string `json:"init_point"`
	Sandboxinitpoint string `json:"sandbox_init_point"`
}

type DatosPago struct {
	Idlugar       uint    `json:"idlugar"`
	TipoDocumento string  `json:"tipoDoc"`
	NroDocumento  string  `json:"nroDoc"`
	Email         string  `json:"email"`
	Importe       float64 `json:"importe"`
	Comision      float64 `json:"comision"`
}

func PreferenciaPago(c echo.Context) error {
	db := database.GetDb()
	//fecha := time.Now() //utils.GetNow()

	datosPago := new(DatosPago)

	if err := c.Bind(datosPago); err != nil {
		response := Respuesta{
			Status:  "error",
			Message: "invalid request body",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Traigo los fatos de la cuenta habilitada del lugar
	cuentas := new(Modelos.Lugares_cuentas_pago)
	db.Raw("SELECT * FROM lugares_cuentas_pago WHERE idlugar=? AND suspendido=0", datosPago.Idlugar).
		Scan(&cuentas)

	// calculo importe de comision
	impComision := (datosPago.Importe * datosPago.Comision) / 100

	// Preparo json para enviar
	requestBody := strings.NewReader(
		`{
			"items": [
				{
					"title": "Compra en Vivir Carlos Paz App",
					"currency_id": "ARS",
					"quantity": 1,
					"unit_price": ` + fmt.Sprintf("%.2f", datosPago.Importe) + `
				}
			],
			"payer": {
				"email": ` + fmt.Sprintf("%v", datosPago.Email) + `
			},
			"payment_methods": {
				"excluded_payment_types": [
					{
						"id": "ticket"
					},
					{
						"id": "atm"
					}
				],
			},
			"marketplace_fee": ` + fmt.Sprintf("%.2f", impComision) + `,
		}`)

	// ,
	// 		"identification": {
	// 			"type": ` + fmt.Sprintf("%v", datosPago.TipoDocumento) + `,
	// 			"number": ` + fmt.Sprintf("%v", datosPago.NroDocumento) + `
	// 		}

	//fmt.Println("Request:", requestBody)

	// Envia los datos
	url := config.UrlPreferenciasMP + "?access_token=" + cuentas.Accesstoken
	res, err := http.Post(
		url,
		"application/json; charset=UTF-8",
		requestBody,
	)

	// Compruebo error en respuesta
	if err != nil {
		return c.JSON(http.StatusOK, Respuesta{
			Status:  "error",
			Message: "Error generando pago",
		})
	}

	// leo respuesta
	response, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// tomo datos para respuesta del json que nos devolvio MP
	var preferencia Preferencia
	json.Unmarshal([]byte(response), &preferencia)
	//fmt.Printf("%s\n", response)

	data := Data{Preferencia: preferencia}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
