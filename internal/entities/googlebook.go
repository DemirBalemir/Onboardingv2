package entities

type GoogleBook struct {
	ID         string `json:"id"`
	VolumeInfo struct {
		Title       string   `json:"title"`
		Authors     []string `json:"authors"`
		Description string   `json:"description"`
	} `json:"volumeInfo"`
}
