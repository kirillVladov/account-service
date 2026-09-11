package di

import "github.com/kirillVladov/account-service/pkg/notification_client"

func (di *DI) NotificationServiceClient() *notification_client.NotificationClient {
	if di.notificatoinClient != nil {
		return di.notificatoinClient
	}

	client, err := notification_client.NewNotificationClient(di.Config().NotificationServiceAddr.Addr)
	if err != nil {
		panic(err)
	}

	di.notificatoinClient = client

	return client
}
