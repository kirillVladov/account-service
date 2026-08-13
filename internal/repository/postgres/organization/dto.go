package organization_repo

import (
	"github.com/kirillVladov/account-service/internal/application/dto"
)

type organization struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func convertToApplication(in organization) dto.Organization {
	return dto.Organization{
		ID:   in.ID,
		Name: in.Name,
	}
}

func convertToRepository(in dto.Organization) organization {
	return organization{
		ID:   in.ID,
		Name: in.Name,
	}
}