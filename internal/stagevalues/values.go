package stagevalues

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	config "github.com/calderwd/operator-helper/internal/config"
)

type Values map[interface{}]interface{}

func (v *Values) Load(config config.ReconcileConfig, nn types.NamespacedName, cr client.Object) error {

	if config.ValuesPath == "" {
		return nil
	}

	b, err := os.ReadFile(config.ValuesPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return v.extractValuesFromCR(cr)
}

func (v Values) extractValuesFromCR(cr client.Object) error {

	values := map[string]interface{}{}
	values["name"] = cr.GetName()
	values["namespace"] = cr.GetNamespace()

	v["cr"] = values
	return nil
}

func (vm Values) AsMapOfString() map[string]interface{} {

	result := make(map[string]interface{})

	var n func(m interface{}) interface{}
	n = func(m interface{}) interface{} {
		result := m

		if em, ok := m.(Values); ok {
			res := make(map[string]interface{})

			for k := range em {
				v := em[k]
				res[k.(string)] = n(v)
			}
			result = res
		}
		return result
	}

	for k := range vm {
		vv := vm[k]
		v := n(vv)
		result[k.(string)] = v
	}

	return result
}

func (v Values) GetString(path string) string {

	vw := v.AsMapOfString()

	spl := strings.Split(path, ".")

	var ptr = vw

	if vv, ok := ptr[spl[0]]; ok {

		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			ptr = vv.(map[string]interface{})
		case reflect.String:
			x := vv.(string)
			return x
		case reflect.Int:
			x := strconv.Itoa(vv.(int))
			return x
		}
	}

	var vv *string

	for _, s := range spl[1:] {
		if ptr, vv = v.walk(ptr, s); vv != nil {
			return *vv
		}
	}
	return ""
}

func (v Values) walk(mv map[string]interface{}, s string) (map[string]interface{}, *string) {
	if vv, ok := mv[s]; ok {
		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			return vv.(map[string]interface{}), nil
		case reflect.String:
			x := vv.(string)
			return mv, &x
		case reflect.Int:
			x := strconv.Itoa(vv.(int))
			return mv, &x
		}
	}
	return mv, nil
}
