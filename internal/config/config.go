package config

import (
	"fmt"
)

const (
	ERROR_STAGES_PATH_EMPTY string = "stages path must be set"
)

type ReconcileConfig struct {
	StagesPath string
	ValuesPath string
}

func (r ReconcileConfig) Validate(config ReconcileConfig) (ReconcileConfig, error) {

	if r.StagesPath == "" {
		return r, fmt.Errorf(ERROR_STAGES_PATH_EMPTY)
	}

	// ValuesPath is optional

	return r, nil
}
