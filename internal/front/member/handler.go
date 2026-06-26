package member

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type MemberHandler struct {
	DBPool        *sqlx.DB
	MemberService *MemberService
}

func NewMemberHandler(dbPool *sqlx.DB) *MemberHandler {
	return &MemberHandler{
		DBPool:        dbPool,
		MemberService: NewMemberService(dbPool),
	}
}

func (m *MemberHandler) GetAllMembers(c *fiber.Ctx) error {
	req := MemberShowRequest{}
	v := utils.NewValidator()

	if err := req.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_all_members_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := m.MemberService.GetAllMembers(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponseWithPaging(
			utils.Translate("get_all_members_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}
