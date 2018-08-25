package main

import (
	"github.com/golang/snappy"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	c_ErrorJson = "-BRQewogICJzd2FnZ2VyIjogIjIuMCIsARQYaW5mbyI6IAEgHCAgInRpdGxlASA0TXkgQXBwbGljYXRpb24FKyQgICJkZXNjcmlwBRN0OiAiLi4vLi4vYmFiZWx0ZW1wbGF0ZXMvZXJyb3IuBRUIOlxuCQIEIEIBJAwncyBFAR0YIEZvcm1hdAEbWCBDb3B5cmlnaHQgKEMpIDIwMTUgVGhlCTAcIEF1dGhvcnMBKVxHZW5lcmF0ZWQgb24gRnJpLCAyNCBBdWcBM0Q4IDIxOjQyOjQ2IEVEVCB2aWEJeQB0BYUAMg32CCBcIn6wAARcIhHpDHZlcnMB-Cg6ICIxLjAiCiAgfSUuDGhvc3QBFRBsb2NhbAUNBRccYmFzZVBhdGgBGwAvCRM0Y29uc3VtZXMiOiBbCiAhPQBhOVAQL2pzb24BVgBdBT8UcHJvZHVjmiwACGF0aAFVAHsJkRRkZWZpbmkhkQUVCWolZjXVIdcEeXAl1hRvYmplY3QFoQEBPtABLZ8FVCxlcyBhIHNpbmdsZSAl01wgbWVzc2FnZSBhbmQgY29kZSB0aGF0IG0lvQhiZSAlExRpemVkXG4FI1RkaXNwbGF5ZWQgdG8gYSBjYWxsZXIuGXwccHJvcGVydGkF9wm5AQEMIkNvZAGsGRIdwBRzdHJpbmcVRAEBQsQAQT88c2VydmljZS1zcGVjaWZpYw3CAbYhdgkBJUEJAQQiTQncPTQBAQAiMTieeAAkdGV4dCBvZiB0aDE1MGluIFVTLUVuZ2xpc2hWfAAQUGFyYW05vgkBFXsQYXJyYXk28gAIaXRlQjAARiUBRQsJAR3qRsUDJTEIbGlzBbkAcAGMEGV0ZXJzIbYZxy38LC4gVGhpcyBjb3VsZCH0IHVzZWQgYnlcbjH-hS4MIHN5cwGjAUMAZ224EUIEcyBhLGXGGV4hmAAuGbEAfQ27CQgNvQQiUynU4twCLkAAVeMZhBxyZXNwb25zZRWyAdmNiAhmb3JRUyAgZmFpbHVyZXM1k2rJAgBuQRBeVAI6jAMBAQwiYWRkZcwIYWxQVh4DAQFVG0pCAAAgmkQAACAZRk3cLXQJASl6CQEtgAkBQm4DAEMJ5QQgaYEvCG1hcEH-GcZteBhkZXRhaWxzhRcAYwkxGXAdYAQiRAkpXjMBDauhNBEBQpEADUMFaBRhaW5zIG_FrwRhbA2GBGVkDZXB-cGHLGlvbiwgc3VjaCBhcw0bTe4Uc1xuIG9yoQcodGFjayB0cmFjZSCVUwR0b0U-EGdpdmVuRR4wZXIgZW52aXJvbm1lblrhAEW-YuAApuwDDCRyZWbBbQQjL6WtRTUEcy8FXjlazbURAUIrAQBBdfel4QBzqdEcb2NjdXJyZWQZSC6iAQhJbm4SQQgynwVWjgAucQMAInK0AQVVFCBwb2ludIU1AGElUARpY7FFwSsEYWfl_jhmcm9tIGFub3RoZXIgdGkBBcGGDGlzIHXBrWmxVrsACFRhZ-56AS16OhMGfTc9JwAgRm0BAX5BNCBjYXRlZ29yaXrhQA3hGUYdRAwiVGltspEGAGZFz0EWFGRhdGUtdAE_cncBAWFhsLUFHYgN2AkIAQYFmgQidClEFtIIUQQIIm5hBaoEImIO5gk0IgogICAgfQogIF0KfQo="
)

var staticFiles = map[string]string{
	"/error.json": c_ErrorJson,
}

// Lookup returns the bytes associated with the given path, or nil if the path was not found.
func Lookup(path string) []byte {
	s, ok := staticFiles[path]
	if ok {
		d, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			log.Print("main.Lookup: ", err)
			return nil
		}
		r, err := snappy.Decode(nil, d)
		if err != nil {
			log.Print("main.Lookup: ", err)
			return nil
		}
		return r
	}
	return nil
}

// ServeHTTP serves the stored file data over HTTP.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasSuffix(p, "/") {
		p += "index.html"
	}
	b := Lookup(p)
	if b != nil {
		mt := mime.TypeByExtension(path.Ext(p))
		if mt != "" {
			w.Header().Set("Content-Type", mt)
		}
		w.Header().Set("Expires", time.Now().AddDate(0, 0, 1).Format(time.RFC1123))
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}
