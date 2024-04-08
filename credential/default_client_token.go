package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/RichXan/douyin-sdk/cache"
	"github.com/RichXan/douyin-sdk/util"
)

// getClientTokenURL 获取token的url
const getClientTokenURL = "https://open.douyin.com/oauth/client_token"
const CacheKeyPrefix = "ClientTokenPrefix"

// DefaultClientToken 默认获取client token方法
type DefaultClientToken struct {
	clientKey      string
	clientSecret   string
	cacheKeyPrefix string
	cache          cache.Cache
	// clientTokenLock 读写锁 同一个clientkey一个
	clientTokenLock *sync.Mutex
}

// ResClientToken struct
type ResClientToken struct {
	Data struct {
		AccessToken string `json:"access_token"`
		Description string `json:"description"`
		ErrorCode   int64  `json:"error_code"`
		ExpiresIn   int64  `json:"expires_in"`
	} `json:"data"`
	Message string `json:"message"`
	Extra   struct {
		LogID string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
}

// NewDefaultClientToken new
func NewDefaultClientToken(clientKey, clientSecret, cacheKeyPrefix string, cache cache.Cache) ClientTokenHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &DefaultClientToken{
		clientKey:       clientKey,
		clientSecret:    clientSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		clientTokenLock: new(sync.Mutex),
	}
}

// GetClientToken 获取client_token,先从cache中获取，没有则从服务端获取
func (ct *DefaultClientToken) GetClientToken() (clientToken string, err error) {
	return ct.GetClientTokenContext(context.Background())
}

// GetClientTokenContext 获取client_token,先从cache中获取，没有则从服务端获取
func (ct *DefaultClientToken) GetClientTokenContext(ctx context.Context) (clientToken string, err error) {
	// 先从cache中取
	clientTokenCacheKey := fmt.Sprintf("%s_client_token_%s", ct.cacheKeyPrefix, ct.clientKey)
	if val := ct.cache.Get(clientTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从抖音服务器上获取到不同token
	ct.clientTokenLock.Lock()
	defer ct.clientTokenLock.Unlock()

	// 双检，防止重复从抖音服务器获取
	if val := ct.cache.Get(clientTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// cache失效，从抖音服务器获取
	var resClientToken *ResClientToken
	if resClientToken, err = ct.GetClientTokenFromServer(); err != nil {
		return
	}

	expires := resClientToken.Data.ExpiresIn - 300
	if err = ct.cache.Set(clientTokenCacheKey, resClientToken.Data.AccessToken, time.Duration(expires)*time.Second); err != nil {
		return
	}

	clientToken = resClientToken.Data.AccessToken
	return
}

// GetClientTokenFromServer 强制从抖音服务器获取token
func (ct *DefaultClientToken) GetClientTokenFromServer() (resClientToken *ResClientToken, err error) {
	return ct.GetClientTokenFromServerContext(context.Background())
}

// GetClientTokenFromServerContext 强制从抖音服务器获取token
func (ct *DefaultClientToken) GetClientTokenFromServerContext(ctx context.Context) (resClientToken *ResClientToken, err error) {
	body := map[string]interface{}{
		"client_key":    ct.clientKey,
		"client_secret": ct.clientSecret,
		"grant_type":    "client_credential",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	header := map[string]string{"Content-Type": "application/json"}

	response, err := util.HTTPPostContext(ctx, getClientTokenURL, bodyBytes, header)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &resClientToken)
	if err != nil {
		return
	}
	if resClientToken.Data.ErrorCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resClientToken.Data.ErrorCode, resClientToken.Message)
		return
	}
	return
}
