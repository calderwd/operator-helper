package rc

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestDynamicClient(t *testing.T) {

	dc, err := GetDynamicClient()

	assert.Nil(t, err)
	assert.NotNil(t, dc)

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}

	config := unstructured.Unstructured{}
	content := map[string]interface{}{
		"data": map[string]interface{}{
			"key1": "value1",
		},
	}
	config.SetUnstructuredContent(content)

	config.SetName("test")
	config.SetNamespace("test")
	config.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"})

	u, err := dc.Client.Resource(gvr).Namespace("test").Create(context.TODO(), &config, v1.CreateOptions{})

	assert.Nil(t, err)
	assert.NotNil(t, u)
}

func TestCreateResourceFromYaml(t *testing.T) {
	dc, err := GetDynamicClient()

	assert.Nil(t, err)
	assert.NotNil(t, dc)

	var fileName string = "../test/resources/dummy-configmap.yaml"

	dir, _ := os.Getwd()

	filePath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)

	b, err := os.ReadFile(filePath)
	assert.Nil(t, err)

	err = dc.CreateResourceFromYaml(string(b))
	assert.Nil(t, err)
}
