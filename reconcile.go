package reconcile

import (
	"context"
	"time"

	rconfig "github.com/calderwd/operator-helper/internal/config"
	rstages "github.com/calderwd/operator-helper/internal/stages"
	stagevalues "github.com/calderwd/operator-helper/internal/stagevalues"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	RequeueAfter      = reconcile.Result{Requeue: true, RequeueAfter: time.Minute}
	RequeueAfterError = reconcile.Result{Requeue: true, RequeueAfter: time.Minute}
)

const (
	ERROR_FAILED_TO_LOAD_CR        string = "failed to load cr"
	ERROR_CONFIG_IS_INVALID        string = "invalid config specified"
	ERROR_FAILED_TO_LOAD_VALUES    string = "failed to load values"
	ERROR_FAILED_TO_LOAD_STAGES    string = "failed to load stages"
	ERROR_FAILED_TO_PROCESS_STAGE  string = "failed to process stage"
	ERROR_FAILED_TO_VALIDATE_STAGE string = "failed to validate stage"
)

type ReconcileConfig rconfig.ReconcileConfig

func (rc ReconcileConfig) toInternalConfig() (rconfig.ReconcileConfig, error) {
	m := rconfig.ReconcileConfig{
		StagesPath: rc.StagesPath,
		ValuesPath: rc.ValuesPath,
	}
	return m.Validate(m)
}

func Reconcile(config ReconcileConfig, c client.Client, nn types.NamespacedName, cr client.Object, log logr.Logger) (ctrl.Result, error) {

	// Load the CR with the current instance
	if r, err := loadCR(c, nn, cr); err != nil {
		log.Error(err, ERROR_FAILED_TO_LOAD_CR)
		return r, err
	}

	var err error
	var ic rconfig.ReconcileConfig

	if ic, err = config.toInternalConfig(); err != nil {
		log.Error(err, ERROR_CONFIG_IS_INVALID)
		return RequeueAfterError, err
	}

	var values stagevalues.Values

	// Extract values from CR
	if err := values.Load(ic, nn, cr); err != nil {
		log.Error(err, ERROR_FAILED_TO_LOAD_VALUES)
		return RequeueAfterError, err
	}

	var stages rstages.Stages

	// Get associated stages
	// Apply values in case where stages are templated
	if err := stages.Load(ic, values); err != nil {
		log.Error(err, ERROR_FAILED_TO_LOAD_STAGES)
		return RequeueAfterError, err
	}

	// For each stage
	//   Process stage
	//   if validation present then wait until pass
	for _, stage := range stages.Stages {
		if err := stage.Process(ic, values, nn); err != nil {
			log.Error(err, ERROR_FAILED_TO_PROCESS_STAGE)
			return RequeueAfterError, err
		}

		if err := stage.Validate(ic, values, nn); err != nil {
			log.Error(err, ERROR_FAILED_TO_VALIDATE_STAGE)
			return RequeueAfterError, err
		}
	}

	return ctrl.Result{}, nil
}

func loadCR(c client.Client, nn types.NamespacedName, cr client.Object) (ctrl.Result, error) {
	err := c.Get(context.TODO(), nn, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return RequeueAfter, err
	}
	return ctrl.Result{}, nil
}
