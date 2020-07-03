package util

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ResponseErr - Error struct for Http response
type ResponseErr struct {
	Error string `json:"error"`
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// JoinString joins multiple strings
func JoinString(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// ResponseErrorJSON builds a Http response.
func ResponseErrorJSON(e error, w http.ResponseWriter, statusCode int) {
	response := ResponseErr{e.Error()}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

// AssignString returns the first non-empty string
// It is equivalent the following in Javascript
// var value = val0 || val1 || val2 || default
func AssignString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// ReportError logs error
func ReportError(err error) error {
	log.Errorf("error %v", err)
	return err
}

// StrContains check if a string is contained in an array of string
func StrContains(strs []string, str string) bool {
	for _, v := range strs {
		if strings.TrimSpace(v) == strings.TrimSpace(str) {
			return true
		}
	}
	return false
}

// StrToInt converts string to an integer format with a default if inproper value retrieved
func StrToInt(str string, defaultNum int) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}
	return defaultNum
}

// StringToBool format various strings to boolean
// strconv.ParseBool only covers `true` and `false` cases
func StringToBool(str string) bool {
	s := strings.ToLower(strings.TrimSpace(str))

	// `1` is true because the default Golang boolean is initialized as false
	if s == "true" || s == "yes" || s == "enable" || s == "enabled" || s == "1" || s == "ok" {
		return true
	}

	return false
}

// QueryParamString get URL query parameter with a default value
func QueryParamString(params url.Values, name, defaultValue string) string {
	if str, ok := params[name]; ok {
		return str[0]
	}
	return defaultValue
}
