package stagevalues

import (
	"testing"

	config "github.com/calderwd/operator-helper/internal/config"
	log "github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

func TestLoad(t *testing.T) {

	config := config.ReconcileConfig{
		ValuesPath: "../test/values.yaml",
	}

	nn := types.NamespacedName{}

	dummy := Dummy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "testnamespace",
		},
	}

	values := Values{}
	values.Load(config, nn, &dummy, log.Logger{})

	assert.Equal(t, values.GetString("title"), "Test")
	assert.Equal(t, values.GetString("app.name"), "my-test")
	assert.Equal(t, values.GetString("app.cpu.request"), "50m")
	assert.Equal(t, values.GetString("middle.replicas"), "2")

	assert.Equal(t, values.GetString("cr.name"), "test")
	assert.Equal(t, values.GetString("cr.namespace"), "testnamespace")
}
