package config

type Config struct {
	ApiToken  string          `json:"apiToken"`
	Libraries []LibraryConfig `json:"libraries"`
}

type LibraryConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
	ID   string `json:"id"`
}
