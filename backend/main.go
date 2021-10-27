package main

import (
	_ "./config_time_zone"

	"./database"
	"golang.org/x/crypto/acme/autocert"

	"./config"

	AlicuotasIvaControlador "./controladores/alicuotas_iva"
	CondicionesIvaControlador "./controladores/condiciones_iva"
	DiasControlador "./controladores/dias"
	EventosControlador "./controladores/eventos"
	EventosImagenesControlador "./controladores/eventos_img"
	EventosRedesControlador "./controladores/eventos_rrss"
	HorasControlador "./controladores/horas"
	Imprimir_qrControlador "./controladores/imprimir_qr"
	InsumosControlador "./controladores/insumos"
	LocalidadesControlador "./controladores/localidades"
	LugaresControlador "./controladores/lugares"
	LugaresCuentasPagoControlador "./controladores/lugares_cuentas_pago"
	LugaresHorariosControlador "./controladores/lugares_horarios"
	LugaresImagenesControlador "./controladores/lugares_img"
	LugaresPalabrasControlador "./controladores/lugares_plc"
	LugaresRedesControlador "./controladores/lugares_rrss"
	LugaresSesionesControlador "./controladores/lugares_sesiones"
	MinutosControlador "./controladores/minutos"
	NotificacionesControlador "./controladores/notificaciones"
	OperadoresControlador "./controladores/operadores"
	OperadoresNivelesControlador "./controladores/operadores_niveles"
	PagosControlador "./controladores/pagos"
	PaisesControlador "./controladores/paises"
	PalabrasClaveControlador "./controladores/palabras_clave"
	PedidosControlador "./controladores/pedidos_control"
	PedidosDetallePago "./controladores/pedidos_detalle_pago"
	PedidosEstadosControlador "./controladores/pedidos_estados"
	PersonasHumanasJuridicasControlador "./controladores/personas_humanas_juridicas"
	ProductosControlador "./controladores/productos"
	ProductosCategoriasControlador "./controladores/productos_categorias"
	ProductosImagenesControlador "./controladores/productos_img"
	ProductosInsumosControlador "./controladores/productos_insumos"
	PromocionesControlador "./controladores/promociones"
	ProvinciasControlador "./controladores/provincias"
	RedesSocialesControlador "./controladores/redes_sociales"
	RegistroUsuariosLugaresControlador "./controladores/registro_usuarios_lugares"
	RubrosControlador "./controladores/rubros"
	SubrubrosControlador "./controladores/subrubros"
	TiposConveniosControlador "./controladores/tipos_convenio"
	TiposDeliveryControlador "./controladores/tipos_delivery"
	TiposEventosControlador "./controladores/tipos_eventos"
	TiposLugaresControlador "./controladores/tipos_lugares"
	TurnosFarmaciasControlador "./controladores/turnos_farmacias"
	UsuariosControlador "./controladores/usuarios"
	UsuariosLugaresFavoritosControlador "./controladores/usuarios_lugares_favoritos"
	UsuariosSesionesControlador "./controladores/usuarios_sesiones"
	ValoracionesControlador "./controladores/valoraciones"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"golang.org/x/crypto/acme/autocert"
)

