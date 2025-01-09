package model

type Attachment struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`

	ContentID uint    `json:"-" gorm:"default:null"`
	Content   Content `json:"content"`

	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime"`
}
