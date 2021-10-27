package pedidos_control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	config "../../config"
	"../../database"
	Modelos "../../modelos"
	"../../utils"
	"github.com/labstack/echo"

	Notificaciones "../notificaciones"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
	Registros       int                       `json:"registros,omitempty"`
	Pedidos         []Modelos.Pedidos_control `json:"pedidos,omitempty"`
	Pedidos_usuario []Pedidos                 `json:"pedidos_usuario,omitempty"`
	Pedido          *Modelos.Pedidos_control  `json:"pedido,omitempty"`
	Preferencia     Preferencia               `json:"preferencia,omitempty"`
}

type Preferencia struct {
	Id               string `json:"id"`
	Initpoint        string `json:"init_point"`
	Sandboxinitpoint string `json:"sandbox_init_point"`
}

type Pedidos struct {
	Id             uint      `json:"id" gorm:"primary_key"`
	Fechaalta      time.Time `json:"fechaalta"`
	Idlugar        uint      `json:"idlugar"`
	Nombrelugar    string    `json:"nombrelugar"`
	Direccionlugar string    `json:"direccionlugar"`
	Celular        string    `json:"celular"`
	Localidad      string    `json:"localidad"`
	Rutafoto       string    `json:"rutafoto"`
	Idtiporetiro   uint      `json:"idtiporetiro"`
	Tipo_retiro    string    `json:"tipo_retiro"`
	Direccion      string    `json:"direccion"`
	Importe        float64   `json:"importe"`
	Idestado       uint      `json:"idestado"`
	Estado         string    `json:"estado"`
	Idtipodelivery int       `json:"idtipodelivery"`
	Tipodelivery   string    `json:"tipodelivery"`
	Tipopago       string    `json:"tipo_pago"`
	Metodopago     string    `json:"metodo_pago"`
	Observaciones  string    `json:"observaciones"`
}

func (Pedidos) TableName() string {
	return "pedidos_control"
}

type Estado_pedido struct {
	Idpedido uint `json:"idpedido"`
	Idestado uint `json:"idestado"`
}

type Usuario_pedido struct {
	Idusuario uint `json:"idusuario"`
}

type Mensaje_Estado struct {
	Mensaje string `json:"mensaje"`
}

