package structs

type LLCTGroup struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
	Color   string   `json:"color"`
}

type LLCTSongMetadata struct {
	Album    string   `json:"album,omitempty"`
	Length   int64    `json:"length,omitempty"`
	BPM      int64    `json:"bpm,omitempty"`
	Released int64    `json:"released,omitempty"`
	Composer []string `json:"composer,omitempty"`
}

type LLCTSongs struct {
	Title    string            `json:"title"`
	TitleKo  *string           `json:"title.ko,omitempty"`
	Artist   interface{}       `json:"artist"`
	Metadata *LLCTSongMetadata `json:"metadata,omitempty"`
}

type LLCTLists struct {
	Groups []LLCTGroup   `json:"groups"`
	Songs  [][]LLCTSongs `json:"songs"`
}
