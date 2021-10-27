package pedidos_estados

import (
	"net/http"

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
	Registros     int                       `json:"registros,omitempty"`
	Estados_lista []Modelos.Pedidos_estados `json:"estados_lista,omitempty"`
	Estados       []Pedidos_estados_ruta    `json:"estados,omitempty"`
	Estado        *Pedidos_estados_ruta     `json:"estado,omitempty"`
}

type Pedidos_estados_ruta struct {
	Idestado uint   `json:"idestado"`
	Estado   string `json:"estado"`
	Detalle  string `json:"detalle"`
}

func Lista(c echo.Context) error {
	db := database.GetDb()

	if c.QueryParam("query") != "" {
		db = db.Where(" estado LIKE ? ", "%"+c.QueryParam("query")+"%")
	}

	// Ejecuto consulta
	var estados []Modelos.Pedidos_estados
	db.Find(&estados)
	data := Data{Estados_lista: estados}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}

func ListaEstadoElegido(c echo.Context) error {
	db := database.GetDb()
	idestado := c.Param("id")

	// Controlo parametro query recibido

	// Ejecuto consulta
	var estados []Pedidos_estados_ruta

	db = db.Raw("SELECT pedidos_estados_ruta.id, pedidos_estados.id AS idestado, pedidos_estados.estado, pedidos_estados.detalle FROM pedidos_estados_ruta JOIN pedidos_estados ON pedidos_estados.id=pedidos_estados_ruta.estfin WHERE pedidos_estados_ruta.estinicio=?", idestado).
		Scan(&estados)

	data := Data{Estados: estados}
	return c.JSON(http.StatusOK, Respuesta{
		Status: "success",
		Data:   data,
	})
}