func Paginacion(c echo.Context) error {
	db := database.GetDb()

	// Controlo valores para filtro y paginacion que llegan de la url
	if c.QueryParam("query") != "" {
		db = db.Where(" pedidos_control.id like ? ", "%"+c.QueryParam("query")+"%").
			Or(" apellido like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nombres like ? ", "%"+c.QueryParam("query")+"%").
			Or(" nombrelugar like ? ", "%"+c.QueryParam("query")+"%").
			Or(" direccion like ? ", "%"+c.QueryParam("query")+"%")
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fechaalta ASC")
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
	var pedidos []Modelos.Pedidos_control
	db.Preload("Usuario").Preload("Propietario").Preload("Lugar").Preload("Tiporetiro").Preload("Tipodelivery").Preload("Estado").Offset(offset).Limit(limite).Find(&pedidos)
	db.Table("pedidos_control").Count(&registros)
	data := Data{Registros: registros, Pedidos: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PedidosTipoDelivery(c echo.Context) error {
	db := database.GetDb()
	idtipodelivery := c.Param("id")
	idestado := c.QueryParam("estado")

	db = db.Select("pedidos_control.*,usuarios_app.apellido,usuarios_app.nombres")
	db = db.Joins("JOIN usuarios_app ON usuarios_app.id=pedidos_control.idusuario").
		Joins("JOIN lugares ON lugares.idlugar=pedidos_control.idlugar")

	if c.QueryParam("query") != "" {
		db = db.Where(" pedidos_control.idtipodelivery=? AND pedidos_control.idestado=? AND (pedidos_control.id like ? OR apellido like ? OR nombres like ? OR nombrelugar like ? OR pedidos_control.direccion like ? )", idtipodelivery, idestado, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" pedidos_control.idtipodelivery=? AND pedidos_control.idestado=?", idtipodelivery, idestado)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fechaalta ASC")
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
	var pedidos []Modelos.Pedidos_control
	db.Preload("Usuario").Preload("Propietario").Preload("Lugar.Localidad.Provincia").Preload("Tiporetiro").Preload("Tipodelivery").Preload("Estado").Offset(offset).Limit(limite).Find(&pedidos)
	db.Table("pedidos_control").Count(&registros)
	data := Data{Registros: registros, Pedidos: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PedidosLugar(c echo.Context) error {
	db := database.GetDb()
	idlugar := c.Param("id")
	idestado := c.QueryParam("estado")

	db = db.Select("pedidos_control.*,usuarios_app.apellido,usuarios_app.nombres")
	db = db.Joins("JOIN usuarios_app ON usuarios_app.id=pedidos_control.idusuario").
		Joins("JOIN lugares ON lugares.idlugar=pedidos_control.idlugar")

	if c.QueryParam("query") != "" {
		db = db.Where(" pedidos_control.idlugar=? AND pedidos_control.idestado=? AND (pedidos_control.id like ? OR apellido like ? OR nombres like ? OR nombrelugar like ? OR pedidos_control.direccion like ? )", idlugar, idestado, "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%", "%"+c.QueryParam("query")+"%")
	} else {
		db = db.Where(" pedidos_control.idlugar=? AND pedidos_control.idestado=?", idlugar, idestado)
	}

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fechaalta ASC")
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
	var pedidos []Modelos.Pedidos_control
	db.Preload("Usuario").Preload("Propietario").Preload("Lugar").Preload("Tiporetiro").Preload("Tipodelivery").Preload("Estado").Offset(offset).Limit(limite).Find(&pedidos)
	db.Table("pedidos_control").Count(&registros)
	data := Data{Registros: registros, Pedidos: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func PedidosUsuario(c echo.Context) error {
	db := database.GetDb()
	idusuario := c.Param("id")

	db = db.Select("pedidos_control.id,pedidos_control.fechaalta,pedidos_control.idlugar,lugares.nombrelugar,lugares.direccion AS direccionlugar,lugares.celular,localidades.nombrelocali AS localidad,(SELECT rutaimg FROM lugares_img WHERE idlugar=pedidos_control.idlugar LIMIT 1) AS rutafoto,pedidos_control.idtiporetiro,tipos_retiro.tipo AS tipo_retiro,pedidos_control.direccion,pedidos_control.importe,pedidos_control.idestado,pedidos_estados.estado,pedidos_control.idtipodelivery,tipos_delivery.tipodelivery,pedidos_detalle_pago.tipo AS tipopago,pedidos_detalle_pago.metodo AS metodopago,pedidos_control.observaciones")
	db = db.Joins("JOIN lugares ON lugares.idlugar=pedidos_control.idlugar").
		Joins("LEFT JOIN localidades ON localidades.idlocalidad=lugares.idlocalidad").
		Joins("LEFT JOIN tipos_retiro ON tipos_retiro.id=pedidos_control.idtiporetiro").
		Joins("LEFT JOIN pedidos_estados ON pedidos_estados.id=pedidos_control.idestado").
		Joins("LEFT JOIN tipos_delivery ON tipos_delivery.id=pedidos_control.idtipodelivery").
		Joins("LEFT JOIN pedidos_detalle_pago ON pedidos_detalle_pago.idpedido=pedidos_control.id")

	var condicion string = "pedidos_control.idusuario=" + fmt.Sprintf("%v", idusuario)

	if c.QueryParam("estado") == "pendientes" {
		condicion = condicion + " AND (pedidos_control.idestado=1 OR pedidos_control.idestado=2 OR pedidos_control.idestado=3 OR pedidos_control.idestado=4)"
	} else {
		condicion = condicion + " AND (pedidos_control.idestado=5 OR pedidos_control.idestado=6 OR pedidos_control.idestado=7)"
	}

	if c.QueryParam("query") != "" {
		condicion = condicion + " AND nombrelugar like " + "'%" + c.QueryParam("query") + "%'"
	}

	db = db.Where(condicion)

	if c.QueryParam("sortField") != "" {
		db = db.Order(c.QueryParam("sortField") + " " + c.QueryParam("sortOrder"))
	} else {
		db = db.Order("fechaalta DESC")
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
	var pedidos []Pedidos
	db.Offset(offset).Limit(limite).Find(&pedidos)
	db.Table("pedidos_control").Count(&registros)
	data := Data{Registros: registros, Pedidos_usuario: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func GetPedido(c echo.Context) error {
	db := database.GetDb()
	id := c.Param("id")

	pedidos := new(Modelos.Pedidos_control)
	db.Preload("Usuario").Preload("Propietario.Condicioniva").Preload("Propietario.Localidad.Provincia").Preload("Lugar.Localidad.Provincia").Preload("Tiporetiro").Preload("Tipodelivery").Preload("Estado").Preload("Detalle.Producto").Preload("Pago").Find(&pedidos, id)

	data := Data{Pedido: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func AltaOld(c echo.Context) error {
	db := database.GetDb()

	pedidos := new(Modelos.Pedidos_control)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(pedidos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Completo informacion de pedido
	pedidos.Fechaalta = time.Now() //utils.GetNow()
	pedidos.Idestado = 1           // Estado por defecto "Recibido"
	total := pedidos.Importe - pedidos.Impodelivery
	pedidos.Impcomision = (total * pedidos.Porcomision) / 100 // Importe de comision
	pedidos.Impneto = total - pedidos.Impcomision             // Neto

	// Inserta registro en la tabla
	if err := db.Omit("Emailusuario").Create(&pedidos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Notifico a usuario y lugar
	mensaje_estado := new(Mensaje_Estado)
	db.Raw("SELECT detalle AS mensaje FROM pedidos_estados WHERE id=?", pedidos.Idestado).Scan(&mensaje_estado)
	Notificaciones.NotificarUsuarios(pedidos.Idusuario, "Pedido", mensaje_estado.Mensaje, config.AccionPedidos, "")
	Notificaciones.NotificarLugar(pedidos.Idlugar, "Nuevo", "Nuevo pedido por un importe de $"+fmt.Sprintf("%.2f", pedidos.Importe))

	// Preparo mensaje de retorno
	data := Data{Pedido: pedidos}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func Alta(c echo.Context) error {
	db := database.GetDb()

	pedidos := new(Modelos.Pedidos_control)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(pedidos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Completo informacion de pedido
	pedidos.Fechaalta = time.Now() //utils.GetNow()
	pedidos.Idestado = 0           // Estado por defecto: sin estado hasta recibir informacion de pago
	total := pedidos.Importe - pedidos.Impodelivery
	pedidos.Impcomision = (total * pedidos.Porcomision) / 100 // Importe de comision
	pedidos.Impneto = total - pedidos.Impcomision             // Neto

	// Completo informacion de pago
	// Idtipopago: 1 - Efectivo; 2 - Mercado Pago
	if pedidos.Idtipopago == 1 {
		var pago Modelos.Pedidos_detalle_pago
		pago.Tipo = "efectivo"
		pago.Metodo = "efectivo"
		pago.Estado = "approved"

		pedidos.Idestado = 1 // Estado "Recibido" al ser pago en efectivo
		pedidos.Pago = pago
	}

	// Inserta registro en la tabla
	if err := db.Omit("Emailusuario").Create(&pedidos).Error; err != nil {
		response := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Creo la preferencia de pago
	var preferencia Preferencia

	// Idtipopago: 1 - Efectivo; 2 - Mercado Pago
	if pedidos.Idtipopago == 2 {
		// Traigo los fatos de la cuenta habilitada del lugar
		cuentas := new(Modelos.Lugares_cuentas_pago)
		db.Raw("SELECT * FROM lugares_cuentas_pago WHERE idlugar=? AND suspendido=0", pedidos.Idlugar).
			Scan(&cuentas)

		requestBody := strings.NewReader(
			`{
				"statement_descriptor":"VIVIR CARLOS PAZ"
				"items": [
					{
						"title": "Compra",
						"currency_id": "ARS",
						"quantity": 1,
						"unit_price": ` + fmt.Sprintf("%.2f", pedidos.Importe) + `
					}
				],
				"payer": {
					"email": ` + fmt.Sprintf("%v", pedidos.Emailusuario) + `
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
				"marketplace_fee": ` + fmt.Sprintf("%.2f", pedidos.Impcomision) + `,
				"external_reference": ` + fmt.Sprintf("%v", pedidos.Id) + `,
				"notification_url": "https://www.vivircarlospaz.com/mercadopago/notificaciones/pedidos?source_news=ipn",
			}`)

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

		json.Unmarshal([]byte(response), &preferencia)
	}

	// Notifico a usuario y lugar solo si el pago es en efectivo
	if pedidos.Idtipopago == 1 {
		mensaje_estado := new(Mensaje_Estado)
		db.Raw("SELECT detalle AS mensaje FROM pedidos_estados WHERE id=?", pedidos.Idestado).Scan(&mensaje_estado)
		Notificaciones.NotificarUsuarios(pedidos.Idusuario, "Pedido", mensaje_estado.Mensaje, config.AccionPedidos, "")
		Notificaciones.NotificarLugar(pedidos.Idlugar, "Nuevo", "Ha ingresado un Nuevo pedido")
	}

	// Preparo mensaje de retorno
	data := Data{Preferencia: preferencia}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func SetEstado(c echo.Context) error {
	db := database.GetDb()

	estado := new(Estado_pedido)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(estado); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Actualiza estado de pedido
	if err := db.Exec("UPDATE pedidos_control SET idestado=? WHERE id=?", estado.Idestado, estado.Idpedido).Error; err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Tomo id usuario del pedido
	var usuario_pedido = new(Usuario_pedido)
	db.Raw("SELECT idusuario FROM pedidos_control WHERE id=?", estado.Idpedido).Scan(&usuario_pedido)

	// Envio notificacion a usuario
	var mensaje_estado = new(Mensaje_Estado)
	db.Raw("SELECT detalle AS mensaje FROM pedidos_estados WHERE id=?", estado.Idestado).Scan(&mensaje_estado)
	Notificaciones.NotificarUsuarios(usuario_pedido.Idusuario, "Pedido", mensaje_estado.Mensaje, config.AccionPedidos, "")

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Estado asignado con éxito",
	})
}
