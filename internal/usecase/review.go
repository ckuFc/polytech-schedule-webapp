package usecase

import (
	"context"
	"errors"
	"polytech_timetable/internal/entity"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/model/converter"
	"polytech_timetable/internal/repository"
	"polytech_timetable/pkg/customerrors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReviewUseCase struct {
	log         *logrus.Logger
	validate    *validator.Validate
	userRepo    *repository.UserRepository
	reviewRepo  *repository.ReviewRepository
	teacherRepo *repository.TeacherRepository
}

func NewReviewUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	userRepo *repository.UserRepository,
	reviewRepo *repository.ReviewRepository,
	teacherRepo *repository.TeacherRepository,
) *ReviewUseCase {
	return &ReviewUseCase{
		log:         log,
		validate:    validate,
		userRepo:    userRepo,
		reviewRepo:  reviewRepo,
		teacherRepo: teacherRepo,
	}
}

func (uc *ReviewUseCase) ListAllTeachers(ctx context.Context) ([]model.TeacherResponse, error) {
	teachers, err := uc.teacherRepo.FindAllWithReviewCount(ctx)
	if err != nil {
		uc.log.Errorf("Failed to list all teachers: %+v", err)
		return nil, customerrors.ErrInternalServerError
	}

	result := make([]model.TeacherResponse, 0, len(teachers))
	for i := range teachers {
		result = append(result, model.TeacherResponse{
			ID:          teachers[i].ID,
			Name:        teachers[i].Name,
			ReviewCount: teachers[i].ReviewCount,
		})
	}

	return result, nil
}

func (uc *ReviewUseCase) GetReviews(ctx context.Context, request *model.GetReviewsRequest, currentUserID int64) (*model.GetReviewsResponse, error) {
	if err := uc.validate.Struct(request); err != nil {
		uc.log.Warnf("Invalid get review request: %+v", err)
		return nil, customerrors.ErrValidation
	}
	teacher, err := uc.teacherRepo.FindByTeacherID(ctx, request.TeacherID)
	if err != nil {
		return nil, customerrors.ErrNotFound
	}

	reviews, err := uc.reviewRepo.FindByTeacherID(ctx, request.TeacherID)
	if err != nil {
		uc.log.Errorf("Failed to get reviews: %+v", err)
		return nil, customerrors.ErrInternalServerError
	}

	reviewResponses := make([]model.ReviewResponse, 0, len(reviews))
	for i := range reviews {
		reviewResponses = append(reviewResponses, converter.ReviewToResponse(&reviews[i], currentUserID))
	}

	return &model.GetReviewsResponse{
		Teacher: converter.TeacherToResponse(&teacher, int64(len(reviews))),
		Reviews: reviewResponses,
	}, nil
}

func (uc *ReviewUseCase) CreateReview(ctx context.Context, request *model.CreateReviewRequest) error {
	if err := uc.validate.Struct(request); err != nil {
		uc.log.Warnf("Invalid create review request: %+v", err)
		return customerrors.ErrValidation
	}

	_, err := uc.teacherRepo.FindByTeacherID(ctx, request.TeacherID)
	if err != nil {
		return customerrors.ErrNotFound
	}

	_, err = uc.reviewRepo.FindByTeacherAndUser(ctx, request.TeacherID, request.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			review := entity.TeacherReview{
				TeacherID: request.TeacherID,
				UserID:    request.UserID,
				Comment:   strings.TrimSpace(request.Comment),
			}

			if err := uc.reviewRepo.CreateReview(ctx, &review); err != nil {
				uc.log.Errorf("Failed to create review: %+v", err)
				return customerrors.ErrInternalServerError
			}
			return nil

		} else {
			uc.log.Errorf("failed to check existing review: %+v", err)
			return customerrors.ErrInternalServerError
		}
	}

	return customerrors.ErrAlreadyExists
}

func (uc *ReviewUseCase) DeleteReview(ctx context.Context, request *model.DeleteReviewRequest) error {
	if err := uc.validate.Struct(request); err != nil {
		uc.log.Warnf("Invalid delete review request: %+v", err)
		return customerrors.ErrValidation
	}

	_, err := uc.reviewRepo.FindByTeacherAndUser(ctx, request.TeacherID, request.UserID)
	if err != nil {
		return customerrors.ErrNotFound
	}
	if err := uc.reviewRepo.DeleteByTeacherAndUser(ctx, request.TeacherID, request.UserID); err != nil {
		uc.log.Errorf("Failed to delete review: %+v", err)
		return customerrors.ErrInternalServerError
	}

	return nil
}
