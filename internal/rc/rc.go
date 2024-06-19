package rc

import (
	"context"

	"gopkg.in/yaml.v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	ctrl "sigs.k8s.io/controller-runtime"
)

type DClient struct {
	Client *dynamic.DynamicClient
	Mapper meta.RESTMapper
}

func GetDynamicClient() (*DClient, error) {

	cfg := ctrl.GetConfigOrDie()

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(cfg)

	apis, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(apis)

	return &DClient{
		Client: dynamicClient,
		Mapper: mapper,
	}, nil
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
	mapping, err := c.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var res dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		res = c.Client.Resource(mapping.Resource).Namespace("test")
	} else {
		res = c.Client.Resource(mapping.Resource)
	}

	_, err = res.Create(context.TODO(), &resource, v1.CreateOptions{})

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
