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

var port *string
var configFile *string
var config ProxyConfig

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}
type Host struct {
	Host   string `json:"host"`
	Target string `json:"target"`
}
type ProxyConfig struct {
	DefaultPort string `json:"defaultPort"`
	Hosts       []Host `json:"hosts"`
}

func NewProxy(target string) *Prox {
	url, _ := url.Parse(target)

	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func readConfig(fileName string) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
}

func handle(w http.ResponseWriter, r *http.Request) {
	hostHeader := r.Host
	for i := range config.Hosts {
		if strings.Contains(config.Hosts[i].Host, hostHeader) {
			proxy := NewProxy(config.Hosts[i].Target)
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

	// flags
	port = flag.String("port", "80", defaultPortUsage)
	configFile = flag.String("config", defaultConfig, defaultConfigUsage)
	flag.Parse()

	readConfig(*configFile)

	fmt.Println("server will run on :", *port)

	http.HandleFunc("/proxyServer", ProxyServer)

	// server redirection
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
