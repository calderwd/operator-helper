package stages

type Resources string

type Stage struct {
	Name      string      `yaml:"name"`
	Resources []Resources `yaml:"resources"`
}
