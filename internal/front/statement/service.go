package statement

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type StatementServiceCreator interface {
	GetAllStatements(req StatementShowRequest) (*StatementListResponse, *responses.ErrorResponse)
	GetStatementsByMemberUUID(req StatementShowRequest, member_uuid string) (*StatementListResponse, *responses.ErrorResponse)
}

type StatementService struct {
	DBPool        *sqlx.DB
	UserContext    *types.UserContext
	StatementRepo StatementRepo
}

func NewStatementService(dbPool *sqlx.DB, userContext *types.UserContext) StatementServiceCreator {
	return &StatementService{
		DBPool:        dbPool,
		UserContext:    userContext,
		StatementRepo:  NewStatementRepoImpl(dbPool, userContext),
	}
}

func (s *StatementService) GetAllStatements(req StatementShowRequest) (*StatementListResponse, *responses.ErrorResponse) {
	return s.StatementRepo.GetAllStatements(req)
}

func (s *StatementService) GetStatementsByMemberUUID(req StatementShowRequest, member_uuid string) (*StatementListResponse, *responses.ErrorResponse) {
	return s.StatementRepo.GetStatementsByMemberUUID(req, member_uuid)
}
