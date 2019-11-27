package plugins

import (
	"fmt"
	"io"
	"time"

	"github.com/ofauchon/loadwizard/core"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

type AdminServer struct {
	inj *core.Injector
}

func NewAdminServer(pInj *core.Injector) *AdminServer {
	a := &AdminServer{}
	a.inj = pInj
	return a
}

func (a *AdminServer) StartServer() {
	fmt.Println("Starting Admin Server")

	gin.DefaultWriter = colorable.NewColorableStderr()
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.GET("/stream", func(c *gin.Context) {
		chanStream := make(chan int64, 10)
		go func() {
			defer close(chanStream)
			a := int64(len(a.inj.Users))
			chanStream <- a
			time.Sleep(time.Second * 1)
		}()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Run()
}
