package stagevalues

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	config "github.com/calderwd/operator-helper/internal/config"
	"github.com/go-logr/logr"
)

var (
	ErrToMarshalCr   = errors.New("unable to unmarshal cr type")
	ErrToUnmarshalCr = errors.New("unable to marshal type")
)

const (
	CR string = "cr"
)

type Values map[interface{}]interface{}

func (v *Values) Load(config config.ReconcileConfig, nn types.NamespacedName, cr client.Object, log logr.Logger) error {

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

	return v.extractValuesFromCR(cr, log)
}

func (v Values) extractValuesFromCR(cr interface{}, log logr.Logger) error {

	if bb, err := yaml.Marshal(cr); err == nil {

		bstr := string(bb)
		log.Info(bstr)

		r := Values{}
		if err = yaml.Unmarshal(bb, r); err != nil {
			return err
		}

		v[CR] = v.sanitise(r)

		// Add shortcuts
		r["name"] = v.GetString("cr.objectmeta.name")
		r["namespace"] = v.GetString("cr.objectmeta.namespace")

	} else {
		return ErrToMarshalCr
	}
	return nil
}

func (v Values) sanitise(values Values) Values {

	values.deleteEntry("objectmeta.managedfields")
	values.deleteEntry("status")
	return values
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

func (v Values) deleteEntry(path string) Values {
	spl := strings.Split(path, ".")

	var ptr = v

	if vv, ok := ptr[spl[0]]; ok {
		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			ptr = vv.(Values)
			for _, s := range spl[1:] {
				ptr, vv = v.walk(ptr, s)
			}
			delete(ptr, spl[len(spl)-1])
		}
	}

	return v
}

func (v Values) GetString(path string) string {

	spl := strings.Split(path, ".")

	var ptr = v

	if vv, ok := ptr[spl[0]]; ok {

		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			ptr = vv.(Values)
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

func (v Values) walk(mv map[interface{}]interface{}, s string) (map[interface{}]interface{}, *string) {
	if vv, ok := mv[s]; ok {
		rv := reflect.Indirect(reflect.ValueOf(vv))

		switch rv.Kind() {
		case reflect.Map:
			return vv.(Values), nil
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
