package main

import (
	"fmt"
	"github.com/ancientlore/kubismus"
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// Service implementation for log tail
type svcImpl struct {
	cancel context.CancelFunc
}

// Start starts the service - it shoud not block
func (i *svcImpl) Start(s service.Service) error {
	// start worker
	go i.doWork()

	return nil
}

// Stop stops the service
func (i *svcImpl) Stop(s service.Service) error {
	// cancel work
	if i.cancel != nil {
		i.cancel()
	}
	return nil
}

func switchPrefix(oldPrefix, newPrefix string, h http.Handler) http.Handler {
	if oldPrefix == newPrefix {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, oldPrefix); len(p) < len(r.URL.Path) {
			r.URL.Path = newPrefix + p
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func addSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

// doWork does the work of the service
func (i *svcImpl) doWork() {
	// setup kubismus status
	kubismus.Setup("babelproxy", addSlash(conf.MediaPath)+"/logo36.png")
	kubismus.Define("Requests", kubismus.COUNT, "HTTP Requests")
	kubismus.Note("Processors", fmt.Sprintf("%d of %d", runtime.GOMAXPROCS(0), runtime.NumCPU()))
	kubismus.Note("Version", BABELPROXY_VERSION)
	kubismus.Note("REST Service Proxy", "http://"+conf.PubAddr+addSlash(conf.RestPath))
	kubismus.Note("Babel Service Proxy", "http://"+conf.PubAddr+addSlash(conf.BabelPath))
	kubismus.Note("Destination Babel Service", conf.BabelProto+"://"+conf.BabelAddr+addSlash(conf.BabelPath))
	kubismus.Note("Swagger UI (REST)", "http://"+conf.PubAddr+addSlash(conf.SwaggerPath)+"?url=/rest.json")
	kubismus.Note("Swagger UI (Babel)", "http://"+conf.PubAddr+addSlash(conf.SwaggerPath)+"?url=/babel.json")

	mux := http.NewServeMux()
	mux.Handle(addSlash(conf.StatusPath), switchPrefix(addSlash(conf.StatusPath), "/", http.HandlerFunc(kubismus.ServeHTTP)))
	mux.Handle(addSlash(conf.MediaPath), switchPrefix(addSlash(conf.MediaPath), "/media/", http.HandlerFunc(ServeHTTP)))

	// initialize HTTP settings for posting
	i.cancel = initHttp()

	// spawn a function that updates the number of goroutines shown in the status page
	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				kubismus.Note("Goroutines", fmt.Sprintf("%d", runtime.NumGoroutine()))
			}
		}
	}()

	// status site
	go http.ListenAndServe(conf.StatusAddr, mux)

	// read files
	midl, err := loadBabelFiles(conf.Args)
	if err != nil {
		log.Fatal(err)
	}

	// handle swagger UI and JSON generator
	router := httprouter.New()
	router.Handler("GET", addSlash(conf.SwaggerPath)+"*path", switchPrefix(addSlash(conf.SwaggerPath), "/swagger/", http.HandlerFunc(ServeHTTP)))
	err = handleSwagger(router)
	if err != nil {
		log.Fatal(err)
	}
	err = handleBabel(router)
	if err != nil {
		log.Fatal(err)
	}

	// add handlers
	err = addHandlers(router, midl)
	if err != nil {
		log.Fatal(err)
	}

	// listen up!
	log.Fatal(http.ListenAndServe(conf.RestAddr, router))
}
