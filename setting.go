package main

type Setting struct {
	Newline       string              `yaml:"Newline-Character"`
	NotNLFeatures []map[string]string `yaml:"Not-NewLine-Features-Pattern"`
	NotNLStrings  []map[string]string `yaml:"Not-NewLine-String-Pattern"`
}
