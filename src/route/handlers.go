package route

import (
	"io/ioutil"
	"net/http"
)

const subDelimiter = "-"

// Init initializes database
func Init() {
}

// StatusPage replies with basic status code
func StatusPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

// ReceiveHandler - the message receiver handler
func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		util.ResponseErrorJSON(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// ResponseErr - Error struct for Http response
type ResponseErr struct {
	Error string `json:"error"`
}
