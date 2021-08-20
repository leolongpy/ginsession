# ginsession

Download and install:

```shell
go get github.com/leolongpy/ginsession
```

Import it in you code:

```go
import "github.com/leolongpy/ginsession"
```

```go
// 测试Session服务的gin demo
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// 初始化全局的MgrObj对象
	ginsession.InitMgr("redis", "127.0.0.1:6379") // Redis版
	option := &ginsession.Option{
		MaxAge:   600,
		Path:     "/",
		Domain:   "127.0.0.1",
		Secure:   false,
		HttpOnly: true,
	}
    // session中间件应该作为一个全局的中间件
	r.Use(ginsession.SessionMiddleware(ginsession.MgrObj, option))
	r.Run()

}
```

```go
func AuthMiddleware(c *gin.Context) {
	tmpSD, _ := c.Get(ginsession.SessionContextName)
	sd := tmpSD.(ginsession.SessionData)
	value, err := sd.Get("isLogin")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/login")
		return
	}
	isLogin, ok := value.(bool)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if !isLogin {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Next()
}

func loginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		toPath := c.DefaultQuery("next", "/index")
		var u UserInfo
		err := c.ShouldBind(&u)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"err": "用户名或密码不能为空",
			})
			return
		}
		if u.Username == "leo" && u.Password == "123456" {
			tmpSD, ok := c.Get(ginsession.SessionContextName)
			if !ok {
				panic("session middleware")
			}
			sd := tmpSD.(ginsession.SessionData)
			sd.Set("isLogin", true)
			sd.Save()
			c.Redirect(http.StatusMovedPermanently, toPath)
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"err": "用户名或密码错误",
			})
			return
		}
	} else {
		c.HTML(http.StatusOK, "login.html", nil)
	}

}
```

