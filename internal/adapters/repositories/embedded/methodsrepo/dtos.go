package methodsrepo

type methodDTO struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
	MinUM   float64  `json:"min_um"`
	MaxUM   float64  `json:"max_um"`
}
