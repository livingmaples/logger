package logger

import (
	"strings"
)

func mergeDefaultsToOutput(a *map[string]interface{}, b map[string]interface{}) *map[string]interface{} {
	builder := strings.Builder{}
	for k, v := range b {
		if _, ok := (*a)[k]; ok {
			builder.WriteString("defaults.")
			builder.WriteString(k)
			(*a)[builder.String()] = v
			continue
		}

		(*a)[k] = v
	}

	return a
}
