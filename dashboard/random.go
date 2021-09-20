package dashboard

import (
	"github.com/so-chiru/llct-server/utils"
)

func generateRandomSongs(size int) *[]Dashboard {
	songs := utils.GenerateRandomSongs(size)
	var musics []MusicComponent = make([]MusicComponent, 0)

	for _, songId := range songs {
		var reason = "랜덤"
		var music = MusicComponent{
			ID:              songId,
			RecommendReason: &reason,
		}

		musics = append(musics, music)
	}

	var results []Dashboard = make([]Dashboard, 0)

	var title = "오늘의 랜덤 곡을 뽑아 볼까요?"
	var dash = Dashboard{
		Type:     "musicset",
		Title:    &title,
		MusicSet: &musics,
	}
	results = append(results, dash)

	return &results
}
