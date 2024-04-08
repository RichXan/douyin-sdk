package goods

import "github.com/RichXan/douyin-sdk/localLife/context"

const (
	getOnlineProductsURL = "https://open.douyin.com/goodlife/v1/goods/product/online/query"
)

type Good struct {
	*context.Context
}

// NewGood 实例化
func NewGood(context *context.Context) *Good {
	user := new(Good)
	user.Context = context
	return user
}
