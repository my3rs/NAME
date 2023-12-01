package model

// type Tag struct {
// 	ID        uint   `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
// 	No        string `gorm:"unique;index" json:"no"`
// 	Text      string `json:"text,omitempty"`
// 	Path      string `json:"path,omitempty"`
// 	CreatedAt int64  `json:"createdAt,omitempty" gorm:"autoCreateTime:milli"`
// 	UpdatedAt int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime:milli"`
// }

type Tag struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID  uint   `json:"parentID"`
	No        string `gorm:"unique;index" json:"no"`
	Text      string `json:"text,omitempty"`
	Path      string `json:"path,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime:milli"`
}
