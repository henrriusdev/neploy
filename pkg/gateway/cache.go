package gateway

import (
	"bytes"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var memoryCache = cache.New(5*time.Minute, 10*time.Minute)

const maxCacheSize = 5 * 1024 * 1024 // 5MB

func SetCache(key string, value any, ttl time.Duration) {
	memoryCache.Set(key, value, ttl)
}

func GetCache(key string) (any, bool) {
	return memoryCache.Get(key)
}

func InvalidateCache(key string) {
	memoryCache.Delete(key)
}

func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + ":" + r.URL.String()

		if val, found := GetCache(key); found {
			w.Header().Set("X-Cache", "HIT")
			w.Write(val.([]byte))
			return
		}

		buf := new(bytes.Buffer)
		cw := &captureWriter{ResponseWriter: w, buf: buf, status: http.StatusOK}

		next.ServeHTTP(cw, r)

		if cw.status == http.StatusOK && buf.Len() <= maxCacheSize {
			SetCache(key, buf.Bytes(), 1*time.Minute)
		}
	})
}

type captureWriter struct {
	http.ResponseWriter
	buf    *bytes.Buffer
	status int
}

func (w *captureWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *captureWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}
