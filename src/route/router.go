package route

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kafkaesque-io/pulsar-beam/src/middleware"
	log "github.com/sirupsen/logrus"
)

// NewRouter - create new router for HTTP routing
func NewRouter(mode *string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range ProxyRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.AuthFunc(handler))

	}
	// TODO rate limit can be added per route basis
	router.Use(middleware.LimitRate)

	log.Infof("router added")
	return router
}
