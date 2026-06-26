package memberbonus

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type JackpotMemberBonusHandler interface {
	GetAllMemberBonuses(c *fiber.Ctx) error
	CreateMemberBonus(c *fiber.Ctx) error
}

type JackpotMemberBonusHandlerImpl struct {
	service func(c *fiber.Ctx) JackpotMemberBonusServiceCreator
	v       *utils.Validator
}

func NewJackpotMemberBonusHandler(db *sqlx.DB) *JackpotMemberBonusHandlerImpl {
	return &JackpotMemberBonusHandlerImpl{
		service: func(c *fiber.Ctx) JackpotMemberBonusServiceCreator {
			raw := c.Locals("UserContext")
			var uCtx types.UserContext
			if contextMap, ok := raw.(types.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "failed to cast UserContext", "warn")
				uCtx = types.UserContext{}
			}

			return NewJackpotMemberBonusService(db, &uCtx)
		},
		v:       utils.NewValidator(),
	}
}

func (h *JackpotMemberBonusHandlerImpl) GetAllMemberBonuses(c *fiber.Ctx) error {
	var req JackpotMemberBonusShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_jackpot_member_bonuses_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetAllMemberBonuses(req)
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
			utils.Translate("get_jackpot_member_bonuses_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}

func (h *JackpotMemberBonusHandlerImpl) CreateMemberBonus(c *fiber.Ctx) error {
	var req CreateJackpotMemberBonusRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("create_jackpot_member_bonus_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).CreateMemberBonus(req)
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
		response.NewResponse(
			utils.Translate("create_jackpot_member_bonus_success", nil, c),
			1000,
			resp,
		),
	)
}
