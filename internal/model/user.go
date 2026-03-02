package model

type SetGroupRequest struct {
	UserID int64  `json:"-"`
	Group  string `json:"group" validate:"required"`
}

type SetNotificationsSettingRequest struct {
	UserID               int64 `json:"-"`
	NotificationsEnabled bool  `json:"notifications_enabled"`
}
