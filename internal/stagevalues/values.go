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

	if b, err := os.ReadFile(config.ValuesPath); err != nil {
		return err
	} else {
		return yaml.Unmarshal(b, v)
	}
}

func (v Values) GetString(path string) string {

	spl := strings.Split(path, ".")

	var ptr map[interface{}]interface{}

	if vv, ok := v[spl[0]]; ok {
		if ptr, ok = vv.(map[interface{}]interface{}); !ok {
			if vvv, ok := vv.(string); ok {
				return vvv
			}
			if vvv, ok := vv.(int); ok {
				return strconv.Itoa(vvv)
			}
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

func (v Values) walk(mv map[interface{}]interface{}, s string) (map[interface{}]interface{}, *string) {
	if vv, ok := mv[s]; ok {
		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			return vv.(map[interface{}]interface{}), nil
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
