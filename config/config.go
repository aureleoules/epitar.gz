package config

type Module struct {
	Name    string            `yaml:"name"`
	Path    string            `yaml:"path"`
	Config  map[string]string `yaml:"config"`
	Enabled bool              `yaml:"enabled"`
}

type Config struct {
	Output  string   `yaml:"output"`
	Modules []Module `yaml:"modules"`
}
