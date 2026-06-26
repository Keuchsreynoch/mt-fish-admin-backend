package fishtype

import (
	"errors"
	response "fish_shooting_admin_backend/pkg/http/response"
	"fish_shooting_admin_backend/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type FishTypeHandler interface {
	GetFishTypes(c *fiber.Ctx) error
}

type FishTypeHandlerImpl struct {
	service FishTypeServiceCreator
}

func NewFishTypeHandler(db *sqlx.DB) *FishTypeHandlerImpl {
	return &FishTypeHandlerImpl{
		service: NewFishTypeService(db),
	}
}

func (h *FishTypeHandlerImpl) GetFishTypes(c *fiber.Ctx) error {
	resp, err := h.service.GetFishTypes()
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
			utils.Translate("get_fish_types_success", nil, c),
			1000,
			FishTypesResponse{FishTypes: resp},
			1,
			len(resp),
			len(resp),
		),
	)
}
