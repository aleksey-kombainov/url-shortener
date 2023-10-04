package entities

type Shortcut struct {
	Id          uint64 `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
