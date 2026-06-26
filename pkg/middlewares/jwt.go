package middlewares

import (
	// "fish_shooting_admin_backend/internal/front/auth"
	"fish_shooting_admin_backend/internal/front/auth"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	response "fish_shooting_admin_backend/pkg/http/response"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func NewJwtMinddleWare(app *fiber.App, DBPool *sqlx.DB) {
	errs := godotenv.Load()
	if errs != nil {
		log.Fatalf("Error loading .env file")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")

	app.Use(func(c *fiber.Ctx) error {
		switch c.Path() {
		case "/api/v1/admin/auth/login", "/api/v1/admin/auth/login/", "/api/v1/front/auth/login", "/api/v1/front/auth/login/":
			return c.Next()
		}

		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			webSocketProtocol := c.Get("Sec-WebSocket-Protocol")
			if webSocketProtocol == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": utils.Translate("missing_ws_auth_protocol", nil, c),
				})
			}

			parts := strings.Split(webSocketProtocol, ",")
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": utils.Translate("invalid_ws_auth_format", nil, c),
				})
			}

			tokenString := strings.TrimSpace(parts[1])
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret_key), nil
			})
			if err != nil || !token.Valid {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": utils.Translate("invalid_jwt_token", nil, c),
				})
			}
			c.Locals("jwt_data", token)
			c.Set("Sec-WebSocket-Protocol", "Bearer")
			return c.Next()
		}

		return jwtware.New(jwtware.Config{
			SigningKey:  jwtware.SigningKey{Key: []byte(secret_key)},
			ContextKey: "jwt_data",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": utils.Translate("session_expired", nil, c),
					})
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": utils.Translate("missing_or_malformed_jwt", nil, c),
				})
			},
		})(c)
	})

	app.Use(func(c *fiber.Ctx) error {
		switch c.Path() {
		case "/api/v1/admin/auth/login", "/api/v1/admin/auth/login/", "/api/v1/front/auth/login", "/api/v1/front/auth/login/":
			return c.Next()
		}

		user_token := c.Locals("jwt_data").(*jwt.Token)
		pclaim := user_token.Claims.(jwt.MapClaims)

		if websocketUpgrade := c.Get("Upgrade"); websocketUpgrade == "websocket" {
			return handleUserContext(c, pclaim, DBPool)
		}

		return handleUserContext(c, pclaim, DBPool)
	})
}

func handleUserContext(c *fiber.Ctx, pclaim jwt.MapClaims, DBPool *sqlx.DB) error {
	user_uuid, ok := pclaim["user_uuid"].(string)
	if !ok || user_uuid == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			utils.Translate("access_denied", nil, c),
			-500,
			fmt.Errorf(utils.Translate("missing_or_invalid_key", map[string]interface{}{"key": "user_uuid"}, c)),
		))
	}

	login_session, ok := pclaim["login_session"].(string)
	if !ok || login_session == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			utils.Translate("access_denied", nil, c),
			-500,
			fmt.Errorf(utils.Translate("missing_or_invalid_key", map[string]interface{}{"key": "login_session"}, c)),
		))
	}

	exp, ok := pclaim["exp"].(float64)
	if !ok {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			utils.Translate("access_denied", nil, c),
			-500,
			fmt.Errorf(utils.Translate("missing_or_invalid_key", map[string]interface{}{"key": "exp"}, c)),
		))
	}

	// get member info
	user_info, err := auth.NewAuthRepoImpl(DBPool).GetUserByUUID(user_uuid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			utils.Translate("access_denied", nil, c),
			-500,
			fmt.Errorf(utils.Translate("get_userinfo_failed", nil, c)),
		))
	}

	// validate session
	if login_session != user_info.LoginSession {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.NewResponseError(
			utils.Translate("access_denied", nil, c),
			-500,
			fmt.Errorf(utils.Translate("session_expired", nil, c)),
		))
	}

	user_context := types.UserContext{
		Id:           user_info.ID,
		UserUuid:     user_info.UserUUID,
		UserName:     user_info.UserName,
		LoginSession: login_session,
		Exp:          time.Unix(int64(exp), 0),
		UserAgent:    string(c.Context().UserAgent()),
		Ip:           string(c.Context().RemoteIP().String()),
		StatusId:     user_info.StatusID,
	}
	c.Locals("UserContext", user_context)

	return c.Next()
}
