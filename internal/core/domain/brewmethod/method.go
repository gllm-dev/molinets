package brewmethod

type Method struct {
	ID      string
	Name    string
	Aliases []string
	MinUM   float64
	MaxUM   float64
}

func (m *Method) Contains(um float64) bool { return um >= m.MinUM && um <= m.MaxUM }

func Matching(methods []Method, um float64) []string {
	out := []string{}
	for i := range methods {
		if methods[i].Contains(um) {
			out = append(out, methods[i].Name)
		}
	}
	return out
}
