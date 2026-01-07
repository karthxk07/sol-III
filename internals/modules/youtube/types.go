package youtube


type Snippet struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CategoryID  int      `json:"categoryId"`
}
