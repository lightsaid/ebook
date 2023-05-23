package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lightsaid/ebook/pkg/logger"
)

// logRequest 访问系统日志
func (a *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			// NOTE: Stauts header 在 a.writeJSON 写入
			logger.InfoLog.Printf("[INFO] %s\t%s\t%v\t%v", r.RequestURI, r.Method, time.Since(start), w.Header().Get("Stauts"))
		}()

		next.ServeHTTP(w, r)
	})
}

// recoverPanic 恐慌处理
func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				a.errorResponse(w, fmt.Errorf("%v", r))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// authenticate 认证
func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLog.Println("-> Todo...")
		next.ServeHTTP(w, r)
	})
}
