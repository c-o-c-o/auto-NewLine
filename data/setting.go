package data

type Setting struct {
	BreakStr      string         `yaml:"BreakString"`
	BreakPatterns []BreakPattern `yaml:"BreakPattern"`
	PriorityWait  float64        `yaml:"PriorityWait"`
	DeleteFixs    []string       `yaml:"DeleteFix"`
}

type BreakPattern struct {
	Priority float64             `yaml:"Priority"`
	Patterns []map[string]string `yaml:"Pattern"`
}
