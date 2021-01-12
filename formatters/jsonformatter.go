package formatters

import (
	"context"
	"encoding/json"
)

type JsonFormatter struct {
	Timestamp bool
}

func (f JsonFormatter) Format(ctx context.Context, data *map[string]interface{}) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		j, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		return j, nil
	}
}
