package models

type ScannerInfo struct {
	ScannerName           string                       `yaml:"ScannerName"`
	Description           string                       `yaml:"Description"`
	Type                  string                       `yaml:"Type"`
	ToolDirectoryName     string                       `yaml:"ToolDirectoryName"`
	SupportedPlatform     map[string]SupportedPlatform `yaml:"SupportedPlatform"`
	ScannerTemplateAuthor string                       `yaml:"ScannerTemplateAuthor"`
	TimeOut               int                          `yaml:"TimeOut"`
}

type SupportedPlatform struct {
	StartCmd string `yaml:"StartCmd"`
}
