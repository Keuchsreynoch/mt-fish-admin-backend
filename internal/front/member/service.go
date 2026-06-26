package member

import (
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type MemberServiceCreator interface {
	GetAllMembers(req MemberShowRequest) (*GetAllMembersResponse, *responses.ErrorResponse)
}

type MemberService struct {
	DBPool     *sqlx.DB
	MemberRepo *MemberRepoImpl
}

func NewMemberService(dbPool *sqlx.DB) *MemberService {
	return &MemberService{
		DBPool:     dbPool,
		MemberRepo: NewMemberRepoImpl(dbPool),
	}
}

func (m *MemberService) GetAllMembers(req MemberShowRequest) (*GetAllMembersResponse, *responses.ErrorResponse) {
	return m.MemberRepo.GetAllMembers(req)
}
