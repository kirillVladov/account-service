package notification_client

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kirillVladov/account-service/internal/gen/grpc"
	"github.com/kirillVladov/account-service/pkg/logger"
)

type NotificationClient struct {
	client pb.MailServiceClient
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &NotificationClient{
		client: pb.NewMailServiceClient(conn),
	}, nil
}

func (c *NotificationClient) SendEmailConfirmation(ctx context.Context, organizationID int64, confirmationLink, email string) error {
	l := logger.FromContext(ctx)

	l.Info("sending email confirmation",
		zap.Int64("organization_id", organizationID),
		zap.String("email", email),
	)

	_, err := c.client.SendEmailConfirmation(ctx, &pb.SendEmailConfirmationRequest{
		OrganizationId:   organizationID,
		ConfirmationLink: confirmationLink,
		Email:            email,
	})
	if err != nil {
		l.Error("failed to send email confirmation",
			zap.Int64("organization_id", organizationID),
			zap.String("email", email),
			zap.Error(err),
		)
		return err
	}

	l.Info("email confirmation sent",
		zap.Int64("organization_id", organizationID),
		zap.String("email", email),
	)

	return nil
}
