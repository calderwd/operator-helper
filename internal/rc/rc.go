package rc

import (
	"context"
	"strings"

	"gopkg.in/yaml.v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type DClient struct {
	Client *dynamic.DynamicClient
}

func GetDynamicClient() (*DClient, error) {

	cfg := ctrl.GetConfigOrDie()

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &DClient{Client: dynamicClient}, nil
}

func (c DClient) CreateResourceFromYaml(rs string) error {

	em := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(rs), &em); err != nil {
		return err
	}

	em = normalizeResource(em)

	resource := unstructured.Unstructured{}
	resource.SetUnstructuredContent(em)

	gvk := resource.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind) + "s", // TODO - need proper fix using mapper
	}

	_, err := c.Client.Resource(gvr).Namespace("test").Create(context.TODO(), &resource, v1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func normalizeResource(em map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	var n func(m interface{}) interface{}
	n = func(m interface{}) interface{} {
		result := m

		if em, ok := m.(map[interface{}]interface{}); ok {
			res := make(map[string]interface{})

			for k := range em {
				v := em[k]
				res[k.(string)] = n(v)
			}
			result = res
		}

		if em, ok := m.([]interface{}); ok {
			res := make([]interface{}, len(m.([]interface{})))
			for ai, ae := range em {
				res[ai] = n(ae)
			}
			result = res
		}

		return result
	}

	for k := range em {
		vv := em[k]
		v := n(vv)
		result[k] = v
	}

	return result
}
