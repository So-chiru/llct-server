package dashboard

type LLCTUpdate struct {
	Updates    int64       `json:"updates"`
	Notices    []Notices   `json:"notices"`
	Dashboards []Dashboard `json:"dashboards"`
}
