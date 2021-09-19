package structs

import "github.com/so-chiru/llct-server/dashboard"

type LLCTUpdate struct {
	UpdateAt  int64                 `json:"updateAt"`
	Notices   []dashboard.Notices   `json:"notices"`
	Dashboard []dashboard.Dashboard `json:"dashboard"`
}
