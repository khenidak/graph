package graph

// General purpose string to vertix conversion
type stringVertix string

func (s stringVertix) Id() VertixId {
	return string(s)
}

var (
	graph_small_data = map[string][]string{
		"zero":  []string{"one", "two"},
		"one":   []string{"zero", "two"},
		"two":   []string{"two", "one", "three"},
		"three": []string{"one", "four", "five"},
		"four":  []string{"five", "three"},
		"five":  []string{"zero", "one"},
	}

	ordered = []string{"zero", "one", "two", "three", "four", "five"}
	mapped  = map[string]uint{
		"zero":  0,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
)
