package route

type JsonResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

type GoColorStruct struct {
	Main     string `json:"main"`
	Sub      string `json:"sub"`
	Text     string `json:"text"`
	MainDark string `json:"mainDark"`
	SubDark  string `json:"subDark"`
	TextDark string `json:"textDark"`
}
