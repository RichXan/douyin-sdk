package localLife

import (
	"github.com/RichXan/douyin-sdk/credential"
	"github.com/RichXan/douyin-sdk/localLife/certificate"
	"github.com/RichXan/douyin-sdk/localLife/config"
	"github.com/RichXan/douyin-sdk/localLife/context"
	"github.com/RichXan/douyin-sdk/localLife/goods"
	"github.com/RichXan/douyin-sdk/localLife/shop"
)

type LocalLife struct {
	ctx         *context.Context
	goods       *goods.Good
	certificate *certificate.Certificate
}

func NewLocalLife(cfg *config.Config) *LocalLife {
	defaultAkHandle := credential.NewDefaultClientToken(cfg.ClientKey, cfg.ClinetSecret, credential.CacheKeyPrefix, cfg.Cache)
	ctx := &context.Context{
		Config:            cfg,
		ClientTokenHandle: defaultAkHandle,
	}
	return &LocalLife{ctx: ctx}
}

// SetClientTokenHandle 自定义 access_token 获取方式
func (localLife *LocalLife) SetClientTokenHandle(accessTokenHandle credential.ClientTokenHandle) {
	localLife.ctx.ClientTokenHandle = accessTokenHandle
}

// GetContext get Context
func (localLife *LocalLife) GetContext() *context.Context {
	return localLife.ctx
}

// 门店管理接口
func (localLife *LocalLife) GetShop() *shop.Shop {
	return shop.NewShop(localLife.ctx)
}

// 团购验券接口
func (localLife *LocalLife) GetCertificate() *certificate.Certificate {
	if localLife.certificate == nil {
		localLife.certificate = certificate.NewCertificate(localLife.ctx)
	}
	return localLife.certificate
}

// GetDevice 获取Good
func (localLife *LocalLife) GetGood() *goods.Good {
	if localLife.goods == nil {
		localLife.goods = goods.NewGood(localLife.ctx)
	}
	return localLife.goods
}
