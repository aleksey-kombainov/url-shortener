package api

type ShortenerRequest struct {
	Url string `json:"url"`
}

type ShortenerResponse struct {
	Result string `json:"result"`
}
