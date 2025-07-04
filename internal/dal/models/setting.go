package models

type UserSetting struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex;not null"` // 配置项名称，如：theme / auto_accept_match
	Value     string `gorm:"not null"`             // 配置值（建议都存为字符串）
	UpdatedAt int64  // 时间戳（可选）
}
