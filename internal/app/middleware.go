package app

// import (
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/lightsaid/ebook/pkg/logger"
// 	"github.com/lightsaid/ebook/pkg/random"
// )

// // logRequest 访问系统日志
// func (a *application) logRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		defer func() {
// 			// NOTE: Stauts header 在 a.writeJSON 写入
// 			logger.InfoLog.Printf("[INFO] %s\t%s\t%v\t%v", r.RequestURI, r.Method, time.Since(start), w.Header().Get("Stauts"))
// 		}()

// 		next.ServeHTTP(w, r)
// 	})
// }

// // recoverPanic 恐慌处理
// func (a *application) recoverPanic(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				a.errorResponse(w, fmt.Errorf("%v", r))
// 			}
// 		}()

// 		next.ServeHTTP(w, r)
// 	})
// }

// // authenticate 认证
// func (a *application) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.InfoLog.Println("-> Todo...")
// 		next.ServeHTTP(w, r)
// 	})
// }

// // rateLimit 请求速率
// func (a *application) rateLimit(next http.Handler) http.Handler {
// 	// 这里用最简单的方式实现，计数器形式，同一个IP在一定时间内允许特定访问次数
// 	// NOTE: 每分钟最大限流是max
// 	const max = 60
// 	type counter struct {
// 		count     int
// 		startSeen time.Time
// 	}

// 	var clients = make(map[string]counter)
// 	var mu sync.RWMutex

// 	go func() {
// 		for {
// 			// 每过一分钟清理一次
// 			time.Sleep(time.Minute)
// 			mu.Lock()
// 			for ip, client := range clients {
// 				if time.Since(client.startSeen) > time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			mu.Unlock()
// 		}
// 	}()

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// ip := r.RemoteAddr
// 		ip := "192.168.10" + fmt.Sprintf("%d", random.RandomInt(1, 3))
// 		mu.Lock()
// 		_, exists := clients[ip]
// 		if exists {
// 			client := clients[ip]
// 			client.count = client.count + 1
// 			clients[ip] = client
// 			if client.count > max {
// 				a.writeJSON(w, http.StatusTooManyRequests, msgWrapp("请求过多"))
// 				mu.Unlock()
// 				return
// 			}
// 		} else {
// 			clients[ip] = counter{startSeen: time.Now(), count: 1}
// 		}
// 		mu.Unlock()
// 		next.ServeHTTP(w, r)
// 	})
// }
