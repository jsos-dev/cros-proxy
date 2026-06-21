package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var excludeHeaders = map[string]bool{
	"host":              true,
	"cf-connecting-ip":  true,
	"cf-ray":            true,
	"cf-visitor":        true,
	"cf-ipcountry":      true,
	"x-forwarded-proto": true,
	"x-real-ip":         true,
	"x-cors-proxy-key":  true,
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("access-control-allow-methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		w.Header().Set("access-control-allow-headers", "*")
		w.Header().Set("access-control-expose-headers", "*")
		w.Header().Set("access-control-max-age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func corsProxy(proxyKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if proxyKey != "" {
			requestKey := r.Header.Get("x-cors-proxy-key")
			if requestKey == "" || requestKey != proxyKey {
				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Unauthorized: Missing or invalid x-cors-proxy-key header")
				return
			}
		}

		pathWithProtocol := strings.TrimPrefix(r.RequestURI, "/")

		if pathWithProtocol == "" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Usage: http://proxy/https://example.com/?query=value")
			return
		}

		var targetURL string
		if strings.HasPrefix(pathWithProtocol, "http://") || strings.HasPrefix(pathWithProtocol, "https://") {
			targetURL = pathWithProtocol
			if r.URL.RawQuery != "" {
				targetURL += "?" + r.URL.RawQuery
			}
		} else {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid target URL. Must start with http:// or https://")
			return
		}

		proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
		if err != nil {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Proxy Error: %v", err)
			return
		}

		for key, values := range r.Header {
			if !excludeHeaders[strings.ToLower(key)] {
				for _, value := range values {
					proxyReq.Header.Set(key, value)
				}
			}
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(proxyReq)
		if err != nil {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Proxy Error: %v", err)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("access-control-allow-methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		w.Header().Set("access-control-allow-headers", "*")
		w.Header().Set("access-control-expose-headers", "*")

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func main() {
	port := flag.String("port", "8080", "Server port")
	key := flag.String("key", "", "API authentication key (empty to disable)")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" && !isFlagSet("port") {
		*port = envPort
	}
	if envKey := os.Getenv("CORS_PROXY_KEY"); envKey != "" && !isFlagSet("key") {
		*key = envKey
	}

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: corsMiddleware(corsProxy(*key)),
	}

	log.Printf("CORS Proxy server starting on :%s", *port)
	if *key != "" {
		log.Printf("Authentication: enabled")
	} else {
		log.Printf("Authentication: disabled")
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
