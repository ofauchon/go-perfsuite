package core

import (
    "fmt"
    "strconv"
    "net/http"

    "github.com/zenazn/goji/graceful"
    "github.com/zenazn/goji/web"
)

const DEFAULT_HTTP_PORT=8000

type ControllerInstance struct{
	Config *ControllerConfig
}

type ControllerConfig struct{
	HttpServerPort	int
}

func NewControllerInstance() (*ControllerInstance) {
	conf := &ControllerConfig{HttpServerPort: DEFAULT_HTTP_PORT}
	ctrl := &ControllerInstance{Config: conf}
	return ctrl
}

func (c *ControllerInstance) StartHttpServer(){
	fmt.Printf("Starting HTTP Server on port %d\n", c.Config.HttpServerPort)
    context := &appContext{controller:c }

     r := web.New()
    r.Get("/", appHandler{context, IndexHandler})

    graceful.ListenAndServe(":" + strconv.Itoa(c.Config.HttpServerPort), r)   
}


func IndexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
    // Our handlers now have access to the members of our context struct.
    // e.g. we can call methods on our DB type via err := a.db.GetPosts()
    fmt.Fprintf(w, "IndexHandler: HTTP Port is %d", a.controller.Config.HttpServerPort)
    return 200, nil
}

