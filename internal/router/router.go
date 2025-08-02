package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"go.uber.org/zap"
)

type Router struct {
	ServeMux  *http.ServeMux
	Env       *env.Environment
	Endpoints []Endpoint
}

type Endpoint struct {
	HTTPMethod string
	Path       string
	Handler    http.Handler
}

func (endpoint *Endpoint) GetFullPath() string {
	return fmt.Sprintf("%s %s", endpoint.HTTPMethod, endpoint.Path)
}

func (endpoint *Endpoint) GetOptionsPath() string {
	return fmt.Sprintf("%s %s", http.MethodOptions, endpoint.Path)
}

func (router *Router) addEndpoint(endpoint Endpoint) {
	router.Endpoints = append(router.Endpoints, endpoint)
}

func (router *Router) AddEndpointWithPrefix(prefix string, endpoint Endpoint) {
	endpoint.Path = fmt.Sprintf("/%s/%s", prefix, endpoint.Path)
	router.addEndpoint(endpoint)
}

func (router *Router) InitEndpoints() {
	for _, endpoint := range router.Endpoints {
		router.ServeMux.Handle(endpoint.GetOptionsPath(),
			router.LogRequest(
				router.Verify(
					router.CorsHandler(),
				),
			),
		)

		router.ServeMux.Handle(endpoint.GetFullPath(),
			router.LogRequest(
				router.Verify(
					endpoint.Handler,
				),
			),
		)
	}
}

func (router *Router) Verify(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r)
		if token == router.Env.ApiToken() {
			originalHandler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized access"))
		}
	})
}

func (router *Router) LogRequest(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequest := map[string]interface{}{
			"requestMethod": strings.ToUpper(r.Method),
			"requestUrl":    r.URL.RequestURI(),
			"userAgent":     r.UserAgent(),
			//"remoteIp":      parseIP(r),
			"rawQuery": r.URL.RawQuery,
			"query":    r.URL.Query(),
			"host":     r.Host,
		}
		logger.GetLogInstance().Info(fmt.Sprintf("%s %s", r.Method, r.URL.String()), zap.Any("httpRequest", httpRequest))
		originalHandler.ServeHTTP(w, r)
	})
}

func getBearerToken(r *http.Request) string {
	authorizationValue := r.Header.Get("Authorization")
	values := strings.Split(authorizationValue, " ")
	if len(values) > 1 {
		return values[1]
	}
	return ""
}

func (router *Router) CorsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "authorization, content-type, accept, accept-language")
		w.Header().Add("Access-Control-Max-Age", "86400")
	})
}

func NewHandlerServer(router *Router) *http.ServeMux {
	router.InitEndpoints()
	return router.ServeMux
}
