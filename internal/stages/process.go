package stages

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/calderwd/operator-helper/internal/config"
	"github.com/calderwd/operator-helper/internal/stagevalues"
	"k8s.io/apimachinery/pkg/types"
)

func (s Stage) Process(config config.ReconcileConfig, values stagevalues.Values, nn types.NamespacedName) error {

	if len(s.Resources) == 0 {
		return nil
	}

	for _, resource := range s.Resources {

		var fileName string = string(resource)

		dir, _ := os.Getwd()

		filePath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)

		tmplate, err := template.ParseFiles(filePath)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmplate.Execute(&buf, values); err != nil {
			return err
		}

		fmt.Println(buf.String())
	}

	return nil
}
