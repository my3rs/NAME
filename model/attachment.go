package model

type Attachment struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`

	ContentID uint    `json:"-" gorm:"default:null"`
	Content   Content `json:"content,omitempty"`

	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime"`
}

func (Attachment) TableName() string {
	return "attachments"
}
