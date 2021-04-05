package route

type JsonResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}
