package jackpothistory

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type JackpotHistoryHandler interface {
	GetAllJackpotHistories(c *fiber.Ctx) error
}

type JackpotHistoryHandlerImpl struct {
	service func(c *fiber.Ctx) JackpotHistoryServiceCreator
	v       *utils.Validator
}

func NewJackpotHistoryHandler(db *sqlx.DB) *JackpotHistoryHandlerImpl {
	return &JackpotHistoryHandlerImpl{
		service: func(c *fiber.Ctx) JackpotHistoryServiceCreator { return NewJackpotHistoryService(db, nil) },
		v:       utils.NewValidator(),
	}
}

func (h *JackpotHistoryHandlerImpl) GetAllJackpotHistories(c *fiber.Ctx) error {
	var req JackpotHistoryShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_jackpot_history_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetAllJackpotHistories(req)
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
			utils.Translate("get_jackpot_history_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}
