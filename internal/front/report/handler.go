package report

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type ReportHandler interface {
	GetMemberReports(c *fiber.Ctx) error
}

type ReportHandlerImpl struct {
	service func(c *fiber.Ctx) ReportServiceCreator
	v       *utils.Validator
}

func NewReportHandler(db *sqlx.DB) *ReportHandlerImpl {
	return &ReportHandlerImpl{
		service: func(c *fiber.Ctx) ReportServiceCreator { return NewReportService(db) },
		v:       utils.NewValidator(),
	}
}

func (h *ReportHandlerImpl) GetMemberReports(c *fiber.Ctx) error {
	var req ReportShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_statements_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetMemberReports(req)
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
