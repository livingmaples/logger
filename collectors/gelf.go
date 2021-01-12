package collectors

import "context"

type GelfCollector struct{}

func (c GelfCollector) Ping(ctx context.Context) bool {
	return true
}

func (c GelfCollector) Write(ctx context.Context) {

}
