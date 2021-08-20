package ginsession

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName  = "session_id"
	SessionContextName = "session"
)

var (
	MgrObj Mgr
)

type Option struct {
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

type SessionData interface {
	GetID() string
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
	Del(key string)
	Save()
	SetExpire(int)
	Load(sessionID string) (err error)
}

type Mgr interface {
	Init(addr string, options ...string)
	GetSessionData(sessionId string) (sd SessionData, err error)
	CreateSession() (sd SessionData)
}

func InitMgr(name string, addr string, option ...string) {
	switch name {
	case "memory":
		MgrObj = NewMemoryMgr()
	case "redis":
		MgrObj = NewRedisMgr()
	}
	MgrObj.Init(addr, option...)
}

func SessionMiddleware(mgrObj Mgr, option *Option) gin.HandlerFunc {
	if mgrObj == nil {
		panic("must call InitMgr before use it.")
	}
	return func(c *gin.Context) {

		var sd SessionData // session data
		sessionID, err := c.Cookie(SessionCookieName)
		fmt.Println(sessionID)
		if err != nil {
			sd = mgrObj.CreateSession()
			sessionID = sd.GetID()

		} else {

			sd, err = mgrObj.GetSessionData(sessionID)
			if err != nil {
				sd = mgrObj.CreateSession()
				sessionID = sd.GetID()

			}

		}
		sd.SetExpire(option.MaxAge)
		c.Set(SessionContextName, sd)
		c.SetCookie(SessionCookieName, sessionID, option.MaxAge, option.Path, option.Domain, option.Secure, option.HttpOnly)
		c.Next()
	}
}
