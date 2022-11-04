package model

type QueryRequest struct {
	PageSize  int    `url:"pageSize" json:"pageSize"`
	PageIndex int    `url:"pageIndex" json:"pageIndex"`
	Order     string `url:"orderBy" json:"orderBy,omitempty"`
}

type PostRequest struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Abstract     string        `json:"abstract"`
	Text         string        `json:"text"`
	AuthorID     uint          `json:"authorID"`
	Status       ContentStatus `json:"status"`
	PublishAt    int64         `json:"publishAt"`
	IsPublic     bool          `json:"isPublic"`
	AllowComment bool          `json:"allowComment"`
}

type DeleteRequest struct {
	IDs []uint
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
