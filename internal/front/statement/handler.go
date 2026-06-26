package statement

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type StatementHandler interface {
	GetAllStatements(c *fiber.Ctx) error
	GetStatementsByMemberUUID(c *fiber.Ctx) error
}

type StatementHandlerImpl struct {
	service func(c *fiber.Ctx) StatementServiceCreator
	v       *utils.Validator
}

func NewStatementHandler(db *sqlx.DB) *StatementHandlerImpl {
	return &StatementHandlerImpl{
		service: func(c *fiber.Ctx) StatementServiceCreator { return NewStatementService(db, nil) },
		v: utils.NewValidator(),
	}
}

func (h *StatementHandlerImpl) GetAllStatements(c *fiber.Ctx) error {
	var req StatementShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_statements_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetAllStatements(req)
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
			utils.Translate("get_statements_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}

func (h *StatementHandlerImpl) GetStatementsByMemberUUID(c *fiber.Ctx) error {
	var req StatementShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_statements_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	memberUUID := c.Params("member_uuid")
	if memberUUID == "" {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_statements_failed", nil, c),
				-1000,
				errors.New("missing member_uuid"),
			),
		)
	}

	resp, err := h.service(c).GetStatementsByMemberUUID(req, memberUUID)
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
			utils.Translate("get_statements_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}
