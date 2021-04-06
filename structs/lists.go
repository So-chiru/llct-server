package structs

type LLCTGroup struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
	Color   string   `json:"color"`
}

type LLCTSongMetadata struct {
	Released int64 `json:"released"`
}

type LLCTSongs struct {
	Title    string           `json:"title"`
	TitleKo  string           `json:"title.ko"`
	Artist   interface{}      `json:"artist"`
	Metadata LLCTSongMetadata `json:"metadata"`
}

type LLCTLists struct {
	Groups []LLCTGroup   `json:"groups"`
	Songs  [][]LLCTSongs `json:"songs"`
}
