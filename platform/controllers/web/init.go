package web

import (
	"github.com/o3labs/openpoint/platform/controllers/web/app"
	"github.com/o3labs/openpoint/platform/route"
)

func Routes() []route.Config {
	appController := app.NewController()

	return []route.Config{
		route.Config{
			Path:   "/web/",
			Func:   appController.Root,
			Method: route.GET,
		},
	}
}
