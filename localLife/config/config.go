package config

import "github.com/RichXan/douyin-sdk/cache"

// config for partner
type Config struct {
	ClientKey    string `json:"client_key"` // 授权应用的key
	ClinetSecret string `json:"client_secret"`	// 授权应用的secret
	AccountID   string // 商家ID，传入时服务商须与该商家满足授权关系
	Cache        cache.Cache
}
