package config

type ReconcileConfig struct {
}

func (r ReconcileConfig) Validate(config ReconcileConfig) (ReconcileConfig, error) {
	return r, nil
}
