package model

type File struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Path string `json:"path"`

	ContentID uint    `json:"-"`
	Content   Content `json:"content"`

	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime:milli"`
}
