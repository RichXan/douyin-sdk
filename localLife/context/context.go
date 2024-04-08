package context

import (
	"github.com/RichXan/douyin-sdk/credential"
	"github.com/RichXan/douyin-sdk/localLife/config"
)

type Context struct {
	*config.Config
	credential.ClientTokenHandle
}
