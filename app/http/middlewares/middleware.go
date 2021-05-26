package middlewares

import "net/http"

// HttpHandleFunc 简写 - func(w http.ResponseWriter, r *http.Request)
type HttpHandleFunc func(w http.ResponseWriter, r *http.Request)
