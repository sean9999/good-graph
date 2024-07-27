package society

type Message struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  []byte `json:"body"`
}
