package model

type QueryRequest struct {
	PageSize         int    `url:"pageSize" json:"pageSize"`
	PageIndex        int    `url:"pageIndex" json:"pageIndex"`
	Order            string `url:"orderBy" json:"orderBy,omitempty"`
	WithReadablePath bool   `url:"path" json:"path,omitempty"`
}
