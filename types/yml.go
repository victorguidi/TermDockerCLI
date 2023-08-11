package types

type Host struct {
	IP       string `yaml:"ip"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Hosts []Host `yaml:"hosts"`
}
