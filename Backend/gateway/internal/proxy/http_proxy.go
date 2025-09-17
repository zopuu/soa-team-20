package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Options struct {
	Target       string
	StripPrefix  string
	DialTimeout  time.Duration
	ProxyTimeout time.Duration
}

func NewHTTPReverseProxy(opt Options) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(opt.Target)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{Timeout: opt.DialTimeout}
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: opt.ProxyTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false}, // change if needed
		Proxy:                 http.ProxyFromEnvironment,
	}

	rp := httputil.NewSingleHostReverseProxy(u)
	// Custom director to strip prefix and fix Host headers
	origDirector := rp.Director
	rp.Director = func(req *http.Request) {
		origDirector(req)
		if opt.StripPrefix != "" && strings.HasPrefix(req.URL.Path, opt.StripPrefix) {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, opt.StripPrefix)
			if req.URL.Path == "" { req.URL.Path = "/" }
		}
		// Forward original host
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
		req.Host = u.Host
	}
	rp.Transport = transport

	// Optional error handler
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		http.Error(w, "Upstream unavailable", http.StatusBadGateway)
	}
	rp.ModifyResponse = func(res *http.Response) error {
		// skini CORS headere koje je setovao downstream, gateway će ih postaviti
		res.Header.Del("Access-Control-Allow-Origin")
		res.Header.Del("Access-Control-Allow-Headers")
		res.Header.Del("Access-Control-Allow-Methods")
		res.Header.Del("Access-Control-Expose-Headers")
		res.Header.Del("Access-Control-Allow-Credentials")
		return nil
	}

	return rp, nil
}
