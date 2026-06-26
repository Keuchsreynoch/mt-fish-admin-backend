package gameconfig

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

type GameConfigHandler interface {
	GetGameConfig(c *fiber.Ctx) error
	UpdateGameConfig(c *fiber.Ctx) error
}

type GameConfigHandlerImpl struct {
	service func(c *fiber.Ctx) GameConfigServiceCreator
	v       *utils.Validator
}

func NewGameConfigHandler(db *sqlx.DB) *GameConfigHandlerImpl {
	return &GameConfigHandlerImpl{
		service: func(c *fiber.Ctx) GameConfigServiceCreator {
			raw := c.Locals("UserContext")
			var uCtx types.UserContext
			if contextMap, ok := raw.(types.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "failed to cast UserContext", "warn")
				uCtx = types.UserContext{}
			}

			return NewGameConfigService(db, &uCtx)
		},
		v:       utils.NewValidator(),
	}
}

func (h *GameConfigHandlerImpl) UpdateGameConfig(c *fiber.Ctx) error {
	var req UpdateGameConfigRequest
	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("update_game_config_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := h.service(c).UpdateGameConfig(req)
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
			utils.Translate("update_game_config_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *GameConfigHandlerImpl) GetGameConfig(c *fiber.Ctx) error {

	resp, err := h.service(c).GetGameConfig()
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
			utils.Translate("get_game_config_success", nil, c),
			1000,
			resp,
		),
	)
}
