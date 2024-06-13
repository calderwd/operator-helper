package stagevalues

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	config "github.com/calderwd/operator-helper/internal/config"
)

type Values map[string]interface{}

func (v Values) Load(config config.ReconcileConfig, nn types.NamespacedName, cr client.Object) error {

	return nil
}
