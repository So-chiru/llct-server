package structs

type LLCTGroup struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
	Color   string   `json:"color"`
}

type LLCTSongs struct {
	Title   string `json:"title"`
	TitleKo string `json:"title.ko"`
	Artist  int    `json:"artist"`
}

type LLCTLists struct {
	Groups []LLCTGroup   `json:"groups"`
	Songs  [][]LLCTSongs `json:"songs"`
}
