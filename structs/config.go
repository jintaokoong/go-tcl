package structs

type Config struct {
	Channels []string `yaml:"channels"`
	Roles    []string `yaml:"roles"`
	Database struct {
		Name             string `yaml:"name"`
		Collection       string `yaml:"collection"`
		ConnectionString string `yaml:"connectionString"`
	} `yaml:"database"`
}
