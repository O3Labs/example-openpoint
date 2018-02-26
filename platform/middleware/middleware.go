package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func UseWithHook(handler http.HandlerFunc, hook func(http.HandlerFunc) http.HandlerFunc, middlewares []func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	if hook != nil {
		handler = hook(handler)
	}

	for _, m := range middlewares {
		handler = m(handler)
	}

	return handler
}

// func Use(handler http.HandlerFunc, middlewares []func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
// 	for _, m := range middlewares {
// 		handler = m(handler)
// 	}
// 	return handler
// }

func AllowIFrame(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "ALLOWALL")
		h(w, r)
	}
}

func JavascriptResult(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-javascript; charset=utf-8")
		h(w, r)
	}
}

func GzipResult(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	}
}

// func CachedResult(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", 3600))
// 		h(w, r)
// 	}
// }

func EventStream(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		h(w, r)
	}
}

func JsonResult(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h(w, r)
	}
}

type LogRecord struct {
	http.ResponseWriter
	Status  int
	Message string
}

func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

func (r *LogRecord) WriteHeader(status int) {
	r.Status = status
	//error is always in json format.
	if status != 0 && status != 200 {
		r.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	r.ResponseWriter.WriteHeader(status)
}
