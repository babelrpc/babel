package main

import (
	"github.com/golang/snappy"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

const (
	c_ErrorJson = "-RRQewogICJzd2FnZ2VyIjogIjIuMCIsARQYaW5mbyI6IAEgHCAgInRpdGxlASA0TXkgQXBwbGljYXRpb24FKyQgICJkZXNjcmlwBRN0OiAiLi4vLi4vYmFiZWx0ZW1wbGF0ZXMvZXJyb3IuBRUIOlxuCQIEIEIBJAwncyBFAR0YIEZvcm1hdAEbWCBDb3B5cmlnaHQgKEMpIDIwMTUgVGhlCTAgIEF1dGhvcnMuASpcR2VuZXJhdGVkIG9uIFRodSwgMjIgT2N0CTQ8MDA6MTM6NTMgRURUIHZpYQl6AHQFhgAyDfcEIFyCsQAEXCIR6gx2ZXJzAfkoOiAiMS4wIgogIH0lLwxob3N0ISAQbG9jYWwFDQUXHGJhc2VQYXRoARsALwVETCJjb25zdW1lcyI6IFsKICAgICJhOVEQL2pzb24BVgBdBT8UcHJvZHVjmiwACGF0aAFVAHsJkRRkZWZpbmkhkgUVCWolZzXWIdgEeXAl1xRvYmplY3QFoQWTOtEBLaAFVCxlcyBhIHNpbmdsZSAl1FwgbWVzc2FnZSBhbmQgY29kZSB0aGF0IG0lvghiZSAlExRpemVkXG4FI1RkaXNwbGF5ZWQgdG8gYSBjYWxsZXIuGXwccHJvcGVydGkF90mABZQIQ29kAawZEh3AFHN0cmluZxVEBTA-xABBQDxzZXJ2aWNlLXNwZWNpZmljDcIBtiF2ATgt0gkLBCJNCdw9NAEZBCJ0LTieeAAkdGV4dCBvZiB0aDE1MGluIFVTLUVuZ2xpc2hWfAAQUGFyYW05vgF5HfMQYXJyYXk28gAIaXRlQjAAHTINrUULCU1FKwkLRsYDZXEIbGlzBbkAcAGMEGV0ZXJzIbYZxy38YC4gVGhpcyBjb3VsZCBiZSB1c2VkIGJ5XG4x_oUvDCBzeXMBowFDAGdtuFE-BHMgYSxlxhleIZgALj0dDbsAfQUIDb0EIlMp1OLcAi5AAAAgUeMZhCByZXNwb25zZSARsgHZjYkMZm9yIE1TICBmYWlsdXJlczWTACJmyQIAbkEQXlQCOowDIZ8MImFkZGXMCGFsUFYeAwUmUZZKQgCBZ45EAAQgID7cAiVsAV8BBCl6AQoyJwNKPQIEQ28F5QQgaYEvCG1hcEH-GcZteBhkZXRhaWxzhRcAYwkxPe6hAwlrBCJECSleMwFN2hkxACBClAYAIA1DACABaBRhaW5zIG_FsAxhbCBkBV0EZWQNlcH6wYgsaW9uLCBzdWNoIGFzDRtN7hRzXG4gb3KhByh0YWNrIHRyYWNlIJVTBHRvZecQZ2l2ZW5FHjRlciBlbnZpcm9ubWVudFbhAEW-_uwDCuwDDCRyZWbBbQQjL0XKRTUEcy8FXhl5frwBBEEgcfel4QBzqdEcb2NjdXJyZWQZSG2CIaQQIklubmUOQggynwVWjgAucQMAInK0AQVVFCBwb2ludIV4AGElUARpY7FFYWIEYWfl_jhmcm9tIGFub3RoZXIgdGkBBcGGDGlzIHXBrWmxVrsACFRhZ_56AQp6AXF9AHPJi303PScAIEaYAgF-oREgY2F0ZWdvcml64UAN4VqwAARpbbIJBwBmElcJQRYUZGF0ZS10AT9ydwEMVGltZWGwtQUdiA3YAH0FCAEGBd4EInQpRBbSCFEECCJuYQWqBCJiDucJCUUBMwxdCn0K"
)

var staticFiles = map[string]string{
	"/error.json": c_ErrorJson,
}

func Lookup(path string) []byte {
	s, ok := staticFiles[path]
	if !ok {
		return nil
	} else {
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
}

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
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}
