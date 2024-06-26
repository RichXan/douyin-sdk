package douyinsdk

import (
	"net/http"
	"os"

	"github.com/RichXan/douyin-sdk/cache"
	"github.com/RichXan/douyin-sdk/localLife"
	"github.com/RichXan/douyin-sdk/localLife/config"
	"github.com/RichXan/douyin-sdk/util"
	log "github.com/sirupsen/logrus"
)

type Douyin struct {
	cache cache.Cache
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func NewDouyin() *Douyin {
	return &Douyin{}
}

// SetCache 设置 cache
func (dy *Douyin) SetCache(cache cache.Cache) {
	dy.cache = cache
}

func (dy *Douyin) GetLocalLife(cfg *config.Config) *localLife.LocalLife {
	if cfg.Cache == nil {
		cfg.Cache = dy.cache
	}
	return localLife.NewLocalLife(cfg)
}

func (dy *Douyin) SetHTTPClient(client *http.Client) {
	util.DefaultHTTPClient = client
}
