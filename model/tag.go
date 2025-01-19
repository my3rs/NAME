package model

// Tag represents a tag entity
type Tag struct {
	ID         uint   `gorm:"primaryKey;comment:标签ID" json:"id"`
	Slug       string `gorm:"unique;index;comment:标签唯一标识符，用于URL" json:"slug"`
	Text       string `json:"text,omitempty" gorm:"comment:标签显示文本"`
	CountCount int    `json:"countCount" gorm:"comment:标签被使用的次数"`
	CreatedAt  int64  `json:"createdAt,omitempty" gorm:"autoCreateTime;comment:创建时间（毫秒时间戳）"`
	UpdatedAt  int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime;comment:更新时间（毫秒时间戳）"`
}
