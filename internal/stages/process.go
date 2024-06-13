package stages

import (
	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/stagevalues"
	"k8s.io/apimachinery/pkg/types"
)

func (s Stage) Process(config config.ReconcileConfig, values stagevalues.Values, nn types.NamespacedName) error {

	return nil
}
