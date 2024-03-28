package context

import (
	"github.com/RichXan/douyin-sdk/open/credential"
	"github.com/RichXan/douyin-sdk/open/localLife/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
