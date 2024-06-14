package reconcile

import (
	"context"
	"testing"

	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/test"
	"github.com/go-logr/zapr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Dummy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

func (in *Dummy) DeepCopyInto(out *Dummy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}

func (in *Dummy) DeepCopy() *Dummy {
	if in == nil {
		return nil
	}
	out := new(Dummy)
	in.DeepCopyInto(out)
	return out
}

func (in *Dummy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func TestReconcileOnEmptyConfig(t *testing.T) {

	nn := types.NamespacedName{
		Namespace: "test",
		Name:      "test",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := test.NewMockClient(mockCtrl)

	l, _ := zap.NewDevelopment()
	ll := zapr.NewLogger(l)

	c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(

		func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			switch c := obj.(type) {
			case *Dummy:
				c.SetName("test")
				c.SetNamespace("test")
			}
			return nil
		},
	).MaxTimes(1)

	dummy := Dummy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
	}

	r, err := Reconcile(ReconcileConfig{}, c, nn, &dummy, ll)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), config.ERROR_STAGES_PATH_EMPTY)
	assert.Equal(t, r, RequeueAfterError)
}

func TestReconcile(t *testing.T) {

	nn := types.NamespacedName{
		Namespace: "test",
		Name:      "test",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := test.NewMockClient(mockCtrl)

	l, _ := zap.NewDevelopment()
	ll := zapr.NewLogger(l)

	c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(

		func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			switch c := obj.(type) {
			case *Dummy:
				c.SetName("test")
				c.SetNamespace("test")
			}
			return nil
		},
	).MaxTimes(1)

	dummy := Dummy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
	}

	rc := ReconcileConfig{
		StagesPath: "internal/test/stages.yaml",
	}
	r, err := Reconcile(rc, c, nn, &dummy, ll)

	assert.Nil(t, err)
	assert.Equal(t, r, ctrl.Result{})
}
