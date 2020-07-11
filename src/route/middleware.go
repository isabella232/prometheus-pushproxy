package route

//middleware includes auth, rate limit, and etc.
import (
	"net/http"
	"strings"

	"github.com/kafkaesque-io/prometheus-pushproxy/src/util"

	log "github.com/sirupsen/logrus"
)

var (
	// Rate is the default global rate limit
	// This rate only limits the rate hitting on endpoint
	// It does not limit the underline resource access
	Rate = NewSema(200)
)

// AuthFunc is a function type to allow pluggable authentication middleware
type AuthFunc func(next http.Handler) http.Handler

// AuthVerifyAPIKey authenticates api key
func AuthVerifyAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1))
		if apiKey == util.GetConfig().DefaultAPIKey || util.GetConfig().DefaultAPIKey == "" {
			next.ServeHTTP(w, r)
		} else {
			log.Warnf("failed authentication with api-key %s", apiKey)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

// TODO: implement CRUD to create and revoke api Keys with the master key

// NoAuth bypasses the auth middleware
func NoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// LimitRate rate limites against http handler
// use semaphore as a simple rate limiter
func LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := Rate.Acquire()
		if err != nil {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		} else {
			next.ServeHTTP(w, r)
		}
		Rate.Release()
	})
}
