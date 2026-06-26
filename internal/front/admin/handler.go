package admin

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	types "fish_shooting_admin_backend/pkg/model"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AdminHandler interface {
	GetDashboard(c *fiber.Ctx) error
	ListUsers(c *fiber.Ctx) error
	ListMenus(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	GetMyMenus(c *fiber.Ctx) error
	GetUserMenus(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	AssignMenus(c *fiber.Ctx) error
}

type AdminHandlerImpl struct {
	service func(c *fiber.Ctx) *AdminServiceImpl
	v       *utils.Validator
}

func NewAdminHandler(dbPool *sqlx.DB) *AdminHandlerImpl {
	return &AdminHandlerImpl{
		service: func(c *fiber.Ctx) *AdminServiceImpl {
			userContext := c.Locals("UserContext")
			var uCtx types.UserContext
			if contextMap, ok := userContext.(types.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "failed to cast UserContext", "warn")
				uCtx = types.UserContext{}
			}

			return NewAdminService(dbPool, &uCtx)
		},
		v: utils.NewValidator(),
	}
}

func (h *AdminHandlerImpl) GetDashboard(c *fiber.Ctx) error {
	resp, errResp := h.service(c).GetDashboard()
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("get_dashboard_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) ListUsers(c *fiber.Ctx) error {
	resp, errResp := h.service(c).ListUsers()
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("get_admin_users_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) ListMenus(c *fiber.Ctx) error {
	resp, errResp := h.service(c).ListMenus()
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("get_admin_menus_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) GetUserInfo(c *fiber.Ctx) error {

	resp, errResp := h.service(c).GetUserInfo()
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
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

func (h *AdminHandlerImpl) GetMyMenus(c *fiber.Ctx) error {
	resp, errResp := h.service(c).GetMyMenus()
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("get_my_menus_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) GetUserMenus(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("get_user_menus_failed", nil, c),
				-1000,
				errors.New(utils.Translate("invalid_user_id", nil, c)),
			),
		)
	}

	resp, errResp := h.service(c).GetUserMenus(userID)
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("get_user_menus_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) CreateUser(c *fiber.Ctx) error {
	req := NewUserRequest{}

	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("create_admin_user_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, errResp := h.service(c).CreateUser(req)
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusCreated).JSON(
		response.NewResponse(
			utils.Translate("create_admin_user_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	req := UpdateAdminUserRequest{}

	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("update_admin_user_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("update_admin_user_failed", nil, c),
				-1000,
				errors.New(utils.Translate("invalid_user_id", nil, c)),
			),
		)
	}

	resp, errResp := h.service(c).UpdateUser(userID, req)
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("update_admin_user_success", nil, c),
			1000,
			resp,
		),
	)
}

func (h *AdminHandlerImpl) AssignMenus(c *fiber.Ctx) error {
	req := AssignMenusRequest{}

	if err := req.Bind(c, h.v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("assign_menus_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate("assign_menus_failed", nil, c),
				-1000,
				errors.New(utils.Translate("invalid_user_id", nil, c)),
			),
		)
	}


	resp, errResp := h.service(c).AssignMenus(userID, req)
	if errResp != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utils.Translate(errResp.MessageID, nil, c),
				-1000,
				errors.New(utils.Translate(errResp.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utils.Translate("assign_menus_success", nil, c),
			1000,
			resp,
		),
	)
}
