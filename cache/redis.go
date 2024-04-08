package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis .redis cache
type Redis struct {
	ctx  context.Context
	conn *redis.Client
}

// RedisOptions redis 连接属性
type RedisOptions struct {
	Addr     string `yml:"addr" json:"addr"`
	Password string `yml:"password" json:"password"`
	Database int    `yml:"database" json:"database"`
}

// NewRedis 实例化
func NewRedis(ctx context.Context, opts *RedisOptions) *Redis {
	//初始化redis，连接地址和端口，密码，数据库名称
	conn := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.Database,
	})

	_, err := conn.Ping(ctx).Result()
	if err != nil {
		return nil
	}

	return &Redis{
		conn: conn,
		ctx:  context.Background(),
	}
}

// SetConn 设置conn
func (r *Redis) SetConn(conn *redis.Client) {
	r.conn = conn
}

// SetRedisCtx 设置redis ctx 参数
func (r *Redis) SetRedisCtx(ctx context.Context) {
	r.ctx = ctx
}

// Get 获取一个值
func (r *Redis) Get(key string) interface{} {
	return r.GetContext(r.ctx, key)
}

// GetContext 获取一个值
func (r *Redis) GetContext(ctx context.Context, key string) interface{} {
	cmd := redis.NewCmd(ctx, "get", key)
	_ = r.conn.Do(ctx, cmd)
	return cmd.Val()
}

// 获取
func (r *Redis) GetString(key string) (string, error) {
	return r.conn.Get(r.ctx, key).Result()
}

// 获取Int
func (r *Redis) GetInt(key string) (int, error) {
	return r.conn.Get(r.ctx, key).Int()
}

func (r *Redis) GetInt64(key string) (int64, error) {
	return r.conn.Get(r.ctx, key).Int64()
}

func (r *Redis) GetFloat64(key string) (float64, error) {
	return r.conn.Get(r.ctx, key).Float64()
}

// byte
func (r *Redis) GetByte(key string) ([]byte, error) {
	return r.conn.Get(r.ctx, key).Bytes()
}

// Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) error {
	return r.SetContext(r.ctx, key, val, timeout)
}

// SetContext 设置一个值
func (r *Redis) SetContext(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.conn.Set(r.ctx, key, val, timeout).Err()
}

// IsExist 判断key是否存在
func (r *Redis) IsExist(key string) bool {
	return r.IsExistContext(r.ctx, key)
}

// IsExistContext 判断key是否存在
func (r *Redis) IsExistContext(ctx context.Context, key string) bool {
	result, _ := r.conn.Exists(ctx, key).Result()

	return result > 0
}

// Delete 删除
func (r *Redis) Delete(key string) error {
	return r.DeleteContext(r.ctx, key)
}

// DeleteContext 删除
func (r *Redis) DeleteContext(ctx context.Context, key string) error {
	return r.conn.Del(ctx, key).Err()
}
