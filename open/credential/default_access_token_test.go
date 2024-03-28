package credential

import (
	"context"
	"fmt"
	"testing"

	"github.com/RichXan/douyin-sdk/cache"
	"github.com/alicebob/miniredis/v2"
)

func TestGetAccessToken(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Error("miniredis.Run Error", err)
	}
	t.Cleanup(server.Close)
	var (
		// timeoutDuration = time.Second
		ctx  = context.Background()
		opts = &cache.RedisOptions{
			Addr: server.Addr(),
		}
		redis = cache.NewRedis(ctx, opts)
	)
	ak := NewDefaultAccessToken("aw05az2qjv******", "7802f4e6f243e659d51135445fe********", "douyin", redis)
	token, err := ak.GetAccessToken()
	if err != nil {
		t.Error("miniredis.Run Error", err)
	}

	fmt.Println(token)
}
