package api

import (
	"github.com/o3labs/openpoint/platform/controllers/api/payment"
	"github.com/o3labs/openpoint/platform/route"
)

func Routes() []route.Config {
	paymentController := payment.NewController()

	return []route.Config{
		route.Config{
			Path:   "/api/v1/pay",
			Func:   paymentController.ProcessPaymentToken,
			Method: route.POST,
		},
	}
}
