package river_transport

import "github.com/google/uuid"

const AccountConfirmationQueue = "account_confirmation_workflow"

type AccountConfirmationEvent struct {
	AccountID        uuid.UUID `json:"id" river:"unique"`
	OrganizationID   int64     `json:"organization_id" river:"unique"`
	ConfirmationLink string    `json:"confirmation_link" river:"unique"`
}

func (AccountConfirmationEvent) Kind() string {
	return AccountConfirmationQueue
}
