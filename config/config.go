package config

var Cfg Config

type ModuleMeta struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Authors     []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"authors"`
}

type Module struct {
	Path         string            `yaml:"path"`
	Options      map[string]string `yaml:"options"`
	Enable       bool              `yaml:"enable"`
	IndexOptions struct {
		Enable bool     `yaml:"enable"`
		Files  []string `yaml:"files"`
	} `yaml:"index"`
	Meta ModuleMeta `yaml:"meta"`
}

type IndexConfig struct {
	Store string `yaml:"store"`
	Sonic struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"sonic"`
}

type Config struct {
	Output  string      `yaml:"output"`
	Index   IndexConfig `yaml:"index"`
	Modules []Module    `yaml:"modules"`
}
