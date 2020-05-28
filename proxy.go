package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	withCors := flag.Bool("cors", false, "enable cors header \"Access-Control-Allow-Origin\":*")
	throttleRPS := flag.Int("throttle", 0, "throttle to n requests per second")
	input, listen, err := parseArgs()
	check(err, "parseArgs")
	target, err := url.ParseRequestURI(input)
	check(err, "Url parse '"+input+"'")
	check(nil, fmt.Sprintf("Proxiing %s through %s", target, listen))
	var handler http.Handler
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorLog = log.New(os.Stdin, "[ERROR PROXY] ", log.LstdFlags)
	handler = proxy
	if *throttleRPS > 0 {
		handler = throttle(*throttleRPS, handler)
	}
	handler = forgeHost(target.Host, handler)
	if *withCors {
		handler = enableCors(handler)
	}
	check(http.ListenAndServe(listen, handler), "http.ListenAndServe")
}

func check(err error, step string) {
	defer log.SetPrefix("")
	if err == nil {
		log.Printf("%s\n", step)
	} else {
		log.Fatalf("[OK] %s: %s", step, err)
	}
}

func parseArgs() (string, string, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		return "", "", fmt.Errorf("usage : ./proxy [-h for flags help] http://domain:port/path 0.0.0.0:1234")
	}
	return args[0], args[1], nil
}

func enableCors(handler http.Handler) http.HandlerFunc {
	log.Println("Cors headers enabled")
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
