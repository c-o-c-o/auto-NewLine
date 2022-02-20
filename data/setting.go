package data

type Setting struct {
	BreakStr      string         `yaml:"BreakString"`
	BreakPatterns []BreakPattern `yaml:"BreakPattern"`
	PriorityWait  float32        `yaml:"PriorityWait"`
	DeleteFixs    []string       `yaml:"DeleteFix"`
}

type BreakPattern struct {
	Priority int                 `yaml:"Priority"`
	Patterns []map[string]string `yaml:"Pattern"`
}
