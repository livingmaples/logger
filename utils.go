package logger

import (
	"fmt"
	"strings"
)

func mergeDefaultsToOutput(a map[string]string, b map[string]interface{}) map[string]string {
	builder := strings.Builder{}
	for k, v := range b {
		if _, ok := a[k]; ok {
			builder.WriteString("defaults.")
			builder.WriteString(k)
			a[builder.String()] = fmt.Sprint(v)
			continue
		}

		a[k] = fmt.Sprint(v)
	}

	return a
}