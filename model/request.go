package model

type QueryRequest struct {
	PageSize         int    `url:"pageSize" json:"pageSize"`
	PageIndex        int    `url:"pageIndex" json:"pageIndex"`
	Order            string `url:"orderBy" json:"orderBy,omitempty"`
	WithReadablePath bool   `url:"path" json:"path,omitempty"`
}

type PostRequest struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Abstract     string        `json:"abstract"`
	Text         string        `json:"text"`
	AuthorID     uint          `json:"authorID"`
	CategoryID   uint          `json:"categoryID"`
	Status       ContentStatus `json:"status"`
	CreatedAt    int64         `json:"createdAt"`
	UpdatedAt    int64         `json:"updatedAt"`
	PublishAt    int64         `json:"publishAt"`
	AllowComment bool          `json:"allowComment"`
	Tags         []Tag         `json:"tags"`
}

type DeleteRequest struct {
	IDs []uint
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
