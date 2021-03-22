package push

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/lisomatrix/channels/channels/core"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
)

func InitializeAppDefault(filePath string) *firebase.App {
	opt := option.WithCredentialsFile(filePath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}

type FirePushNotificationHandler struct {
	firebaseApp *firebase.App
	client *messaging.Client
	queue chan *core.PushRequestItem
}

func (handler *FirePushNotificationHandler) EnqueueRequest(request *core.PushRequestItem) {
	handler.queue <- request
}

func (handler *FirePushNotificationHandler) sendMessages() {

	for {
		request, _ := <-handler.queue
		sendMulticast(request, handler.client)
	}
}

func NewPushNotificationHandler(filePath string) *FirePushNotificationHandler {
	handler := new(FirePushNotificationHandler)

	handler.firebaseApp = InitializeAppDefault(filePath)
	handler.queue = make(chan *core.PushRequestItem, 100)

	client, err := handler.firebaseApp.Messaging(context.Background())

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to get messaging client %v\n", err)
		return nil
	}

	handler.client = client

	// Start goroutine for sending push notifications
	go handler.sendMessages()

	return handler
}

func sendMulticast(pushRequestItem *core.PushRequestItem, client *messaging.Client) {

	// Create a list containing up to 100 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	tokens := make([]string, 0)

	for _, clientID := range pushRequestItem.ClientIDs {
		deviceTokens, err := core.GetEngine().GetDeviceRepository().GetClientDeviceTokens(clientID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to get client devices %v\n", err)
			continue
		}

		if len(tokens) + len(deviceTokens) > 500 {
			break
		}

		for _, token := range deviceTokens {
			if token != "" {
				tokens = append(tokens, token)
			}
		}
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"channelID": pushRequestItem.ChannelID,
			"timestamp":  strconv.FormatInt(pushRequestItem.Timestamp, 10),
			"eventType": pushRequestItem.EventType,
			"payload": pushRequestItem.Payload,
		},
		Tokens: tokens,
	}

	_, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to send push notifications %v\n", err)
	}
}