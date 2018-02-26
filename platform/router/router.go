package router

import (
	"github.com/o3labs/openpoint/platform/controllers/web"

	"github.com/o3labs/openpoint/platform/controllers/api"
	"github.com/o3labs/openpoint/platform/route"

	"github.com/gorilla/mux"
)

func RegisterRouting(r *mux.Router) {
	route.Router(r).Register(
		api.Routes(),
		web.Routes(),
	)
}
