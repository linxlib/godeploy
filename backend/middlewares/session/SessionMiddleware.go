package session

import (
	"github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/memory"
	"github.com/fasthttp/session/v2/providers/redis"
	"github.com/linxlib/conv"
	"github.com/linxlib/fw"
	"time"
)

func NewSessionMiddleware(providerName string, config map[string]any) fw.IMiddlewareGlobal {
	s := &SessionMiddleware{
		MiddlewareCtl: fw.NewMiddlewareCtl("Session", "Session"),
		providerName:  providerName,
		config:        config,
	}
	encoder := session.Base64Encode
	decoder := session.Base64Decode
	var provider session.Provider
	var err error
	switch providerName {
	case "memory":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		provider, err = memory.New(memory.Config{})
	case "redis":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		if config == nil {
			panic("config is nil for redis provider")
		}
		var kp string
		var addr string
		var poolSize int
		var ConnMaxIdleTime time.Duration
		if tmp, ok := config["Redis_KeyPrefix"]; ok {
			kp = conv.String(tmp)
		} else {
			kp = "session"
		}
		if tmp, ok := config["Redis_Addr"]; ok {
			addr = conv.String(tmp)
		} else {
			addr = "127.0.0.1:6379"
		}
		if tmp, ok := config["Redis_PoolSize"]; ok {
			poolSize = conv.Int(tmp)
		} else {
			poolSize = 8
		}
		if tmp, ok := config["Redis_ConnMaxIdleTime"]; ok {
			ConnMaxIdleTime, _ = time.ParseDuration(conv.String(tmp))
		} else {
			ConnMaxIdleTime = 30 * time.Second
		}

		provider, err = redis.New(redis.Config{
			KeyPrefix:       kp,
			Addr:            addr,
			PoolSize:        poolSize,
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
	return s
}

var _ fw.IMiddlewareCtl = (*SessionMiddleware)(nil)

type IAuth interface {
	CheckSession(store *session.Store) bool
}

type SessionMiddleware struct {
	*fw.MiddlewareCtl
	providerName string
	config       map[string]any
	session      *session.Session
}

func (s *SessionMiddleware) CloneAsMethod() fw.IMiddlewareMethod {
	return s.CloneAsCtl()
}
func (s *SessionMiddleware) HandlerIgnored(next fw.HandlerFunc) fw.HandlerFunc {
	return func(context *fw.Context) {
		store, _ := s.session.Get(context.GetFastContext())
		context.Map(store)
		next(context)
		if len(store.GetAll().KV) > 0 {
			s.session.Save(context.GetFastContext(), store)
		}
	}
}
func (s *SessionMiddleware) HandlerMethod(next fw.HandlerFunc) fw.HandlerFunc {

	return func(context *fw.Context) {
		store, _ := s.session.Get(context.GetFastContext())

		if ctl, ok := s.GetCtlRValue().Interface().(IAuth); ok {
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
		context.Map(store)
		next(context)
		s.session.Save(context.GetFastContext(), store)

	}
}

func (s *SessionMiddleware) CloneAsCtl() fw.IMiddlewareCtl {
	return NewSessionMiddleware(s.providerName, s.config)
}

func (s *SessionMiddleware) HandlerController(base string) []*fw.RouteItem {
	return fw.EmptyRouteItem(s)
}
