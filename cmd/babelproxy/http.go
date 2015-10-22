package main

import (
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

var (
	// set up connection and transport for HTTP
	myTransport *http.Transport
	myClient    *http.Client
	myRequests  chan CReq
)

// main calls initHttp to set up the HTTP transport and HTTP client
// to have keep-alives on, pool connections, and use a timeout
func initHttp() context.CancelFunc {
	myTransport = &http.Transport{DisableKeepAlives: false, MaxIdleConnsPerHost: conf.PoolSize, DisableCompression: true, ResponseHeaderTimeout: time.Duration(conf.Timeout)}
	myClient = &http.Client{Transport: myTransport}
	myRequests = make(chan CReq)
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < conf.PoolSize; i++ {
		go postThread(ctx, myRequests)
	}
	return cancel
}

type CReq struct {
	req *http.Request
	w   http.ResponseWriter
	rsp chan error
}

func doHttp(req *http.Request, w http.ResponseWriter) error {
	r := CReq{req: req, w: w, rsp: make(chan error)}
	myRequests <- r
	s := <-r.rsp
	return s
}

func postThread(ctx context.Context, in <-chan CReq) {
	done := ctx.Done()
	for {
		select {
		case r, ok := <-in:
			if !ok {
				return
			}
			resp, err := myClient.Do(r.req)
			if err != nil {
				http.Error(r.w, err.Error(), 500)
			} else {
				// TODO: What about defined response codes?
				// TODO: Other response types?
				for k, v := range resp.Header {
					for _, vi := range v {
						r.w.Header().Add(k, vi)
					}
				}
				r.w.WriteHeader(resp.StatusCode)
				//if resp.ContentLength > 0 {
				io.Copy(r.w, resp.Body)
				//}
				resp.Body.Close()
			}
			r.rsp <- err
		case <-done:
			return
		}
	}
}
