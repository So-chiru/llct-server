package structs

import "github.com/so-chiru/llct-server/dashboard"

type LLCTUpdate struct {
	Updates    int64                 `json:"updates"`
	Notices    []dashboard.Notices   `json:"notices"`
	Dashboards []dashboard.Dashboard `json:"dashboards"`
}
