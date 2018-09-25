package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	input, listen, err := parseArgs()
	check(err, "parseArgs")
	target, err := url.ParseRequestURI(input)
	check(err, "Url parse '"+input+"'")
	check(nil, fmt.Sprintf("Proxiing %s through %s", target, listen))
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorLog = log.New(os.Stdin, "[ERROR PROXY] ", log.LstdFlags)
	check(http.ListenAndServe(listen, proxy), "http.ListenAndServe")
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
