package credential

import (
	"context"
	"testing"

	"github.com/RichXan/douyin-sdk/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

// TestGetTicketFromServer .
func TestGetTicketFromServer(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Error("miniredis.Run Error", err)
	}
	t.Cleanup(server.Close)
	ctx := context.Background()
	opts := &cache.RedisOptions{
		Addr: server.Addr(),
	}
	redis := cache.NewRedis(ctx, opts)

	defer gock.Off()
	gock.New(getClientTokenURL).Reply(200).JSON(map[string]interface{}{
		"data": map[string]interface{}{
			"access_token": "mock-access-token",
			"description":  "mock-description",
			"error_code":   0,
			"expires_in":   7200,
		},
		"message": "mock-message",
		"extra": map[string]interface{}{
			"logid": "mock-logid",
			"now":   10,
		},
	})
	ct := NewDefaultClientToken("arg-ck", "arg-cs", "arg-cache-key-prefix", redis)

	clientToken, err := ct.GetClientToken()
	assert.Nil(t, err)
	assert.Equal(t, "mock-access-token", clientToken, "they should be equal")
}
