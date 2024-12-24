package main

import (
	"net/http"
	"strings"
)

// 定义一个授权中间件
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取请求者的IP地址
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		// 检查IP地址是否为在白名单内
		for _, allowIP := range cfg.AllowIPs {
			if ip == allowIP || strings.HasPrefix(ip, allowIP+":") {
				next.ServeHTTP(w, r)
				return
			}
		}

		// 这里只是一个示例，实际应用中可能需要更复杂的验证逻辑
		authToken := r.Header.Get("Authorization")
		if authToken != "Bearer valid_token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
