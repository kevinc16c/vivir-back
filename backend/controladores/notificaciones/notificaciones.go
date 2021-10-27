package notificaciones

import (
	"context"
	"fmt"
	"net/http"

	config "../../config"
	"../../database"
	"github.com/labstack/echo"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type Respuesta struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type Tokens struct {
	Token string `json:"token"`
}

type Mensaje_topicos struct {
	Topico      string `json:"topico"`
	Titulo      string `json:"titulo"`
	Mensaje     string `json:"mensaje"`
	Vencimiento string `json:"vencimiento"`
}

func NotificarUsuarios(idusuario uint, titulo string, mensaje string, accion string, vencimiento string) {
	db := database.GetDb()

	// Traigo los tokens del usuario a notificar
	var tokens []Tokens
	db.Raw("SELECT token FROM usuarios_app_sesiones WHERE idusuario=? AND token!=''", idusuario).Scan(&tokens)
	if len(tokens) > 0 {

		// Registro a firebase
		opt := option.WithCredentialsFile(config.RutaConfigFirebase)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error initializing app: %v\n", err)
		}

		// Obtengo un messaging.Client para la app.
		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			fmt.Printf("error getting Messaging client: %v\n", err)
		}

		// creo el mensaje
		msg := map[string]string{
			"title":       titulo,
			"message":     mensaje,
			"action":      accion,
			"vencimiento": vencimiento,
		}

		if len(tokens) == 1 { // SI HAY UN SOLO TOKEN PARA NOTIFICAR

			// Token para enviar mensaje
			registrationToken := tokens[0].Token

			message := &messaging.Message{
				Data:  msg,
				Token: registrationToken,
			}

			// envia el mensaje al token correspondiente
			response, err := client.Send(ctx, message)
			if err != nil {
				fmt.Printf("error " + fmt.Sprint(err))
			}

			fmt.Println("Successfully sent message:", response)

		} else { // SI HAY MAS DE UN TOKEN PARA NOTIFICAR

			// Lista de tockens a enviar mensaje
			var registrationTokens []string
			for i := 0; i < len(tokens); i++ {
				registrationTokens = append(registrationTokens, tokens[i].Token)
			}

			message := &messaging.MulticastMessage{
				Data:   msg,
				Tokens: registrationTokens,
			}

			br, err := client.SendMulticast(context.Background(), message)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}

			fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
		}
	}
}

func NotificarLugar(idlugar uint, titulo string, mensaje string) {
	db := database.GetDb()

	// Traigo los tokens del usuario a notificar
	var tokens []Tokens
	db.Raw("SELECT token FROM lugares_sesiones WHERE idlugar=? AND token!=''", idlugar).Scan(&tokens)
	if len(tokens) > 0 {

		// Registro a firebase
		opt := option.WithCredentialsFile(config.RutaConfigFirebase)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error initializing app: %v\n", err)
		}

		// Obtengo un messaging.Client para la app.
		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			fmt.Printf("error getting Messaging client: %v\n", err)
		}

		// creo el mensaje
		msg := map[string]string{
			"title": titulo,
			"body":  mensaje,
		}

		if len(tokens) == 1 { // SI HAY UN SOLO TOKEN PARA NOTIFICAR

			// Token para enviar mensaje
			registrationToken := tokens[0].Token

			message := &messaging.Message{
				Data:  msg,
				Token: registrationToken,
			}

			// envia el mensaje al token correspondiente
			response, err := client.Send(ctx, message)
			if err != nil {
				fmt.Printf("error " + fmt.Sprint(err))
			}

			fmt.Println("Successfully sent message:", response)

		} else { // SI HAY MAS DE UN TOKEN PARA NOTIFICAR

			// Lista de tockens a enviar mensaje
			var registrationTokens []string
			for i := 0; i < len(tokens); i++ {
				registrationTokens = append(registrationTokens, tokens[i].Token)
			}

			message := &messaging.MulticastMessage{
				Data:   msg,
				Tokens: registrationTokens,
			}

			br, err := client.SendMulticast(context.Background(), message)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}

			fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
		}
	}
}

func NotificarTopico(c echo.Context) error {

	mensaje_topicos := new(Mensaje_topicos)

	// Toma los datos del body del post y controla que los datos hayan llegado bien
	if err := c.Bind(mensaje_topicos); err != nil {
		respuesta := Respuesta{
			Status:  "error",
			Message: "invalid request body " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, respuesta)
	}

	// Registro a firebase
	opt := option.WithCredentialsFile(config.RutaConfigFirebase)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}

	// Obtengo un messaging.Client para la app.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Printf("error getting Messaging client: %v\n", err)
	}

	// Creo el mensaje
	msg := map[string]string{
		"title":       mensaje_topicos.Titulo,
		"message":     mensaje_topicos.Mensaje,
		"action":      mensaje_topicos.Topico,
		"vencimiento": mensaje_topicos.Vencimiento,
	}

	message := &messaging.Message{
		Data:  msg,
		Topic: mensaje_topicos.Topico,
	}

	// Envia notificaciones a usuarios subscriptos a topicos
	response, err := client.Send(ctx, message)
	if err != nil {
		fmt.Printf("error " + fmt.Sprint(err))
	}

	return c.JSON(http.StatusOK, Respuesta{
		Status:  "success",
		Message: "Successfully sent message: " + response,
	})
}
