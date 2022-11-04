package model

type Tag struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	No        string `gorm:"unique;index" json:"no"`
	Text      string `json:"text"`
	Path      string `json:"path"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
}

type TagTreeNode struct {
	Tag      Tag
	Children []Tag
}
