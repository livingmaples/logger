package formatters

import (
	"encoding/json"
)

type JsonFormatter struct {
	Timestamp bool
}

func (f JsonFormatter) Format(level string, data map[string]string) ([]byte, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return j, nil
}
