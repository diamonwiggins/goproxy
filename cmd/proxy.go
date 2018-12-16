package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var config ProxyConfig
var port *string

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

func handle(w http.ResponseWriter, r *http.Request) {
	hostHeader := r.Host
	for i := range config.Hosts {
		if strings.Contains(hostHeader, config.Hosts[i].Host) {
			proxy := NewProxy(config.Hosts[i].Backend)
			w.Header().Set("X-GoProxy", "GoProxy")
			proxy.proxy.ServeHTTP(w, r)
			return
		}
	}
	w.Write([]byte("403: Host forbidden " + hostHeader))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Reverse proxy Server Running. Accepting at port:" + *port))
}

func main() {
	const (
		defaultPortUsage   = "default server port, ':9000'"
		defaultConfig      = "config.json"
		defaultConfigUsage = "default config path, './config.json'"
	)

	var configFile *string

	// flags
	port = flag.String("port", "9000", defaultPortUsage)
	configFile = flag.String("config", defaultConfig, defaultConfigUsage)
	flag.Parse()

	//read config file
	jsonFile, err := os.Open(*configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	fmt.Println("server will run on :", *port)

	http.HandleFunc("/proxyServer", ProxyServer)

	// server redirection
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
