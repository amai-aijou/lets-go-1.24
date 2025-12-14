package main

import (
	"fmt" // 6.4
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// 6.3: logRequest() method to capture IP address, method, URI, and HTTP version for requests made
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip		= r.RemoteAddr
			proto	= r.Proto
			method	= r.Method
			uri		= r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

// 6.4: Allow us to recover from panics and restart the application; also, notify users of the panic via HTTP 500
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 6.4: Create deferred function (which will always be run in event of panic)
		defer func() {
			// 6.4: Use the built-in recover() function to check for panic and return panic value
			pv := recover()

			if pv != nil {
				//6.4: If panic happened, set "Connection: close' header on the response
				w.Header().Set("Connection", "close")
				//6.4: Call the app.serverError helper method to return a 500 Internal Server response
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
