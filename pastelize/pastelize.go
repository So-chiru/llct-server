package pastelize

import (
	"image"

	"github.com/dayvonjersen/vibrant"
	"github.com/teacat/noire"
)

type GoColorStruct struct {
	Main     string `json:"main"`
	Sub      string `json:"sub"`
	Text     string `json:"text"`
	MainDark string `json:"mainDark"`
	SubDark  string `json:"subDark"`
	TextDark string `json:"textDark"`
}

type GoColorStructV2 struct {
	Primary        string `json:"primary"`
	PrimaryLight   string `json:"primaryLight"`
	PrimaryDark    string `json:"primaryDark"`
	Secondary      string `json:"secondary"`
	SecondaryLight string `json:"secondaryLight"`
	SecondaryDark  string `json:"secondaryDark"`
}

type ColorData struct {
	Size  int
	Color noire.Color
}

func GeneratePalette(img image.Image) (*GoColorStructV2, error) {
	// var d clusters.Observations

	// for x := 0; x <= image.Bounds().Max.X; x++ {
	// 	for y := 0; y <= image.Bounds().Max.Y; y++ {
	// 		r, g, b, _ := image.At(x, y).RGBA()

	// 		d = append(d, clusters.Coordinates{float64(r) / 255, float64(g) / 255, float64(b) / 255})
	// 	}
	// }

	// km, err := kmeans.NewWithOptions(0.05, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// clusters, _ := km.Partition(d, 6)

	// for i := 0; i < 2; i++ {
	// 	clusters.Recenter()
	// }

	// var colors []ColorData
	// for _, c := range clusters {
	// 	if len(c.Observations) < 10000 {
	// 		continue
	// 	}

	// 	var rgb = noire.NewRGB(c.Center[0], c.Center[1], c.Center[2])

	// 	colors = append(colors, ColorData{
	// 		Size:  len(c.Observations),
	// 		Color: rgb,
	// 	})

	// 	fmt.Println(rgb, rgb.Luminanace())
	// }

	// sort.Slice(colors, func(i int, j int) bool {
	// 	return colors[i].Size > colors[j].Size && colors[i].Color.Luminanace() < colors[j].Color.Luminanace()
	// })

	// var sortLuminance []ColorData
	// sortLuminance = append(sortLuminance, colors...)

	// sort.Slice(sortLuminance, func(i int, j int) bool {
	// 	return sortLuminance[i].Color.Luminanace() > sortLuminance[j].Color.Luminanace()
	// })

	// // fmt.Println(colors)

	// var str = GoColorStructV2{
	// 	Primary:   "#" + colors[0].Color.Hex(),
	// 	Secondary: "#" + sortLuminance[len(sortLuminance)-1].Color.Hex(),
	// }

	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		return nil, err
	}

	var pal = palette.ExtractAwesome()

	var str = GoColorStructV2{
		Primary:        pal["Vibrant"].Color.String(),
		PrimaryLight:   pal["LightVibrant"].Color.String(),
		PrimaryDark:    pal["DarkVibrant"].Color.String(),
		Secondary:      pal["Muted"].Color.String(),
		SecondaryLight: pal["LightMuted"].Color.String(),
		SecondaryDark:  pal["DarkMuted"].Color.String(),
	}

	return &str, nil
}
