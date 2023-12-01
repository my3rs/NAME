package model

const (
	CommentStatus_Approved   = 0 // 正常评论
	CommentStatus_Unreviewed = 1 // 等待审核
	CommentStatus_Refused    = 2 // 审核未通过
	CommentStatus_Trash      = 3 // 放入垃圾箱
)

type Comment struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ContentID uint    `json:"-" gorm:"default:null"`
	Content   Content `json:"-"`

	CreatedAt int `json:"createdAt"  gorm:"autoCreateTime:milli"`

	AuthorID   uint   `json:"-" gorm:"default:null"`
	AuthorName string `json:"authorName" gorm:"default:null"`

	ParentID uint `json:"parentID" gorm:"default:null"`

	Path   string `json:"path"`
	Mail   string `json:"mail"`
	URL    string `json:"url"`
	Text   string `json:"text"`
	Status uint   `json:"status"`
	IP     string `json:"-"`
	Agent  string `json:"agent"`

	Points uint `json:"points"`
}
