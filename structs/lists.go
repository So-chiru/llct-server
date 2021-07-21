package structs

type LLCTGroup struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
	Color   string   `json:"color"`
}

type LLCTSongMetadata struct {
	Album    int64    `json:"album"`
	Length   int64    `json:"length"`
	BPM      int64    `json:"bpm"`
	Released int64    `json:"released"`
	Composer []string `json:"composer"`
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
