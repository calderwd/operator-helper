package stages

import (
	"testing"

	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/stagevalues"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {

	config := config.ReconcileConfig{
		StagesPath: "../test/stages.yaml",
	}

	values := stagevalues.Values{}

	stages := Stages{}
	stages.Load(config, values)

	assert.Equal(t, stages.Title, "stages-test")
	assert.Equal(t, stages.Version, "1.0")
	assert.Equal(t, len(stages.Stages), 2)
	assert.Equal(t, stages.Stages[0].Name, "middleware-install")
	assert.Equal(t, len(stages.Stages[0].Resources), 3)
	assert.Equal(t, stages.Stages[1].Name, "app-install")
	assert.Equal(t, len(stages.Stages[1].Resources), 3)
}
