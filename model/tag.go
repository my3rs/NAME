package model

type Tag struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;type:bigint" json:"id"`
	ParentID  uint   `json:"parentID"`
	No        string `gorm:"unique;index" json:"no"`
	Text      string `json:"text,omitempty"`
	Path      string `gorm:"type:ltree" json:"path,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime:milli"`

	// 以下是别名字段
	ReadablePath string `gorm:"-" json:"readablePath,omitempty"`
}
