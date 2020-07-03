package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kafkaesque-io/pulsar-beam/src/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Route - HTTP Route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	AuthFunc    mux.MiddlewareFunc
}

// Routes list of HTTP Routes
type Routes []Route

// ProxyRoutes definition
var ProxyRoutes = Routes{
	Route{
		"Prometeus metrics",
		http.MethodGet,
		"/metrics",
		promhttp.Handler().ServeHTTP,
		AuthVerifyAPIKey,
	},
	Route{
		"status",
		"GET",
		"/status",
		StatusPage,
		middleware.NoAuth,
	},
	Route{
		"Receive",
		"POST",
		"/v1/proxy",
		ReceiveHandler,
		AuthVerifyAPIKey,
	},
}
