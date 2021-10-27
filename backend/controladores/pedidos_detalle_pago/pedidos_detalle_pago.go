package pedidos_detalle_pago

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	config "../../config"
	"../../database"
	Modelos "../../modelos"
	"../../utils"
	"github.com/labstack/echo"

	Notificaciones "../notificaciones"
)

type Respuesta struct {
	Status string `json:"status"`
}

type Usuario_pedido struct {
	Idusuario uint `json:"idusuario"`
}

type Lugar_pedido struct {
	Idlugar uint `json:"idlugar"`
}

type Mensaje_Estado struct {
	Mensaje string `json:"mensaje"`
}

type Pago struct {
	Id             uint   `json:"id"`
	Referencia     string `json:"external_reference"`
	Metodo         string `json:"payment_method_id"`
	Tipo           string `json:"payment_type_id"`
	Estado         string `json:"status"`
	Detalle_estado string `json:"status_detail"`
}

func NotificacionPago(c echo.Context) error {
	db := database.GetDb()

	// Tomo parametros recibidos en url
	topic := c.QueryParam("topic")
	id := c.QueryParam("id")

	log.Print("NOTIFICACION " + topic + " " + id)

	if topic == "payment" {

		// Consulto el pago
		url := config.UrlPagosMP + fmt.Sprintf("%v", id) + "?access_token=" + config.AccesTokenMP
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s\n", "error en respuesta")
		}

		// leo respuesta
		response, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		fmt.Printf("%s\n", response)

		// tomo datos para respuesta del json que nos devolvio MP
		var pago Pago
		json.Unmarshal([]byte(response), &pago)

		// actualizo estado del pedido a recibido
		if err := db.Exec("UPDATE pedidos_control SET idestado=1 WHERE id=?", utils.ParseInt(pago.Referencia)).Error; err != nil {
			fmt.Print(err)
		}

		// Guardo informacion de pago
		detalle_pago := new(Modelos.Pedidos_detalle_pago)
		detalle_pago.Idpedido = utils.ParseInt(pago.Referencia)
		detalle_pago.Idpago = pago.Id
		detalle_pago.Tipo = pago.Tipo
		detalle_pago.Metodo = pago.Metodo
		detalle_pago.Estado = pago.Estado
		detalle_pago.Detalle = pago.Detalle_estado

		// verifico si el pago ya existe
		var dp2 []Modelos.Pedidos_detalle_pago
		db.Raw("SELECT * FROM pedidos_detalle_pago WHERE idpedido=? AND idpago=?", detalle_pago.Idpedido, detalle_pago.Idpago).Scan(&dp2)
		if len(dp2) > 0 { // si existe un pago creado lo actualizo

			if err := db.Exec("UPDATE pedidos_detalle_pago SET estado=?, detalle=? WHERE idpedido=? AND idpago=?", detalle_pago.Estado, detalle_pago.Detalle, detalle_pago.Idpedido, detalle_pago.Idpago).Error; err != nil {
				fmt.Print(err)
			}

		} else { // si no existe pago lo agrego

			if err := db.Create(&detalle_pago).Error; err != nil {
				fmt.Print(err)
			}

			// Notifico a usuario y lugar
			// Tomo id usuario del pedido
			var usuario_pedido = new(Usuario_pedido)
			db.Raw("SELECT idusuario FROM pedidos_control WHERE id=?", utils.ParseInt(pago.Referencia)).Scan(&usuario_pedido)

			// Tomo id lugar del pedido
			var lugar_pedido = new(Lugar_pedido)
			db.Raw("SELECT idlugar FROM pedidos_control WHERE id=?", utils.ParseInt(pago.Referencia)).Scan(&lugar_pedido)

			// Envio notificacion a usuario
			var mensaje_estado = new(Mensaje_Estado)
			db.Raw("SELECT detalle AS mensaje FROM pedidos_estados WHERE id=1").Scan(&mensaje_estado)
			Notificaciones.NotificarUsuarios(usuario_pedido.Idusuario, "Pedido", mensaje_estado.Mensaje, config.AccionPedidos, "")
			Notificaciones.NotificarLugar(lugar_pedido.Idlugar, "Nuevo", "Ha ingresado un Nuevo pedido")

		}

	}

	if topic == "merchant_orders" {
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
	})
}
