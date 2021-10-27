package config

const SecretJwt = "GiRqTa{vcp*2020yTmra"
const SecretPanelJwt = "GiRqTa{vcp*488kmdida"
const ApiKey = "4e3063e14993105b000a62fce89efb43b73e186bc43a1f71c7eea5d80782e1e1"

// Clave secreta para encriptar
const SecretKeyEncrypt = "GiRqTa8vcp142020"

// fecha por defecto de expiracion de token de usuarios
const ExpJWTUser = "31122100000000"

// Certificados SSL
const RutaCertificadoSSL = "/home/ubuntu/vivircarlospaz/backend/vivircarlospaz.com.crt"
const RutaKeySSL = "/home/ubuntu/vivircarlospaz/backend/vivircarlospaz.com.key"

// Api Key para firebase
const RutaConfigFirebase = "/home/ubuntu/vivircarlospaz/backend/vivir-carlos-paz-firebase.json"

// Accion de notificaciones
const AccionPedidos = "pedidos"

// Directorio del proyecto
const DirProyecto = "/home/ubuntu/vivircarlospaz"

const DirPlantillaPdfQr = "/home/ubuntu/vivircarlospaz/backend/static/plantilla_pdf_qr.jpg"

// Directorios de imagenes
const UrlImgLugares = "/img/lugares/"
const DirImgLugares = "/home/ubuntu/vivircarlospaz/img/lugares/"
const UrlImgEventos = "/img/eventos/"
const DirImgEventos = "/home/ubuntu/vivircarlospaz/img/eventos/"

const UrlQrLugares = "/img/lugares/qr/"
const DirQrLugares = "/home/ubuntu/vivircarlospaz/img/lugares/qr/"
const DirPlantillaQrLugares = "/home/ubuntu/vivircarlospaz/backend/static/plantilla_pdf_qr.jpg"

const UrlImgPromociones = "/img/promociones/"
const DirImgPromociones = "/home/ubuntu/vivircarlospaz/img/promociones/"

const UrlImgProductos = "/img/productos/"
const DirImgProductos = "/home/ubuntu/vivircarlospaz/img/productos/"

// Mercado Pago
const AccesTokenMP = "APP_USR-1182502572387360-121019-e4420c1be995b659ea4a90d0b15f4fba-685782144"
const UrlNotificacionesPagosPedidosMP = "https://www.vivircarlospaz.com/api/v1/pedidos/pagos/notificaciones?source_news=webhooks"
const UrlPreferenciasMP = "https://api.mercadopago.com/checkout/preferences"
const UrlPagosMP = "https://api.mercadopago.com/v1/payments/"
