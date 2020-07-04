package route

import (
	"net/http"

	"github.com/gorilla/mux"
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
		"Application Prometeus metrics",
		http.MethodGet,
		"/metrics",
		promhttp.Handler().ServeHTTP,
		AuthVerifyAPIKey,
	},
	Route{
		"Receive",
		http.MethodPost,
		"/v1/proxy/{instance}",
		ReceiveHandler,
		AuthVerifyAPIKey,
	},
	Route{
		"All proxy Prometheus metrics",
		http.MethodGet,
		"/proxy-metrics",
		ProxyMetricsHandler,
		AuthVerifyAPIKey,
	},
	Route{
		"Individual instance proxy Prometheus metrics",
		http.MethodGet,
		"/proxy-metrics/{instance}",
		ProxyMetricsHandler,
		AuthVerifyAPIKey,
	},
}
