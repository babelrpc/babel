package main

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os/exec"
)

func handleSwagger(router *httprouter.Router) error {
	// handle RESTful version of API
	args := []string{
		"-flat",
		"-rest",
		"-error",
		"-basepath", conf.RestPath,
		"-host", conf.PubAddr,
		"-title", conf.Title,
		"-version", conf.RestVersion,
	}
	args = append(args, conf.Args...)

	cmd := exec.Command("babel2swagger", args...)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err := cmd.Run()
	if err != nil {
		log.Print("Cannot generate swagger: ", err)
		log.Print("Result: ", serr.String())
	} else {
		b := sout.Bytes()
		router.HandlerFunc("GET", "/rest.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})
	}

	return nil
}

func handleBabel(router *httprouter.Router) error {
	// handle Babel version of API
	args := []string{
		"-flat",
		"-basepath", conf.BabelPath,
		"-host", conf.PubAddr,
		"-title", conf.Title,
		"-version", conf.BabelVersion,
	}
	args = append(args, conf.Args...)

	cmd := exec.Command("babel2swagger", args...)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err := cmd.Run()
	if err != nil {
		log.Print("Cannot generate babel swagger: ", err)
		log.Print("Result: ", serr.String())
	} else {
		b := sout.Bytes()
		router.HandlerFunc("GET", "/babel.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})
	}

	return nil
}
