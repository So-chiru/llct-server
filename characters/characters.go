package characters

import "fmt"

type CharacterData struct {
	Color  string
	Musics []string
}

var CharacterDataCollection = map[string]CharacterData{
	"코사카 호노카": {
		Color:  "#FFA400",
		Musics: []string{"015", "016", "074", "044"},
	},
	"아야세 에리": {
		Color:  "#41B6E6",
		Musics: []string{"028", "031", "058"},
	},
	"미나미 코토리": {
		Color:  "#B2B4B2",
		Musics: []string{"014", "025", "048"},
	},
	"소노다 우미": {
		Color:  "#003DA5",
		Musics: []string{"013", "048"},
	},
	"호시조라 린": {
		Color:  "#FEDD00",
		Musics: []string{"050", "072", "078"},
	},
	"니시키노 마키": {
		Color:  "#EE2737",
		Musics: []string{"050", "077"},
	},
	"토죠 노조미": {
		Color:  "#84329B",
		Musics: []string{"058", "083"},
	},
	"코이즈미 하나요": {
		Color:  "#00AB84",
		Musics: []string{"049", "025"},
	},
	"야자와 니코": {
		Color:  "#E31C79",
		Musics: []string{"055", "077"},
	},
	"타카미 치카": {
		Color:  "#FF7F32",
		Musics: []string{"159", "1113"},
	},
	"사쿠라우치 리코": {
		Color:  "#FB637E",
		Musics: []string{"164", "1118"},
	},
	"마츠우라 카난": {
		Color:  "#00C7B1",
		Musics: []string{"169", "1126"},
	},
	"쿠로사와 다이아": {
		Color:  "#E4002B",
		Musics: []string{"167", "1125"},
	},
	"와타나베 요우": {
		Color:  "#00B5E2",
		Musics: []string{"165", "1130"},
	},
	"츠시마 요시코": {
		Color:  "#B1B3B3",
		Musics: []string{"163", "1132"},
	},
	"쿠니키다 하나마루": {
		Color:  "#FFCD00",
		Musics: []string{"162", "1127"},
	},
	"오하라 마리": {
		Color:  "#9B26B6",
		Musics: []string{"168", "1131"},
	},
	"쿠로사와 루비": {
		Color:  "#E93CAC",
		Musics: []string{"166", "1119"},
	},
	"우에하라 아유무": {
		Color:  "#ED7D95",
		Musics: []string{"22", "212", "229", "245"},
	},
	"나카스 카스미": {
		Color:  "#E7D600",
		Musics: []string{"23", "213", "230", "246"},
	},
	"오사카 시즈쿠": {
		Color:  "#3FA4C6",
		Musics: []string{"24", "214", "231", "247"},
	},
	"아사카 카린": {
		Color:  "#495EC6",
		Musics: []string{"25", "215", "232", "248"},
	},
	"미야시타 아이": {
		Color:  "#FF5800",
		Musics: []string{"26", "216", "233", "249"},
	},
	"코노에 카나타": {
		Color:  "#B365AE",
		Musics: []string{"27", "217", "234", "250"},
	},
	"유키 세츠나": {
		Color:  "#D81C2F",
		Musics: []string{"28", "218", "235", "251"},
	},
	"엠마 베르데": {
		Color:  "#8EC225",
		Musics: []string{"29", "219", "236", "252"},
	},
	"텐노지 리나": {
		Color:  "#969FB5",
		Musics: []string{"210", "220", "237", "253"},
	},
	"미후네 시오리코": {
		Color:  "#36B482",
		Musics: []string{"238"},
	},
	"미아 테일러": {
		Color:  "#A9A89A",
		Musics: []string{},
	},
	"쇼우 란쥬": {
		Color:  "#F69992",
		Musics: []string{},
	},
	"카즈노 세이라": {
		Color:  "#ACC8EC",
		Musics: []string{"127", "157", "158"},
	},
	"카즈노 리아": {
		Color:  "#DEE6EC",
		Musics: []string{"127", "157", "158"},
	},
	"시부야 카논": {
		Color:  "#FF7E24",
		Musics: []string{"31", "39", "311"},
	},
	"탕 쿠쿠": {
		Color:  "#A0FFF6",
		Musics: []string{"31", "311"},
	},
	"아라시 치사토": {
		Color:  "#FF6E92",
		Musics: []string{"31"},
	},
	"헤안나 스미레": {
		Color:  "#74F467",
		Musics: []string{"31"},
	},
	"하즈키 렌": {
		Color:  "#0100A0",
		Musics: []string{"31"},
	},
}

func GetCharacterData(c string) (*CharacterData, error) {
	v, exists := CharacterDataCollection[c]
	if !exists {
		return nil, fmt.Errorf("not exists")
	}

	return &v, nil
}
