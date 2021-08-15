package structs

type Config struct {
	Channels []string `yaml:"channels"`
	Roles    []string `yaml:"roles"`
}
