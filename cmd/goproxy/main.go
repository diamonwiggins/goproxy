package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var config ProxyConfig
var port *string

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

	fmt.Println("goproxy running on :", *port)

	http.HandleFunc("/proxyServer", ProxyServer)

	// server redirection
	http.HandleFunc("/", rootHandle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	lrw := NewLoggingResponseWriter(w)

	hostHeader := r.Host
	for i := range config.Hosts {
		if strings.Contains(hostHeader, config.Hosts[i].Host) {
			proxy := NewProxy(config.Hosts[i].Backend)
			lrw.Header().Set("X-GoProxy", "GoProxy")
			proxy.proxy.ServeHTTP(lrw, r)

			currentTime := time.Now().UTC()
			statusCode := lrw.statusCode
			fmt.Printf("%s - \"%s %s %s %d\" \"%s\" - %s\n", currentTime.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, r.Proto, statusCode, r.Header.Get("User-Agent"), r.RemoteAddr)
			return
		}
	}
}
