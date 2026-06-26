package auth

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

type AuthHandler struct {
	service     func(c *fiber.Ctx) *AuthService
	userContext func(c *fiber.Ctx) types.UserContext
}

func NewAuthHandler(db_pool *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		service: func(c *fiber.Ctx) *AuthService {
			return NewAuthService(db_pool)
		},
		userContext: func(c *fiber.Ctx) types.UserContext {
			raw := c.Locals("UserContext")
			var uCtx types.UserContext
			if contextMap, ok := raw.(types.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "failed to cast UserContext", "warn")
				uCtx = types.UserContext{}
			}

			return uCtx
		},
	}
}

func (au *AuthHandler) Login(c *fiber.Ctx) error {
	var login_request LoginRequest
	v := utils.NewValidator()

	if err := login_request.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("login_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	loginID := login_request.LoginID
	if loginID == "" {
		loginID = login_request.UserName
	}

	resp, err := au.service(c).Login(loginID, login_request.Password)
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
			utils.Translate("login_success", nil, c),
			1000,
			resp,
		),
	)
}

func (au *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	userCtx := au.userContext(c)
	if userCtx.UserUuid == "" {
		return c.Status(http.StatusUnauthorized).JSON(
			response.NewResponseError(
				utils.Translate("access_denied", nil, c),
				-1000,
				errors.New(utils.Translate("missing_user_context", nil, c)),
			),
		)
	}

	resp, err := au.service(c).GetUserInfo(userCtx.UserUuid)
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
			utils.Translate("get_userinfo_success", nil, c),
			1000,
			resp,
		),
	)
}

func (au *AuthHandler) Register(c *fiber.Ctx) error {
	var register_request RegisterRequest
	v := utils.NewValidator()

	if err := register_request.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("register_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := au.service(c).Register(register_request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(err.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusCreated).JSON(
		response.NewResponse(
			utils.Translate("register_success", nil, c),
			1000,
			resp,
		),
	)
}
