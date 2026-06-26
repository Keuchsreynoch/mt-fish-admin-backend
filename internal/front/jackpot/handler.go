package jackpot

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

type JackpotHandler interface {
	GetCurrentJackpotGlobal(c *fiber.Ctx) error
	GetJackpotGlobalUpdateForm(c *fiber.Ctx) error
	GetAllCompanyTopups(c *fiber.Ctx) error
	CreateCompanyTopup(c *fiber.Ctx) error
	GetAllLedgers(c *fiber.Ctx) error
	UpdateJackpotRate(c *fiber.Ctx) error
	UpdateJackpotGlobal(c *fiber.Ctx) error
}

type JackpotHandlerImpl struct {
	service func(c *fiber.Ctx) JackpotServiceCreator
	v       *utils.Validator
}

func NewJackpotHandler(db *sqlx.DB) *JackpotHandlerImpl {
	return &JackpotHandlerImpl{
		service: func(c *fiber.Ctx) JackpotServiceCreator {
			raw := c.Locals("UserContext")
			var uCtx types.UserContext
			if contextMap, ok := raw.(types.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "failed to cast UserContext", "warn")
				uCtx = types.UserContext{}
			}

			return NewJackpotService(db, &uCtx)
		},
		v:       utils.NewValidator(),
	}
}

func (h *JackpotHandlerImpl) GetCurrentJackpotGlobal(c *fiber.Ctx) error {
	resp, err := h.service(c).GetCurrentJackpotGlobal()
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
			utils.Translate("get_jackpot_global_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *JackpotHandlerImpl) GetJackpotGlobalUpdateForm(c *fiber.Ctx) error {
	resp, err := h.service(c).GetJackpotGlobalUpdateForm()
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
			utils.Translate("get_jackpot_global_form_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *JackpotHandlerImpl) GetAllCompanyTopups(c *fiber.Ctx) error {
	var req JackpotCompanyTopupShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_jackpot_company_topups_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetAllCompanyTopups(req)
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
			utils.Translate("get_jackpot_company_topups_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}

func (h *JackpotHandlerImpl) CreateCompanyTopup(c *fiber.Ctx) error {
	var req CreateJackpotCompanyTopupRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("create_jackpot_company_topup_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).CreateCompanyTopup(req)
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
			utils.Translate("create_jackpot_company_topup_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *JackpotHandlerImpl) GetAllLedgers(c *fiber.Ctx) error {
	var req JackpotLedgerShowRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_jackpot_ledger_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).GetAllLedgers(req)
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
			utils.Translate("get_jackpot_ledger_success", nil, c),
			1000,
			resp,
			req.PageOptions.Page,
			req.PageOptions.Perpage,
			resp.Total,
		),
	)
}

func (h *JackpotHandlerImpl) UpdateJackpotRate(c *fiber.Ctx) error {
	var req UpdateJackpotRateRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("update_jackpot_rate_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).UpdateJackpotRate(req)
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
			utils.Translate("update_jackpot_rate_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *JackpotHandlerImpl) UpdateJackpotGlobal(c *fiber.Ctx) error {
	var req UpdateJackpotGlobalRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("update_jackpot_global_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).UpdateJackpotGlobal(req)
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
			utils.Translate("update_jackpot_global_success", nil, c),
			1000,
			resp,
		),
	)
}
