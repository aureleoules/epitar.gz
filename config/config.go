package config

var Cfg Config

type IndexOptions struct {
	Enable bool     `yaml:"enable"`
	Files  []string `yaml:"files"`
}

type Module struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Options      map[string]string `yaml:"options"`
	Enable       bool              `yaml:"enable"`
	IndexOptions IndexOptions      `yaml:"index"`
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
