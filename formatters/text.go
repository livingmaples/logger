package formatters

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TextFormatter struct {
	Timestamp bool
}

// Format given data as Text format. This method may have bad performance and need to refactor
func (f TextFormatter) Format(ctx context.Context, level string, data map[string]string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		b := strings.Builder{}
		if f.Timestamp {
			now := time.Now()
			year, month, day := now.Date()
			b.Write([]byte(strconv.Itoa(year)))
			b.Write([]byte("/"))
			b.Write([]byte(strconv.Itoa(int(month))))
			b.Write([]byte("/"))
			b.Write([]byte(strconv.Itoa(day)))
			b.Write([]byte(" "))

			hour, min, sec := now.Clock()
			b.Write([]byte(strconv.Itoa(hour)))
			b.Write([]byte(":"))
			b.Write([]byte(strconv.Itoa(min)))
			b.Write([]byte(":"))
			b.Write([]byte(strconv.Itoa(sec)))
			b.Write([]byte("-> "))
		}

		for _, k := range orderedMapKeys(data) {
			b.Write([]byte(k))
			b.Write([]byte("=\""))
			b.Write([]byte(data[k]))
			b.Write([]byte("\" "))
		}

		return []byte(b.String()), nil
	}
}

func orderedMapKeys(v map[string]string) []string {
	keys := make([]string, 0)
	for k := range v {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
