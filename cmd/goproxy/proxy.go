package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	backend *url.URL
	proxy   *httputil.ReverseProxy
}
type Host struct {
	Host    string `json:"host"`
	Backend string `json:"backend"`
}
type ProxyConfig struct {
	DefaultPort string `json:"defaultPort"`
	Hosts       []Host `json:"hosts"`
}

func NewProxy(backend string) *Proxy {
	url, _ := url.Parse(backend)

	return &Proxy{backend: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}
func ProxyServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Reverse proxy Server Running. Accepting at port:" + *port))
}
