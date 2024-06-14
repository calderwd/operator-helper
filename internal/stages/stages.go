package stages

import (
	"os"

	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/stagevalues"
	"gopkg.in/yaml.v2"
)

type Stages struct {
	Title   string  `yaml:"title"`
	Version string  `yaml:"version"`
	Stages  []Stage `yaml:"stage"`
}

func (s *Stages) Load(config config.ReconcileConfig, values stagevalues.Values) error {

	if b, err := os.ReadFile(config.StagesPath); err != nil {
		return err
	} else {
		return yaml.Unmarshal(b, s)
	}
}
