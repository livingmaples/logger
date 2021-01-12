package formatters

import (
	"context"
	"encoding/json"
)

type JsonFormatter struct {
	Timestamp bool
}

func (f JsonFormatter) Format(ctx context.Context, level string, data map[string]string) ([]byte, error) {
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
