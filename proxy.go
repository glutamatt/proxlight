package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	input, listen, err := parseArgs()
	check(err, "parseArgs")
	target, err := url.ParseRequestURI(input)
	check(err, "Url parse '"+input+"'")
	check(nil, fmt.Sprintf("Proxiing %s through %s", target, listen))
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorLog = log.New(os.Stdin, "[ERROR PROXY] ", log.LstdFlags)

	check(http.ListenAndServe(listen, enableCors(forgeHost(target.Host, throttle(10, proxy)))), "http.ListenAndServe")
}

func check(err error, step string) {
	defer log.SetPrefix("")
	if err == nil {
		log.SetPrefix("[OK] ")
		log.Println(step)
	} else {
		log.SetPrefix("[ERROR] ")
		log.Fatalf("%s : %s", step, err.Error())
	}
}

func parseArgs() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf("usage : ./proxy http://domain:port/path 0.0.0.0:1234")
	}
	return os.Args[1], os.Args[2], nil
}

func enableCors(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(resp, req)
	})
}

func forgeHost(hostname string, handler http.Handler) http.HandlerFunc {
	log.Println("Forge request host with", hostname)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		req.Host = hostname
		req.URL.Host = hostname
		req.Header.Set("Host", hostname)
		handler.ServeHTTP(resp, req)
	})
}

func throttle(reqPerSec int, handler http.Handler) http.HandlerFunc {
	ticker := time.NewTicker(time.Second / time.Duration(reqPerSec))
	log.Printf("Throttle %d req per sec", reqPerSec)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		<-ticker.C
		handler.ServeHTTP(resp, req)
	})
}
