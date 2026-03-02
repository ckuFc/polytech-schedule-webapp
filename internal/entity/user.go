package entity

type User struct {
	ID                   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	TelegramID           int64  `gorm:"column:telegram_id;unique"`
	UserName             string `gorm:"column:username"`
	FirstName            string `gorm:"column:firstname"`
	LastName             string `gorm:"column:lastname"`
	Group                string `gorm:"column:group"`
	PhotoURL             string `gorm:"column:photo_url"`
	NotificationsEnabled bool   `gorm:"column:notifications_enabled;default:true"`
	CreatedAt            int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt            int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}
