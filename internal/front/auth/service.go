package auth

import (
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type AuthServiceCreator interface {
	Login(login_id string, password string) (*LoginReponse, *responses.ErrorResponse)
	GetUserInfo(userUUID string) (*UserInfoWithMenusResponse, *responses.ErrorResponse)
	Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse)
}

type AuthService struct {
	DBPool   *sqlx.DB
	AuthRepo *AuthRepoImpl
}

func NewAuthService(db_pool *sqlx.DB) *AuthService {
	return &AuthService{
		DBPool:   db_pool,
		AuthRepo: NewAuthRepoImpl(db_pool),
	}
}

func (au *AuthService) Login(login_id string, password string) (*LoginReponse, *responses.ErrorResponse) {
	return au.AuthRepo.Login(login_id, password)
}

func (au *AuthService) GetUserInfo(userUUID string) (*UserInfoWithMenusResponse, *responses.ErrorResponse) {
	return au.AuthRepo.GetUserInfoWithMenus(userUUID)
}

func (au *AuthService) Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse) {
	return au.AuthRepo.Register(register_req)
}
