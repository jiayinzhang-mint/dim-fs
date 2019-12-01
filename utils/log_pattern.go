package utils

import (
	"log"
	"net/http"
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
)

// LogError for err
func LogError(msg interface{}) {
	msgString, _ := InterfaceToString(msg)
	log.Printf(errorColor, "[ERR] "+msgString)
}

// LogInfo for info
func LogInfo(msg interface{}) {
	msgString, _ := InterfaceToString(msg)
	log.Printf(infoColor, "[INFO] "+msgString)
}

// LogDebug for denig
func LogDebug(msg interface{}) {
	msgString, _ := InterfaceToString(msg)
	log.Printf(debugColor, "[CALL] "+msgString)
}

// LogRequest for server log
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientIP := r.RemoteAddr
		method := r.Method
		url := r.RequestURI
		LogDebug(clientIP + " " + method + " " + url)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
