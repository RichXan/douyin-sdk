package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/RichXan/douyin-sdk/cache"
	"github.com/RichXan/douyin-sdk/response"
	"github.com/RichXan/douyin-sdk/util"
	// "honnef.co/go/tools/lintcmd/cache"
	// "github.com/silenceper/wechat/v2/cache"
)

const (
	// accessTokenURL 获取access_token的接口
	accessTokenURL = "https://open.douyin.com/oauth/client_token/"
	// stableAccessTokenURL 获取沙盒的接口
	sandBoxAccessTokenURL = "https://open-sandbox.douyin.com/oauth/client_token/"
	CacheKeyPrefix        = "douyin_"
)

var Headers = map[string]string{
	// "Content-Type": "application/json",
}

// DefaultAccessToken 默认AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewDefaultAccessToken new DefaultAccessToken
func NewDefaultAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &DefaultAccessToken{
		appID:           appID,
		appSecret:       appSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAccessToken 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从cache中取
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.appID)

	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		if accessToken = val.(string); accessToken != "" {
			return
		}
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从douyin服务器上获取到不同token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()

	// 双检，防止重复从douyin服务器获取
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		if accessToken = val.(string); accessToken != "" {
			// fmt.Println(err)
			return
		}
	}

	// cache失效，从douyin服务器获取
	var resAccessToken ResAccessToken
	if resAccessToken, err = ak.GetAccessTokenDirectly(ctx); err != nil {
		fmt.Println(err)
		return
	}
	if resAccessToken.AccessToken == "" {
		err = fmt.Errorf("token is nil")
		return
	}
	if resAccessToken.ExpiresIn == 0 {
		err = fmt.Errorf("ExpiresIn is nil")
		return
	}
	expires := resAccessToken.ExpiresIn - 1500

	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

	accessToken = resAccessToken.AccessToken
	return
}

func (ak *DefaultAccessToken) GetAccessTokenDirectly(ctx context.Context) (resAccessToken ResAccessToken, err error) {
	b, err := util.PostJSONContext(ctx, accessTokenURL, map[string]interface{}{
		"grant_type":    "client_credential",
		"client_key":    ak.appID,
		"client_secret": ak.appSecret,
		// "force_refresh": forceRefresh,
	}, Headers)
	if err != nil {
		return
	}
	// if b != nil {
	// 	err = fmt.Errorf(string(b))
	// 	return
	// }
	res := response.Response{}
	res.Data = &resAccessToken
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}

	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.Message)
		return
	}
	return
}

// StableAccessToken 获取稳定版接口调用凭据(与getAccessToken获取的调用凭证完全隔离，互不影响)
// 不强制更新access_token,可用于不同环境不同服务而不需要分布式锁以及公用缓存，避免access_token争抢
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-access-token/getStableAccessToken.html
type SandBoxAccessToken struct {
	appID          string
	appSecret      string
	cacheKeyPrefix string
	cache          cache.Cache
}

// NewStableAccessToken new StableAccessToken
func NewSandBoxAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache is need")
	}
	return &SandBoxAccessToken{
		appID:          appID,
		appSecret:      appSecret,
		cache:          cache,
		cacheKeyPrefix: cacheKeyPrefix,
	}
}

// GetAccessToken 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *SandBoxAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *SandBoxAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从cache中取
	accessTokenCacheKey := fmt.Sprintf("%s_sand_box_access_token_%s", ak.cacheKeyPrefix, ak.appID)
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// cache失效，从douyin服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ak.GetAccessTokenDirectly(ctx, false)
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 300
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

	accessToken = resAccessToken.AccessToken
	return
}

// GetAccessTokenDirectly 从douyin获取access_token
func (ak *SandBoxAccessToken) GetAccessTokenDirectly(ctx context.Context, forceRefresh bool) (resAccessToken ResAccessToken, err error) {
	b, err := util.PostJSONContext(ctx, sandBoxAccessTokenURL, map[string]interface{}{
		"grant_type":    "client_credential",
		"client_key":    ak.appID,
		"client_secret": ak.appSecret,
		// "force_refresh": forceRefresh,
	}, Headers)
	if err != nil {
		return
	}
	var res response.Response
	res.Data = &resAccessToken
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	resAccessToken = res.Data.(ResAccessToken)
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get stable access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.Message)
		return
	}
	return
}

// // WorkAccessToken 企业douyinAccessToken 获取
// type WorkAccessToken struct {
// 	CorpID          string
// 	CorpSecret      string
// 	cacheKeyPrefix  string
// 	cache           cache.Cache
// 	accessTokenLock *sync.Mutex
// }

// // NewWorkAccessToken new WorkAccessToken
// func NewWorkAccessToken(corpID, corpSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
// 	if cache == nil {
// 		panic("cache the not exist")
// 	}
// 	return &WorkAccessToken{
// 		CorpID:          corpID,
// 		CorpSecret:      corpSecret,
// 		cache:           cache,
// 		cacheKeyPrefix:  cacheKeyPrefix,
// 		accessTokenLock: new(sync.Mutex),
// 	}
// }

// // GetAccessToken 企业douyin获取access_token,先从cache中获取，没有则从服务端获取
// func (ak *WorkAccessToken) GetAccessToken() (accessToken string, err error) {
// 	return ak.GetAccessTokenContext(context.Background())
// }

// // GetAccessTokenContext 企业douyin获取access_token,先从cache中获取，没有则从服务端获取
// func (ak *WorkAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
// 	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从douyin服务器上获取到不同token
// 	ak.accessTokenLock.Lock()
// 	defer ak.accessTokenLock.Unlock()
// 	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.CorpID)
// 	val := ak.cache.Get(accessTokenCacheKey)
// 	if val != nil {
// 		accessToken = val.(string)
// 		return
// 	}

// 	// cache失效，从douyin服务器获取
// 	var resAccessToken ResAccessToken
// 	resAccessToken, err = GetTokenFromServerContext(ctx, fmt.Sprintf(workAccessTokenURL, ak.CorpID, ak.CorpSecret))
// 	if err != nil {
// 		return
// 	}

// 	expires := resAccessToken.ExpiresIn - 1500
// 	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

// 	accessToken = resAccessToken.AccessToken
// 	return
// }

// GetTokenFromServer 强制从douyin服务器获取token
func GetTokenFromServer(url string) (resAccessToken ResAccessToken, err error) {
	return GetTokenFromServerContext(context.Background(), url)
}

// GetTokenFromServerContext 强制从douyin服务器获取token
func GetTokenFromServerContext(ctx context.Context, url string) (resAccessToken ResAccessToken, err error) {
	var body []byte
	body, err = util.HTTPGetContext(ctx, url, Headers)
	if err != nil {
		return
	}
	var res response.Response
	res.Data = resAccessToken
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	resAccessToken = res.Data.(ResAccessToken)
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.Message)
		return
	}
	return
}
