package session

import (
	"github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/memory"
	"github.com/fasthttp/session/v2/providers/redis"
	"github.com/linxlib/fw"
	"time"
)

type Options struct {
	Provider string       `yaml:"provider" default:"memory"`
	Memory   MemoryConfig `yaml:"memory"`
	Redis    RedisConfig  `yaml:"redis"`
}

type MemoryConfig struct{}
type RedisConfig struct {
	KeyPrefix       string        `yaml:"keyPrefix" default:"session_"`
	Addr            string        `yaml:"addr" default:"127.0.0.1:6379"`
	Username        string        `yaml:"username" default:""`
	Password        string        `yaml:"password" default:""`
	DB              int           `yaml:"db" default:"0"`
	PoolSize        int           `yaml:"poolSize" default:"8"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" default:"5m"`
}

func NewSessionMiddleware() fw.IMiddlewareGlobal {
	s := &SessionMiddleware{
		MiddlewareCtl: fw.NewMiddlewareCtl("Session", "Session"),
	}

	return s
}

var _ fw.IMiddlewareCtl = (*SessionMiddleware)(nil)

type IAuth interface {
	CheckSession(store *session.Store) bool
}

type SessionMiddleware struct {
	*fw.MiddlewareCtl
	options *Options
	session *session.Session
}

// DoInitOnce is called once when the application starts.
// It loads the config for this middleware and sets up the session store.
// It will panic if the provider is invalid or if there is an error setting up the provider.
func (s *SessionMiddleware) DoInitOnce() {
	s.options = new(Options)
	s.LoadConfig("session", s.options)

	encoder := session.Base64Encode
	decoder := session.Base64Decode
	var provider session.Provider
	var err error
	switch s.options.Provider {
	case "memory":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		provider, err = memory.New(memory.Config{})
	case "redis":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode

		var kp string = s.options.Redis.KeyPrefix
		var addr string = s.options.Redis.Addr
		var poolSize int = s.options.Redis.PoolSize
		var db int = s.options.Redis.DB
		var ConnMaxIdleTime time.Duration = s.options.Redis.ConnMaxIdleTime

		provider, err = redis.New(redis.Config{
			KeyPrefix:       kp,
			Addr:            addr,
			PoolSize:        poolSize,
			DB:              db,
			ConnMaxIdleTime: ConnMaxIdleTime,
		})
	default:
		panic("Invalid provider")
	}
	if err != nil {
		panic(err)
	}
	cfg := session.NewDefaultConfig()
	cfg.EncodeFunc = encoder
	cfg.DecodeFunc = decoder
	s.session = session.New(cfg)
	if err = s.session.SetProvider(provider); err != nil {
		panic(err)
	}
}

// Execute will check the session and if unauthorized, return 401
// If authorized, it will set the session to the context and call the next handler
// If the handler is ignored, it will only save the session if it contains data
// Otherwise, it will always save the session
func (s *SessionMiddleware) Execute(ctx *fw.MiddlewareContext) fw.HandlerFunc {
	return func(context *fw.Context) {
		//fmt.Printf("%s.%s\n", ctx.ControllerName, ctx.MethodName)
		store, _ := s.session.Get(context.GetFastContext())

		if !ctx.Ignored {
			if ctl, ok := ctx.GetRValue().Interface().(IAuth); ok {
				if !ctl.CheckSession(store) {
					context.JSON(401, map[string]interface{}{
						"code":    401,
						"message": "session unauthorized",
						"data":    nil,
					})
					store.Flush()
					s.session.Save(context.GetFastContext(), store)
					return
				}
			}
		}

		context.Map(store)
		ctx.Next(context)
		if ctx.Ignored {
			if len(store.GetAll().KV) > 0 {
				s.session.Save(context.GetFastContext(), store)
			}
		} else {
			s.session.Save(context.GetFastContext(), store)
		}
	}
}
