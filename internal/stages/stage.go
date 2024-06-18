package stages

type Resource string

type Stage struct {
	Name      string     `yaml:"name"`
	Resources []Resource `yaml:"resources"`
}
