package usecase

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/encryptions"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/security"
	"github.com/invopop/validation"
)

type usecase struct {
	repo   users.Repository
	appJwt *security.JWT
	crypto *encryptions.Crypto
}

type Dependencies struct {
	UseCases
	Repositories
	Pkgs
}

type UseCases struct{}

type Repositories struct {
	Repo users.Repository
}

type Pkgs struct {
	AppJwt *security.JWT
	Crypto *encryptions.Crypto
}

func (d Dependencies) toUsecase() *usecase {
	return &usecase{
		repo:   d.Repo,
		appJwt: d.AppJwt,
		crypto: d.Crypto,
	}
}

func NewUseCase(d Dependencies) users.Usecase {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.AppJwt, validation.Required),
		validation.Field(&d.Crypto, validation.Required),
		validation.Field(&d.Repo, validation.Required),
	)

	if err != nil {
		panic(err)
	}

	return d.toUsecase()
}
