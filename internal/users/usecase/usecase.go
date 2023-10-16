package usecase

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
)

type (
	Usecase interface {
		CreateUser(ctx context.Context, user entities.CreateUser) (int64, error)
		UpdateUser(ctx context.Context, updateData entities.UpdateUsers) error
		UpdateUserStatus(ctx context.Context, req entities.UpdateUserStatus) error
	}

	usecase struct {
		repo repository.Repository
		log  *zerolog.Logger
	}
)

func NewUseCase(repo repository.Repository, log *zerolog.Logger) Usecase {
	return &usecase{repo, log}
}

func (uc *usecase) CreateUser(ctx context.Context, payload entities.CreateUser) (int64, error) {
	userID, err := uc.repo.SaveNewUser(ctx, entities.NewCreateUser(payload))
	if err != nil {
		uc.log.Z().Err(err).Msg("[usecase]CreateUser.SaveNewUser")

		return userID, err
	}

	return userID, nil
}

func (uc *usecase) UpdateUser(ctx context.Context, updateData entities.UpdateUsers) error {
	return uc.repo.Atomic(ctx, &sql.TxOptions{}, func(tx repository.Repository) error {

		if _, err := tx.GetUserByID(ctx, updateData.UserID, entities.LockingOpt{PessimisticLocking: true}); err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.GetUserByID")

			return err
		}

		if err := tx.UpdateUserByID(ctx, entities.NewUpdateUsers(updateData)); err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.UpdateUserByID")

			return err
		}

		return nil
	})
}

func (uc *usecase) UpdateUserStatus(ctx context.Context, req entities.UpdateUserStatus) error {
	_, err := uc.repo.GetUserByID(ctx, req.UserID, entities.LockingOpt{})
	if err != nil {
		uc.log.Z().Err(err).Msg("[usecase]UpdateUserStatus.GetUserByID")

		return err
	}

	if err := uc.repo.UpdateUserStatusByID(ctx, entities.NewUpdateUserStatus(req)); err != nil {
		uc.log.Z().Err(err).Msg("[usecase]UpdateUserStatus.UpdateUserStatusByID")

		return err
	}

	return nil
}
