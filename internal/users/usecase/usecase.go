package usecase

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/database"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/jmoiron/sqlx"
)

type (
	Usecase interface {
		CreateUser(context.Context, *entities.Users) (int64, error)
		UpdateUser(context.Context, *entities.UpdateUsers) error
	}

	usecase struct {
		repo repository.Repository
		log  *zerolog.Logger
		dbTx *sqlx.DB
	}
)

func NewUseCase(repo repository.Repository, log *zerolog.Logger, txConn *sqlx.DB) Usecase {
	return &usecase{repo, log, txConn}
}

func (uc *usecase) CreateUser(ctx context.Context, user *entities.Users) (int64, error) {
	userID, err := uc.repo.SaveNewUser(ctx, entities.NewUser(*user))
	if err != nil {
		uc.log.Z().Err(err).Msg("[usecase]CreateUser.SaveNewUser")

		return userID, err
	}

	return userID, nil
}

func (uc *usecase) UpdateUser(ctx context.Context, updateData *entities.UpdateUsers) error {
	return func(dbTx *sqlx.DB) error {
		txConn, err := uc.dbTx.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		defer func() {
			if err := new(database.SQL).EndTx(txConn, err); err != nil {
				uc.log.Z().Err(err).Msg("[usecase]UpdateUser.EndTx")
			}
		}()

		repoTx := repository.NewRepository(txConn, uc.log)

		if _, err := repoTx.GetUserByID(ctx, updateData.UserID, entities.LockingOpt{ForUpdate: true}); err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.GetUserByID")
			return err
		}

		if err = repoTx.UpdateUserByID(ctx, entities.NewUpdateUsers(*updateData)); err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.UpdateUserByID")
			return err
		}

		return nil
	}(uc.dbTx)
}
