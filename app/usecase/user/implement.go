package user_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/flash-cards-vocab/backend/pkg/helpers"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	// "github.com/opentracing/opentracing-go"
	// "github.com/pkg/errors"
	"github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	// "github.com/AleksK1NG/api-mc/config"
	// "github.com/AleksK1NG/api-mc/internal/auth"
	// "github.com/AleksK1NG/api-mc/internal/entity"
	// "github.com/AleksK1NG/api-mc/pkg/httpErrors"
	// "github.com/AleksK1NG/api-mc/pkg/logger"
	// "github.com/AleksK1NG/api-mc/pkg/utils"
)

// const (
// 	basePrefix    = "api-auth:"
// 	cacheDuration = 3600
// )

// Auth UseCase
// type authUC struct {
// 	cfg       *config.Config
// 	authRepo  auth.Repository
// 	redisRepo auth.RedisRepository
// 	awsRepo   auth.AWSRepository
// 	logger    logger.Logger
// }

type usecase struct {
	userRepo repository.UserRepository
}

func New(userRepo repository.UserRepository) UseCase {
	return &usecase{
		userRepo: userRepo,
	}
}

// Auth UseCase constructor
// func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, awsRepo auth.AWSRepository, log logger.Logger) auth.UseCase {
// 	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, awsRepo: awsRepo, logger: log}
// }

// Create new user
func (uc *usecase) Register(ctx context.Context, user entity.User) (*entity.UserWithToken, error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	// defer span.Finish()

	existsUser, err := uc.userRepo.CheckIfUserExistsByEmail(user.Email)
	if err != nil {
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error1")
		// return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}
	if existsUser {
		return nil, fmt.Errorf("%w: %v", ErrUserExistsAlready, "User exists already")
	}

	if err = user.PrepareCreate(); err != nil {
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Error while preparing user data to be created")
	}

	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	createdUser.Password = ""

	token, err := helpers.GenerateJWTToken(createdUser)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error3")
	}

	return &entity.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Login user, returns user model with jwt token
func (uc *usecase) Login(ctx context.Context, user entity.UserLogin) (*entity.UserWithToken, error) {
	foundUser, err := uc.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("%w: %v", ErrUserPasswordMismatch, "Password mismatch")
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	foundUser.Password = ""

	token, err := helpers.GenerateJWTToken(foundUser)
	if err != nil {
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	return &entity.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Update existing user
// func (uc *usecase) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Update")
// 	// defer span.Finish()

// 	if err := user.PrepareUpdate(); err != nil {
// 		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareUpdate"))
// 	}

// 	updatedUser, err := u.authRepo.Update(ctx, user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	updatedUser.SanitizePassword()

// 	if err = u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.UserID.String())); err != nil {
// 		u.logger.Errorf("AuthUC.Update.DeleteUserCtx: %s", err)
// 	}

// 	updatedUser.SanitizePassword()

// 	return updatedUser, nil
// }

// Delete new user
// func (uc *usecase) Delete(ctx context.Context, userID uuid.UUID) error {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Delete")
// 	// defer span.Finish()

// 	if err := u.authRepo.Delete(ctx, userID); err != nil {
// 		return err
// 	}

// 	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(userID.String())); err != nil {
// 		u.logger.Errorf("AuthUC.Delete.DeleteUserCtx: %s", err)
// 	}

// 	return nil
// }

// Get user by id
// func (uc *usecase) GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
// 	// defer span.Finish()

// 	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
// 	if err != nil {
// 		u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
// 	}
// 	if cachedUser != nil {
// 		return cachedUser, nil
// 	}

// 	user, err := u.authRepo.GetByID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID.String()), cacheDuration, user); err != nil {
// 		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
// 	}

// 	user.SanitizePassword()

// 	return user, nil
// }

// Find users by name
// func (uc *usecase) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*entity.UsersList, error) {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.FindByName")
// 	// defer span.Finish()

// 	return u.authRepo.FindByName(ctx, name, query)
// }

// // Get users with pagination
// func (uc *usecase) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error) {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetUsers")
// 	// defer span.Finish()

// 	return u.authRepo.GetUsers(ctx, pq)
// }

// Upload user avatar
// func (uc *usecase) UploadAvatar(ctx context.Context, userID uuid.UUID, file entity.UploadInput) (*entity.User, error) {
// 	// span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
// 	// defer span.Finish()

// 	uploadInfo, err := u.awsRepo.PutObject(ctx, file)
// 	if err != nil {
// 		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.PutObject"))
// 	}

// 	avatarURL := u.generateAWSMinioURL(file.BucketName, uploadInfo.Key)

// 	updatedUser, err := u.authRepo.Update(ctx, &entity.User{
// 		UserID: userID,
// 		Avatar: &avatarURL,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	updatedUser.SanitizePassword()

// 	return updatedUser, nil
// }

// func (uc *usecase) GenerateUserKey(userID string) string {
// 	return fmt.Sprintf("%s: %s", basePrefix, userID)
// }

// func (uc *usecase) generateAWSMinioURL(bucket string, key string) string {
// 	return fmt.Sprintf("%s/minio/%s/%s", uc.cfg.AWS.MinioEndpoint, bucket, key)
// }
