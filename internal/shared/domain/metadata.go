package domain

import "github.com/parxyws/nego-gin/pkg/util"

type Metadata struct {
	PageLimit  int        `json:"page_limit,omitempty"`
	SortOrder  string     `json:"sort_order,omitempty"`
	Pagination Pagination `json:"pagination,omitempty'"`
}

type Pagination struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type CursorService struct {
	Self util.Cursor
	Next util.Cursor
	Prev util.Cursor
}
