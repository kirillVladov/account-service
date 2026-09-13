package notification

import (
	"context"

	"github.com/kirillVladov/account-service/pkg/notification_client"
)

type Gateway struct {
	client *notification_client.NotificationClient
}

func New(client *notification_client.NotificationClient) *Gateway {
	return &Gateway{client: client}
}

func (g *Gateway) SendEmailConfirmation(ctx context.Context, organizationID int64, confirmationLink, email string) error {
	return g.client.SendEmailConfirmation(ctx, organizationID, confirmationLink, email)
}
