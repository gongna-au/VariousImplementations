package ChannelApplicationScenario

import "context"

func RequestUniqueId() {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	go Worker(parentCtx, "parent", "parentworker1")
	go Worker(parentCtx, "parent", "parentworker2")
	go Worker(parentCtx, "parent", "parentworker3")

	sonCtxl, sonCancel := context.WithCancel(parentCtx)
	go Worker(sonCtxl, "son", "sonworker1")
	go Worker(sonCtxl, "son", "sonworker2")
	go Worker(sonCtxl, "son", "sonworker3")
	sonCancel()
	parentCancel()

}

func Worker(ctx context.Context, key string, value string) {
	ctx = context.WithValue(ctx, key, value)
}
