package model

type CommentStatus uint

const (
	CommentStatusApproved   CommentStatus = 0 // 正常评论
	CommentStatusUnreviewed CommentStatus = 1 // 等待审核
	CommentStatusRefused    CommentStatus = 2 // 审核未通过
	CommentStatusTrash      CommentStatus = 3 // 放入垃圾箱
)

type Comment struct {
	ID uint `json:"id" gorm:"primaryKey;comment:评论ID"`

	// belongs to 关系
	ContentID uint    `json:"-" gorm:"default:null;comment:关联的内容ID"`
	Content   Content `json:"content,omitempty" gorm:"default:null,OnDelete:CASCADE;comment:关联的内容"`

	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime;comment:创建时间（毫秒时间戳）"`

	AuthorID   uint   `json:"-" gorm:"default:null;comment:评论作者ID"`
	AuthorName string `json:"authorName" gorm:"default:null;comment:评论作者名称"`

	ParentID uint   `json:"parentID" gorm:"default:null;comment:父评论ID，用于回复功能"`
	Path     string `json:"-" gorm:"comment:存储从根评论到当前评论的路径，格式如 1.2.3"`

	Mail   string        `json:"mail" gorm:"comment:评论者邮箱"`
	URL    string        `json:"url" gorm:"comment:评论者网站"`
	Text   string        `json:"text" gorm:"comment:评论内容"`
	Status CommentStatus `json:"status" gorm:"default:0;comment:评论状态：0正常、1待审核、2拒绝、3垃圾箱"`
	IP     string        `json:"-" gorm:"comment:评论者IP地址"`
	Agent  string        `json:"agent" gorm:"comment:评论者浏览器User-Agent"`

	Points uint `json:"points" gorm:"comment:评论获得的点数"`

	Children []Comment `json:"children,omitempty" gorm:"-"` // 子评论列表，不存储在数据库中
}

func (Comment) TableName() string {
	return "comments"
}
