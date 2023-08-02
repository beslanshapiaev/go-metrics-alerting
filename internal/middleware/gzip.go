package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			gzwr := gzipResponseWriter{ResponseWriter: w, Writer: gzw}
			next.ServeHTTP(gzwr, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
