package river_consumer

import "github.com/google/uuid"

type AccountConfirmationEvent struct {
	AccountID      uuid.UUID `json:"id" river:"unique"`
	OrganizationID int64     `json:"organization_id" river:"unique"`
}

func (AccountConfirmationEvent) Kind() string {
	return "account_confirmation_workflow"
}
