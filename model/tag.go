package model

// Tag represents a tag entity
type Tag struct {
	// ID is the unique identifier for the tag
	ID uint `gorm:"primaryKey;autoIncrement;type:bigint;comment:标签ID" json:"id"`
	// No is the unique identifier for the tag, used in URLs
	No string `gorm:"unique;index;comment:标签唯一标识符，用于URL" json:"no"`
	// Text is the display text for the tag
	Text string `json:"text,omitempty" gorm:"comment:标签显示文本"`
	// CountCount is the number of times the tag has been used
	CountCount int `json:"countCount" gorm:"comment:标签被使用的次数"`
	// CreatedAt is the creation time of the tag (in milliseconds)
	CreatedAt int64 `json:"createdAt,omitempty" gorm:"autoCreateTime;comment:创建时间（毫秒时间戳）"`
	// UpdatedAt is the last update time of the tag (in milliseconds)
	UpdatedAt int64 `json:"updatedAt,omitempty" gorm:"autoUpdateTime;comment:更新时间（毫秒时间戳）"`
}
