package usecase

import (
	"context"
	"polytech_timetable/internal/metrics"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/repository"
	"polytech_timetable/pkg/customerrors"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	log      *logrus.Logger
	validate *validator.Validate
	userRepo *repository.UserRepository
}

func NewUserUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	userRepo *repository.UserRepository,
) *UserUseCase {
	return &UserUseCase{
		log:      log,
		validate: validate,
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) SetGroup(ctx context.Context, request *model.SetGroupRequest) error {
	if err := uc.validate.Struct(request); err != nil {
		uc.log.Warnf("Invalid set group request: %+v", err)
		return customerrors.ErrValidation
	}

	_, err := uc.userRepo.FindByUserID(ctx, request.UserID)
	if err != nil {
		return customerrors.ErrNotFound
	}

	if err := uc.userRepo.SetUserGroup(ctx, request.UserID, request.Group); err != nil {
		uc.log.Errorf("Failed to set group: %+v", err)
		return customerrors.ErrInternalServerError
	}

	metrics.GroupChangesTotal.WithLabelValues(request.Group).Inc()
	return nil
}

func (uc *UserUseCase) UpdateSettings(ctx context.Context, request *model.SetNotificationsSettingRequest) error {
	if err := uc.validate.Struct(request); err != nil {
		uc.log.Warnf("Invalid set notifications setting request: %+v", err)
		return customerrors.ErrValidation
	}

	_, err := uc.userRepo.FindByUserID(ctx, request.UserID)
	if err != nil {
		return customerrors.ErrNotFound
	}

	if err := uc.userRepo.UpdateSettings(ctx, request.UserID, request.NotificationsEnabled); err != nil {
		uc.log.Errorf("Failed to update notifications: %+v", err)
		return customerrors.ErrInternalServerError
	}
	return nil
}
