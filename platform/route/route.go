package route

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/o3labs/openpoint/platform/middleware"
)

func Use(middlewares ...func(http.HandlerFunc) http.HandlerFunc) []func(http.HandlerFunc) http.HandlerFunc {
	return middlewares
}

type Group struct {
	Path    string
	Configs []Config
}

type Config struct {
	Path        string
	Func        http.HandlerFunc
	Hook        func(http.HandlerFunc) http.HandlerFunc
	Middlewares []func(http.HandlerFunc) http.HandlerFunc
	Method      string
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type MyRoute struct {
	r *mux.Router
}

func Router(r *mux.Router) MyRoute {
	return MyRoute{r}
}

func NewGroup(path string) Group {
	g := Group{}
	g.Path = path
	return g
}

func (g *Group) Add(config Config) {
	config.Path = g.Path + config.Path
	g.Configs = append(g.Configs, config)
}

func (m MyRoute) Register(configs ...[]Config) {
	for _, config := range configs {
		for _, route := range config {
			m.r.HandleFunc(route.Path, middleware.UseWithHook(route.Func, route.Hook, route.Middlewares)).Methods(route.Method)
			log.Printf("Registered %+v %+v\n", route.Method, route.Path)
		}
	}
}
