package obs

import "net/http"

// Transport koji automatski propagira korelacione headere ka downstream servisima.
type CorrTransport struct{ Base http.RoundTripper }

func (t CorrTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.Base == nil {
		t.Base = http.DefaultTransport
	}
	if r.Header.Get("X-Request-ID") == "" {
		if id := ReqIDFrom(r); id != "" {
			r.Header.Set("X-Request-ID", id)
		}
	}
	if r.Header.Get("X-Trace-Id") == "" {
		if id := TraceIDFrom(r); id != "" {
			r.Header.Set("X-Trace-Id", id)
		}
	}
	return t.Base.RoundTrip(r)
}

func NewHTTPClient() *http.Client {
	return &http.Client{Transport: CorrTransport{Base: http.DefaultTransport}}
}
