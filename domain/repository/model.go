package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(domain.Account, string) (domain.Model, error)
	GetByName(domain.Account, domain.ModelName) (domain.Model, error)
	List(domain.Account, ResourceListOption) ([]domain.Model, error)
	FindUserModels([]UserResourceListOption) ([]domain.Model, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error
}
