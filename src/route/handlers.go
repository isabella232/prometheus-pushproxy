package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kafkaesque-io/prometheus-pushproxy/src/util"
)

// MetricsCache is the cache for Producer objects
var MetricsCache = util.NewCache(util.CacheOption{
	TTL:            time.Duration(300) * time.Second, //TODO: add configurable TTL
	CleanInterval:  time.Duration(302) * time.Second,
	ExpireCallback: func(key string, value interface{}) {},
})

// Init initializes database
func Init() {
}

// HealthHandler replies with basic status code
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

// ReceiveHandler - the message receiver handler
func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	instance := mux.Vars(r)["instance"]
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		util.ResponseErrorJSON(err, w, http.StatusInternalServerError)
		return
	}

	ttl := "300" // default seconds
	if keys, ok := r.URL.Query()["ttl"]; ok {
		ttl = keys[0]
	}

	MetricsCache.SetWithTTL(instance, string(bytes), time.Duration(util.StrToInt(ttl, 0))*time.Second)
	w.WriteHeader(http.StatusOK)
	return
}

// ProxyMetricsHandler exposes received metrics
func ProxyMetricsHandler(w http.ResponseWriter, r *http.Request) {
	instance := mux.Vars(r)["instance"]
	var data string

	if instance != "" {
		obj, ok := MetricsCache.Get(instance)
		if !ok {
			util.ResponseErrorJSON(fmt.Errorf("instance %s not found", instance), w, http.StatusNotFound)
			return
		}
		data = fmt.Sprintf("%v", obj)
	} else {
		var err error
		data, err = getProxyMetrics()
		if err != nil {
			util.ResponseErrorJSON(err, w, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.Write([]byte(data))
	w.WriteHeader(http.StatusOK)
	return
}

func getProxyMetrics() (string, error) {
	var rc string
	for _, v := range MetricsCache.It() {
		rc = rc + fmt.Sprintf("%v", v.Data) + "\n"
	}

	return strings.TrimSuffix(rc, "\n"), nil
}
