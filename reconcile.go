package reconcile

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	RequeueAfter = reconcile.Result{Requeue: true, RequeueAfter: time.Minute}
)

type ReconcileConfig struct {
}

func Reconcile(config ReconcileConfig, c client.Client, nn types.NamespacedName, cr client.Object, log logr.Logger) (ctrl.Result, error) {

	err := c.Get(context.TODO(), nn, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to load CR")
		return RequeueAfter, err
	}

	return ctrl.Result{}, nil
}
