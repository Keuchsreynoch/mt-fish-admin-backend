package handler

import (
	"fish_shooting_admin_backend/internal/front/admin"
	"fish_shooting_admin_backend/internal/front/auth"
	"fish_shooting_admin_backend/internal/front/cointransaction"
	gameconfig "fish_shooting_admin_backend/internal/front/game_config"
	fishtype "fish_shooting_admin_backend/internal/front/fish_type"
	"fish_shooting_admin_backend/internal/front/jackpot"
	jackpothistory "fish_shooting_admin_backend/internal/front/jackpot_history"
	"fish_shooting_admin_backend/internal/front/member"
	memberbonus "fish_shooting_admin_backend/internal/front/member_bonus"
	"fish_shooting_admin_backend/internal/front/report"
	"fish_shooting_admin_backend/internal/front/statement"
	"fish_shooting_admin_backend/internal/front/websocket"
	"fish_shooting_admin_backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// group all the module factories
type ServiceHandler struct {
	Front *FrontService
}

// register modules route 
type FrontService struct {
	AuthRoute            *auth.AuthRoute
	AdminRoute           *admin.AdminRoute
	MemberRoute          *member.MemberRoute
	StatementRoute       *statement.StatementRoute
	ReportRoute          *report.ReportRoute
	CoinTransactionRoute *cointransaction.CoinTransactionRoute
	GameConfigRoute      *gameconfig.GameConfigRoute
	JackpotRoute         *jackpot.JackpotRoute
	JackpotHistoryRoute  *jackpothistory.JackpotHistoryRoute
	WebSocketRoute            *websocket.WebSocketRoute
	TransitionRoute           *cointransaction.CoinTransactionRoute
	FishTypeRoute             *fishtype.FishTypeRoute
	MemberBonusRoute *memberbonus.JackpotMemberBonusRoute
}

func NewFrontService(app *fiber.App, db *sqlx.DB) *FrontService {

	au := auth.NewRoute(db, app).RegisterAuthRoute()
	middlewares.NewJwtMinddleWare(app, db)
	adminRoute := admin.NewAdminRoute(db, app).RegisterAdminRoute()
	memberRoute := member.NewMemberRoute(app, db).RegisterMemberRoute()
	statementRoute := statement.NewStatementRoute(app, db).RegisterStatementRoute()
	reportRoute := report.NewReportRoute(app, db).RegisterReportRoute()
	coinTransactionRoute := cointransaction.NewCoinTransactionRoute(app, db).RegisterCoinTransactionRoute()
	gameConfigRoute := gameconfig.NewGameConfigRoute(app, db).RegisterGameConfigRoute()
	jackpotRoute := jackpot.NewJackpotRoute(app, db).RegisterJackpotRoute()
	jackpotHistoryRoute := jackpothistory.NewJackpotHistoryRoute(app, db).RegisterJackpotHistoryRoute()
	web_socket := websocket.NewRoute(app, db).RegisterWebSocketRoute()
	fish_type := fishtype.NewFishTypeRoute(app, db).RegisterFishTypeRoute()
	member_bonus := memberbonus.NewJackpotMemberBonusRoute(app, db).RegisterJackpotMemberBonusRoute()

	return &FrontService{
		AuthRoute:            au,
		AdminRoute:           adminRoute,
		MemberRoute:          memberRoute,
		StatementRoute:       statementRoute,
		ReportRoute:          reportRoute,
		CoinTransactionRoute: coinTransactionRoute,
		GameConfigRoute:      gameConfigRoute,
		JackpotRoute:         jackpotRoute,
		JackpotHistoryRoute:  jackpotHistoryRoute,
		WebSocketRoute:            web_socket,
		FishTypeRoute:             fish_type,
		MemberBonusRoute: member_bonus,
	}
}

func NewServiceHandlers(app *fiber.App, pool *sqlx.DB) *ServiceHandler {
	return &ServiceHandler{
		Front: NewFrontService(app, pool),
	}
}