func main() {
	database.InitDb()

	e := echo.New()

	// Configuro los static's
	e.Static("/", "/home/ubuntu/vivircarlospaz/frontend/build")
	e.Static("/img", "/home/ubuntu/vivircarlospaz/img")
	//e.Static("/static", "/home/ubuntu/vivircarlospaz/static")

	// certificados ssl
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("www.vivircarlospaz.com")
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	//
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Bienvenido!")
	// })

	/***************************************************************************************************************
	*** RUTAS PARA NOTIFICACIONES DE PAGOS DE MERCADO PAGO *********************************************************
	***************************************************************************************************************/
	mp := e.Group("/mercadopago")
	mp.POST("/notificaciones/pedidos", PedidosDetallePago.NotificacionPago)

	/***************************************************************************************************************
	*** RUTAS PARA GRUPO ADMIN *************************************************************************************
	***************************************************************************************************************/
	adminPublico := e.Group("/admin")
	admin := e.Group("/admin")
	admin.Use(middleware.JWT([]byte(config.SecretJwt)))

	// Operadores
	adminPublico.POST("/login", OperadoresControlador.Login)
	admin.GET("/autenticacion/operador", OperadoresControlador.GetAutenticacionOperador)
	admin.GET("/operadores", OperadoresControlador.Paginacion)
	admin.GET("/operadores/lista", OperadoresControlador.Lista)
	admin.GET("/operadores/:id", OperadoresControlador.GetOperador)
	admin.POST("/operadores", OperadoresControlador.Alta)
	admin.PUT("/operadores", OperadoresControlador.Modificar)
	admin.PUT("/operadores/baja/:id", OperadoresControlador.Baja)
	admin.PUT("/operadores/habilitar/:id", OperadoresControlador.Habilitar)
	admin.PUT("/operadores/clave", OperadoresControlador.Contrasena)
	admin.PUT("/operadores/cambiarClave", OperadoresControlador.CambiarContrasena)

	// Operadores niveles
	admin.GET("/niveles_operadores/lista", OperadoresNivelesControlador.Lista)
	admin.GET("/niveles_operadores/:id", OperadoresNivelesControlador.GetNivel)

	// Paises
	admin.GET("/paises", PaisesControlador.Paginacion)
	admin.GET("/paises/lista", PaisesControlador.Lista)
	admin.GET("/paises/:id", PaisesControlador.GetPais)
	admin.POST("/paises", PaisesControlador.Alta)
	admin.PUT("/paises", PaisesControlador.Modificar)

	// Provincias
	admin.GET("/provincias", ProvinciasControlador.Paginacion)
	admin.GET("/provincias/lista", ProvinciasControlador.Lista)
	admin.GET("/provincias/:id", ProvinciasControlador.GetProvincia)
	admin.POST("/provincias", ProvinciasControlador.Alta)
	admin.PUT("/provincias", ProvinciasControlador.Modificar)

	// Localidades
	admin.GET("/localidades", LocalidadesControlador.Paginacion)
	admin.GET("/localidades/lista", LocalidadesControlador.Lista)
	admin.GET("/localidades/:id", LocalidadesControlador.GetLocalidad)
	admin.POST("/localidades", LocalidadesControlador.Alta)
	admin.PUT("/localidades", LocalidadesControlador.Modificar)

	// Alicuotas de IVA
	admin.GET("/alicuotas/lista", AlicuotasIvaControlador.Lista)
	admin.GET("/alicuotas/:id", AlicuotasIvaControlador.GetAlicuota)

	// Condiciones de IVA
	admin.GET("/condicionesiva/lista", CondicionesIvaControlador.Lista)
	admin.GET("/condicionesiva/:id", CondicionesIvaControlador.GetCondicion)

	// Personas humanas juridicas
	admin.GET("/propietarios", PersonasHumanasJuridicasControlador.Paginacion)
	admin.GET("/propietarios/lista", PersonasHumanasJuridicasControlador.Lista)
	admin.GET("/propietarios/:id", PersonasHumanasJuridicasControlador.GetPersona)
	admin.POST("/propietarios", PersonasHumanasJuridicasControlador.Alta)
	admin.PUT("/propietarios", PersonasHumanasJuridicasControlador.Modificar)
	admin.PUT("/propietarios/baja/:id", PersonasHumanasJuridicasControlador.Baja)
	admin.PUT("/propietarios/habilitar/:id", PersonasHumanasJuridicasControlador.Habilitar)
	admin.PUT("/propietarios/clave", PersonasHumanasJuridicasControlador.Contrasena)
	admin.PUT("/propietarios/cambiarClave", PersonasHumanasJuridicasControlador.CambiarContrasena)
	admin.POST("/propietarios/emailclave/:id", PersonasHumanasJuridicasControlador.EnviarClaveEmail)

	// Rubros
	admin.GET("/rubros", RubrosControlador.Paginacion)
	admin.GET("/rubros/lista", RubrosControlador.Lista)
	admin.GET("/rubros/:id", RubrosControlador.GetRubro)
	admin.POST("/rubros", RubrosControlador.Alta)
	admin.PUT("/rubros", RubrosControlador.Modificar)
	// admin.DELETE("/rubros/:id", RubrosControlador.Baja)

	// Subrubros
	admin.GET("/subrubros", SubrubrosControlador.Paginacion)
	admin.GET("/subrubros/lista", SubrubrosControlador.Lista)
	admin.GET("/subrubros/rubro/:idrubro", SubrubrosControlador.SubrubrosRubro)
	admin.GET("/subrubros/:id", SubrubrosControlador.GetSubrubro)
	admin.POST("/subrubros", SubrubrosControlador.Alta)
	admin.PUT("/subrubros", SubrubrosControlador.Modificar)

	// Tipos de convenios
	admin.GET("/convenios/lista", TiposConveniosControlador.Lista)
	admin.GET("/convenios/:id", TiposConveniosControlador.GetTipoConvenio)

	// Tipos de delivery
	admin.GET("/tiposdelivery/lista", TiposDeliveryControlador.Lista)
	admin.GET("/tiposdelivery/:id", TiposDeliveryControlador.GetTipoDelivery)

	// Tipos de lugares
	admin.GET("/tiposlugares/lista", TiposLugaresControlador.Lista)
	admin.GET("/tiposlugares/:id", TiposLugaresControlador.GetTipoLugar)

	// Tipos de eventos
	admin.GET("/tiposeventos/lista", TiposEventosControlador.Lista)
	admin.GET("/tiposeventos/:id", TiposEventosControlador.GetTipoEvento)

	// Dias
	admin.GET("/dias/lista", DiasControlador.Lista)
	admin.GET("/dias/:id", DiasControlador.GetDia)

	// Dias
	admin.GET("/horas/lista", HorasControlador.Lista)
	admin.GET("/horas/:id", HorasControlador.GetHora)

	// Minutos
	admin.GET("/minutos/lista", MinutosControlador.Lista)
	admin.GET("/minutos/:id", MinutosControlador.GetMinuto)

	// Redes Sociales
	admin.GET("/redes", RedesSocialesControlador.Paginacion)
	admin.GET("/redes/lista", RedesSocialesControlador.Lista)
	admin.GET("/redes/:id", RedesSocialesControlador.GetRedSocial)
	admin.POST("/redes", RedesSocialesControlador.Alta)
	admin.PUT("/redes", RedesSocialesControlador.Modificar)
	admin.PUT("/redes/baja/:id", RedesSocialesControlador.Baja)
	admin.PUT("/redes/habilitar/:id", RedesSocialesControlador.Habilitar)

	// Palabras claves
	admin.GET("/palabras", PalabrasClaveControlador.Paginacion)
	admin.GET("/palabras/:idrubro/lista", PalabrasClaveControlador.Lista)
	admin.GET("/palabras/:id", PalabrasClaveControlador.GetPalabraClave)
	admin.POST("/palabras", PalabrasClaveControlador.Alta)
	admin.PUT("/palabras", PalabrasClaveControlador.Modificar)
	admin.PUT("/palabras/baja/:id", PalabrasClaveControlador.Baja)
	admin.PUT("/palabras/habilitar/:id", PalabrasClaveControlador.Habilitar)

	// Lugares
	admin.GET("/lugares", LugaresControlador.Paginacion)
	admin.GET("/lugares/lista", LugaresControlador.Lista)
	admin.GET("/lugares/subrubros/:id", LugaresControlador.ListaLugaresSubrubros)
	admin.GET("/lugares/:id", LugaresControlador.GetLugar)
	admin.POST("/lugares", LugaresControlador.Alta)
	admin.PUT("/lugares", LugaresControlador.Modificar)
	admin.PUT("/lugares/baja/:id", LugaresControlador.Baja)
	admin.PUT("/lugares/habilitar/:id", LugaresControlador.Habilitar)
	admin.GET("/lugares/qr/generar/:id", LugaresControlador.GenerarQr)
	admin.GET("/lugares/qr/imprimir/:id", LugaresControlador.ImprimirQr)

	// Lugares redes sociales
	admin.GET("/lugares/:id/redes", LugaresRedesControlador.Paginacion)
	admin.GET("/lugaresredes/:id", LugaresRedesControlador.GetRed)
	admin.POST("/lugaresredes", LugaresRedesControlador.Alta)
	admin.PUT("/lugaresredes", LugaresRedesControlador.Modificar)
	admin.DELETE("/lugaresredes/:id", LugaresRedesControlador.Baja)

	// Lugares horarios
	admin.GET("/lugares/:id/horarios", LugaresHorariosControlador.GetHorarios)
	admin.GET("/horarios/:id", LugaresHorariosControlador.GetHorario)
	admin.POST("/horarios", LugaresHorariosControlador.Alta)
	admin.PUT("/horarios", LugaresHorariosControlador.Modificar)
	admin.DELETE("/horarios/:id", LugaresHorariosControlador.Baja)

	// Lugares imagenes
	admin.GET("/lugares/:id/imagenes", LugaresImagenesControlador.Paginacion)
	admin.GET("/imagenes/:id", LugaresImagenesControlador.GetImagen)
	admin.POST("/imagenes", LugaresImagenesControlador.Alta)
	admin.DELETE("/imagenes/:id", LugaresImagenesControlador.Baja)

	// Lugares palabras
	admin.GET("/lugares/:id/palabras", LugaresPalabrasControlador.Paginacion)
	admin.GET("/lugarespalabras/:id", LugaresPalabrasControlador.GetPalabraClave)
	admin.POST("/lugarespalabras", LugaresPalabrasControlador.Alta)
	admin.DELETE("/lugarespalabras/:id", LugaresPalabrasControlador.Baja)

	// Eventos
	admin.GET("/eventos", EventosControlador.Paginacion)
	admin.GET("/eventos/:id", EventosControlador.GetEvento)
	admin.POST("/eventos", EventosControlador.Alta)
	admin.PUT("/eventos", EventosControlador.Modificar)
	admin.DELETE("/eventos/:id", EventosControlador.Baja)

	// Eventos redes sociales
	admin.GET("/eventos/:id/redes", EventosRedesControlador.Paginacion)
	admin.GET("/eventosredes/:id", EventosRedesControlador.GetRed)
	admin.POST("/eventosredes", EventosRedesControlador.Alta)
	admin.PUT("/eventosredes", EventosRedesControlador.Modificar)
	admin.DELETE("/eventosredes/:id", EventosRedesControlador.Baja)

	// Eventos imagenes
	admin.GET("/eventos/:id/imagenes", EventosImagenesControlador.Paginacion)
	admin.GET("/eventos/imagenes/:id", EventosImagenesControlador.GetImagen)
	admin.POST("/eventos/imagenes", EventosImagenesControlador.Alta)
	admin.DELETE("/eventos/imagenes/:id", EventosImagenesControlador.Baja)

	// Promociones
	admin.GET("/promociones", PromocionesControlador.Paginacion)
	admin.GET("/promociones/:id", PromocionesControlador.GetPromocion)
	admin.POST("/promociones", PromocionesControlador.Alta)
	admin.PUT("/promociones", PromocionesControlador.Modificar)
	admin.PUT("/promociones/baja/:id", PromocionesControlador.Baja)
	admin.PUT("/promociones/habilitar/:id", PromocionesControlador.Habilitar)
	admin.PUT("/promociones/cambiarimagen", PromocionesControlador.CambiarImagen)

	// Productos categorias
	admin.GET("/categorias", ProductosCategoriasControlador.Paginacion)
	admin.GET("/categorias/rubro/:idrubro", ProductosCategoriasControlador.CategoriasRubro)
	admin.GET("/categorias/:id", ProductosCategoriasControlador.GetCategoria)
	admin.POST("/categorias", ProductosCategoriasControlador.Alta)
	admin.PUT("/categorias", ProductosCategoriasControlador.Modificar)
	admin.PUT("/categorias/baja/:id", ProductosCategoriasControlador.Baja)
	admin.PUT("/categorias/habilitar/:id", ProductosCategoriasControlador.Habilitar)

	// Farmacias de turnos
	admin.GET("/farmacias/turnos", TurnosFarmaciasControlador.Paginacion)
	admin.GET("/farmacias/turnos/:id", TurnosFarmaciasControlador.GetTurno)
	admin.POST("/farmacias/turnos", TurnosFarmaciasControlador.Alta)
	admin.DELETE("/farmacias/turnos/baja/:id", TurnosFarmaciasControlador.Baja)

	// Pedidos
	admin.GET("/pedidos/delivery/:id", PedidosControlador.PedidosTipoDelivery)
	admin.GET("/pedidos/:id", PedidosControlador.GetPedido)
	admin.PUT("/pedidos", PedidosControlador.SetEstado)

	// Estados de pedidos
	admin.GET("/estados/lista", PedidosEstadosControlador.Lista)
	admin.GET("/estados/lista/:id", PedidosEstadosControlador.Lista)

	// Notificaciones
	admin.POST("/notificaciones/topicos", NotificacionesControlador.NotificarTopico)

	// Registro de usuarios a lugares
	admin.GET("/registro/usuarios/lugares/:id", RegistroUsuariosLugaresControlador.Consultar)

	// Usuarios app
	admin.GET("/usuarios/lista", UsuariosControlador.Lista)

	/***************************************************************************************************************
	*** RUTAS PARA GRUPO PANEL DE COMERCIO *************************************************************************
	***************************************************************************************************************/
	panelPublico := e.Group("/panel")
	panel := e.Group("/panel")
	panel.Use(middleware.JWT([]byte(config.SecretPanelJwt)))

	// Operadores
	panelPublico.POST("/login", PersonasHumanasJuridicasControlador.Login)
	panel.GET("/autenticacion/propietario", PersonasHumanasJuridicasControlador.GetAutenticacionPropietario)
	panel.PUT("/propietarios/crearClave", PersonasHumanasJuridicasControlador.Contrasena)
	panel.PUT("/propietarios/cambiarClave", PersonasHumanasJuridicasControlador.CambiarContrasena)

	panel.GET("/qr/imprimir/:id", Imprimir_qrControlador.QrImprimir)

	// Lugares
	panel.GET("/lugares/propietarios/:id", LugaresControlador.PaginacionLugaresPropietario)
	panel.GET("/lugares/lista/propietarios/:id", LugaresControlador.ListaLugaresPropietario)
	panel.GET("/lugares/qr/generar/:id", LugaresControlador.GenerarQr)
	panel.GET("/lugares/qr/imprimir/:id", LugaresControlador.ImprimirQr)
	panel.PUT("/lugares/abierto/:id", LugaresControlador.MarcarAbierto)
	panel.PUT("/lugares/cerrado/:id", LugaresControlador.MarcarCerrado)

	// Lugares Sesiones
	panel.POST("/lugares/sesiones", LugaresSesionesControlador.Alta)
	panel.POST("/lugares/sesiones/baja", LugaresSesionesControlador.Baja)

	// Lugares cuentas de pago
	panel.GET("/lugares/cuentas/:id", LugaresCuentasPagoControlador.CuentasPagoLugar)
	panel.GET("/cuentas/:id", LugaresCuentasPagoControlador.GetCuenta)
	panel.POST("/cuentas", LugaresCuentasPagoControlador.Alta)
	panel.DELETE("/cuentas/:id", LugaresCuentasPagoControlador.Baja)
	panel.PUT("/cuentas/suspender/:id", LugaresCuentasPagoControlador.Suspender)
	panel.PUT("/cuentas/:idlugar/habilitar/:id", LugaresCuentasPagoControlador.Habilitar)

	// Alicuotas de IVA
	panel.GET("/alicuotas/lista", AlicuotasIvaControlador.Lista)
	panel.GET("/alicuotas/:id", AlicuotasIvaControlador.GetAlicuota)

	// Insumos
	panel.GET("/insumos/lugar/:id", InsumosControlador.Paginacion)
	panel.GET("/insumos/:id", InsumosControlador.GetInsumo)
	panel.POST("/insumos", InsumosControlador.Alta)
	panel.PUT("/insumos", InsumosControlador.Modificar)
	panel.DELETE("/insumos/:id", InsumosControlador.Baja)
	panel.PUT("/insumos/suspender/:id", InsumosControlador.Suspender)
	panel.PUT("/insumos/habilitar/:id", InsumosControlador.Habilitar)

	// Productos
	panel.GET("/productos/propietarios/:id", ProductosControlador.PaginacionProductosPropietario)
	panel.GET("/productos/lugares/:id", ProductosControlador.PaginacionProductosLugar)
	panel.GET("/productos/lista/propietarios/:id", ProductosControlador.ListaProductosPropietario)
	panel.GET("/productos/lista/lugares/:id", ProductosControlador.ListaProductosLugar)
	panel.GET("/productos/:id", ProductosControlador.GetProducto)
	panel.POST("/productos", ProductosControlador.Alta)
	panel.PUT("/productos", ProductosControlador.Modificar)
	panel.PUT("/productos/baja/:id", ProductosControlador.Baja)
	panel.PUT("/productos/habilitar/:id", ProductosControlador.Habilitar)
	panel.PUT("/productos/agotado/:id", ProductosControlador.Agotado)
	panel.PUT("/productos/enstock/:id", ProductosControlador.EnStock)

	// Productos categorias
	panel.GET("/categorias/rubro/:idrubro", ProductosCategoriasControlador.CategoriasRubro)

	// Productos imagenes
	panel.GET("/productos/:id/imagenes", ProductosImagenesControlador.Paginacion)
	panel.GET("/productos/imagenes/:id", ProductosImagenesControlador.GetImagen)
	panel.POST("/productos/imagenes", ProductosImagenesControlador.Alta)
	panel.DELETE("/productos/imagenes/:id", ProductosImagenesControlador.Baja)

	// Productos insumos
	panel.GET("/productos/:id/insumos", ProductosInsumosControlador.Paginacion)
	panel.POST("/productos/:id/insumos", ProductosInsumosControlador.SetInsumos)

	// Promociones
	panel.GET("/promociones/lugar/:id", PromocionesControlador.PaginacionPromocionesLugar)
	panel.GET("/promociones/:id", PromocionesControlador.GetPromocion)
	panel.POST("/promociones", PromocionesControlador.Alta)
	panel.PUT("/promociones", PromocionesControlador.Modificar)
	panel.PUT("/promociones/baja/:id", PromocionesControlador.Baja)
	panel.PUT("/promociones/habilitar/:id", PromocionesControlador.Habilitar)
	panel.PUT("/promociones/cambiarimagen", PromocionesControlador.CambiarImagen)

	// Pedidos
	panel.GET("/pedidos/lugar/:id", PedidosControlador.PedidosLugar)
	panel.GET("/pedidos/:id", PedidosControlador.GetPedido)
	panel.PUT("/pedidos", PedidosControlador.SetEstado)

	// Estados de pedidos
	panel.GET("/estados/lista", PedidosEstadosControlador.Lista)
	panel.GET("/estados/lista/:id", PedidosEstadosControlador.ListaEstadoElegido)

	/***************************************************************************************************************
	*** RUTAS PARA GRUPO API ***************************************************************************************
	***************************************************************************************************************/
	v1 := e.Group("/api/v1")
	v1.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == config.ApiKey, nil
	}))

	// Usuarios
	v1.POST("/login", UsuariosControlador.Login)
	v1.POST("/usuarios", UsuariosControlador.Alta)
	v1.PUT("/usuarios", UsuariosControlador.Modificar)
	v1.POST("/usuarios/cuenta/redes", UsuariosControlador.CuentaRedesSociales)
	v1.PUT("/usuarios/crearClave", UsuariosControlador.CrearContrasena)
	v1.PUT("/usuarios/cambiarClave", UsuariosControlador.CambiarContrasena)

	// Lugares favoritos usuarios
	v1.GET("/usuarios/:id/favoritos", UsuariosLugaresFavoritosControlador.Paginacion)
	v1.POST("/usuarios/favoritos", UsuariosLugaresFavoritosControlador.Alta)
	v1.DELETE("/usuarios/favoritos/:id", UsuariosLugaresFavoritosControlador.Baja)

	// Sesiones de usuarios
	v1.POST("/usuarios/sesiones", UsuariosSesionesControlador.Alta)
	v1.POST("/usuarios/sesiones/baja", UsuariosSesionesControlador.Baja)

	// Registro de personas en lugares
	v1.POST("/registro/usuarios/lugares", RegistroUsuariosLugaresControlador.Alta)

	// Paises
	v1.GET("/paises/lista", PaisesControlador.Lista)
	v1.GET("/paises/:id", PaisesControlador.GetPais)

	// Provincias
	v1.GET("/provincias/lista", ProvinciasControlador.Lista)
	v1.GET("/provincias/:id", ProvinciasControlador.GetProvincia)

	// Localidades
	v1.GET("/localidades/lista", LocalidadesControlador.Lista)
	v1.GET("/localidades/:id", LocalidadesControlador.GetLocalidad)

	// Rubros
	v1.GET("/rubros/lista", RubrosControlador.Lista)
	v1.GET("/rubros/sitiosutiles", RubrosControlador.RubrosSitiosUtiles)
	v1.GET("/rubros/:id", RubrosControlador.GetRubro)

	// Subrubros
	v1.GET("/subrubros/lista", SubrubrosControlador.Lista)
	v1.GET("/subrubros/rubro/:idrubro", SubrubrosControlador.SubrubrosRubro)
	v1.GET("/subrubros/:id", SubrubrosControlador.GetSubrubro)

	// Lugares
	v1.GET("/lugares", LugaresControlador.LugaresTipo)
	v1.GET("/lugares/busqueda", LugaresControlador.LugaresBusqueda)
	v1.GET("/lugares/rubro", LugaresControlador.LugaresRubro)
	v1.GET("/lugares/:id", LugaresControlador.GetDetalleLugar)

	// Lugares imagenes
	v1.GET("/lugares/:id/imagenes", LugaresImagenesControlador.Paginacion)

	// Eventos
	v1.GET("/eventos", EventosControlador.Eventos)
	v1.GET("/eventos/:id", EventosControlador.GetEvento)

	// Productos
	v1.GET("/productos/lugares/:id", ProductosControlador.ListaProductosLugar)
	v1.GET("/productos/:id", ProductosControlador.GetProductoDetalle)

	// Promociones
	v1.GET("/promociones", PromocionesControlador.GetPromocionesExclusivas)
	v1.GET("/promociones/:id", PromocionesControlador.GetPromocion)
	v1.PUT("/promociones/canje", PromocionesControlador.Canje)

	// Valoraciones
	v1.GET("/valoraciones/lugares/:id", ValoracionesControlador.PaginacionValoracionLugar)
	v1.GET("/valoraciones/:idusuario/:idlugar", ValoracionesControlador.GetValoracionUsuarioLugar)
	v1.POST("/valoraciones", ValoracionesControlador.Alta)
	v1.DELETE("/valoraciones/:id", ValoracionesControlador.Baja)

	// Pedidos
	v1.GET("/pedidos/usuario/:id", PedidosControlador.PedidosUsuario)
	v1.GET("/pedidos/:id", PedidosControlador.GetPedido)
	v1.POST("/pedidos", PedidosControlador.Alta)
	v1.PUT("/pedidos", PedidosControlador.SetEstado)

	// Pagos
	v1.POST("/pagos/preferencia", PagosControlador.PreferenciaPago)

	e.Logger.Fatal(e.Start(":5600"))
	//e.Logger.Fatal(e.Start(":80"))
	//e.Logger.Fatal(e.StartAutoTLS(":443"))
	//e.Logger.Fatal(e.StartTLS(":443", config.RutaCertificadoSSL, config.RutaKeySSL))
}
