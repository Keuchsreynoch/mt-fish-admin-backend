package websocket

import "time"

const (
	TopicFrontNotificationCreated = "front.notification.created"
	TopicFrontCoinUpdate = "front.coin.updated"
)

type FrontNotificationPayload struct {
	ID                 int       `json:"id"`
	MemberID           int       `json:"member_id"`
	NotificationTypeID int       `json:"notification_type_id"`
	Context            string    `json:"context"`
	Subject            string    `json:"subject"`
	Description        string    `json:"description"`
	StatusID           int16     `json:"status_id"`
	CreatedAt          time.Time `json:"created_at"`
}

type FrontCoinUpdatePayload struct {
	MemberID   int    `json:"member_id"`
	CoinAmount string `json:"coin_amount"`
}
