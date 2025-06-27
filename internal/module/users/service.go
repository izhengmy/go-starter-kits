package users

import (
	"app/internal/errorx"
	"app/internal/module/commons/authutil"
	"app/pkg/auth"
	"app/pkg/crypt"
	"app/pkg/jwtx"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserService struct {
	logger     *zap.Logger
	repository *UserRepository
}

func NewUserService(logger *zap.Logger, repository *UserRepository) *UserService {
	return &UserService{
		logger:     logger,
		repository: repository,
	}
}

func (s UserService) Register(request RegisterRequest) error {
	if exists := s.repository.ExistsByUsername(request.Username); exists {
		return errorx.NewServiceError("用户名已被占用")
	}

	hashedPassword, err := crypt.PasswordHash(request.Password)
	if err != nil {
		return errors.Wrap(err, "UserService.Register() failed to hash password")
	}

	user := User{
		Username: request.Username,
		Password: hashedPassword,
	}

	return s.repository.Create(user)
}

func (s UserService) Login(request LoginRequest) (*LoginResponse, error) {
	user, err := s.repository.FindByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errorx.NewServiceError("用户名或密码错误")
	}

	if !crypt.PasswordVerify(request.Password, user.Password) {
		return nil, errorx.NewServiceError("用户名或密码错误")
	}

	authUser := user.ToAuthenticationUser()
	token, err := auth.GenerateJWTToken(authUser)
	if err != nil {
		return nil, errors.Wrap(err, "UserService.Login() failed to generate jwt token")
	}

	return &LoginResponse{
		Token:     token,
		TokenType: jwtx.TokenType,
	}, nil
}

func (s UserService) Profile(authUser *authutil.AuthenticationUser[uint]) (*ProfileResponse, error) {
	if authUser == nil {
		return nil, errorx.NewServiceError("用户不存在")
	}

	user, err := s.repository.FindById(authUser.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errorx.NewServiceError("用户不存在")
	}

	return &ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
