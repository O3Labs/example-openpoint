package app

import (
	"net/http"

	"github.com/o3labs/openpoint/platform/helper"
)

type ControllerInterface interface {
	Root(w http.ResponseWriter, r *http.Request)
}

type Controller struct{}

func NewController() ControllerInterface {
	return &Controller{}
}

func (c *Controller) Root(w http.ResponseWriter, r *http.Request) {
	template := helper.LoadTemplateWithContentFile("web/app/index.html")

	out := struct {
		App string
	}{
		App: "Open Point",
	}

	//sample call
	// neoclient := neorpc.NewClient("http://localhost:30333")
	// result := neoclient.GetContractState("ce575ae1bb6153330d20c560acb434dc5755241b")
	// log.Printf("%v", result)
	template.Execute(w, out)
}
