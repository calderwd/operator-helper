package stages

import (
	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/stagevalues"
)

type Stages []Stage

func (s Stages) Load(config config.ReconcileConfig, values stagevalues.Values) error {

	return nil
}
