package admin

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type AdminService interface {
	GetDashboard() (*DashboardSummaryResponse, *responses.ErrorResponse)
	ListUsers() (*UserListResponse, *responses.ErrorResponse)
	ListMenus() (*MenuListResponse, *responses.ErrorResponse)
	GetUserInfo(userUUID string) (*UserInfoResponse, *responses.ErrorResponse)
	GetUserMenus(userID int) (*UserMenusResponse, *responses.ErrorResponse)
	CreateUser(req NewUserRequest) (*NewUser, *responses.ErrorResponse)
	UpdateUser(userID int, req UpdateAdminUserRequest) (*UpdateAdminUserResponse, *responses.ErrorResponse)
	AssignMenus(userID int, req AssignMenusRequest) (*UserMenusResponse, *responses.ErrorResponse)
}

type AdminServiceImpl struct {
	DBPool    *sqlx.DB
	AdminRepo AdminRepo
	UserContext *types.UserContext
}

func NewAdminService(dbPool *sqlx.DB, uCtx *types.UserContext) *AdminServiceImpl {
	return &AdminServiceImpl{
		DBPool:    dbPool,
		AdminRepo: NewAdminRepoImpl(dbPool, uCtx),
		UserContext: uCtx,
	}
}

func (s *AdminServiceImpl) ListUsers() (*UserListResponse, *responses.ErrorResponse) {
	return s.AdminRepo.ListUsers()
}

func (s *AdminServiceImpl) GetDashboard() (*DashboardSummaryResponse, *responses.ErrorResponse) {
	return s.AdminRepo.GetDashboard()
}

func (s *AdminServiceImpl) ListMenus() (*MenuListResponse, *responses.ErrorResponse) {
	return s.AdminRepo.ListMenus()
}

func (s *AdminServiceImpl) GetUserInfo() (*UserInfoResponse, *responses.ErrorResponse) {
	return s.AdminRepo.GetUserInfo(s.UserContext.UserUuid)
}

func (s *AdminServiceImpl) GetUserMenus(userID int) (*UserMenusResponse, *responses.ErrorResponse) {
	return s.AdminRepo.GetUserMenus(userID)
}

func (s *AdminServiceImpl) GetMyMenus() (*UserMenusResponse, *responses.ErrorResponse) {
	return s.AdminRepo.GetUserMenus(s.UserContext.Id)
}

func (s *AdminServiceImpl) CreateUser(req NewUserRequest) (*NewUser, *responses.ErrorResponse) {
	return s.AdminRepo.CreateUser(req)
}

func (s *AdminServiceImpl) UpdateUser(userID int, req UpdateAdminUserRequest) (*UpdateAdminUserResponse, *responses.ErrorResponse) {
	return s.AdminRepo.UpdateUser(userID, req)
}

func (s *AdminServiceImpl) AssignMenus(userID int, req AssignMenusRequest) (*UserMenusResponse, *responses.ErrorResponse) {
	return s.AdminRepo.AssignMenus(userID, req)
}
