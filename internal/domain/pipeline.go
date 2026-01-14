package domain

type Pipeline struct {
	Image string `yaml:"image"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name     string            `yaml:"name"`
	Commands []string          `yaml:"commands"`
	Env      map[string]string `yaml:"env"`
}
