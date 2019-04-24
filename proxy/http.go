package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (p *Proxy) Router() http.Handler {
	rtr := chi.NewRouter()
	// p.router = rtr

	rtr.Use(middleware.RequestID)
	rtr.Use(middleware.RealIP)
	rtr.Use(middleware.Logger)
	rtr.Use(middleware.Recoverer)

	// proxy.HandleFunc("/", proxy.handleOK)
	rtr.HandleFunc("/_cache/", p.handleStatus)
	rtr.HandleFunc("/_cache/client", p.handleServerCache)
	rtr.HandleFunc("/_cache/server", p.handleClientCache)
	rtr.HandleFunc("/*", p.handleServeFile)

	return rtr
}

func (p *Proxy) handleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Last ping : %s", time.Since(p.lastping))
	return
}

func (p *Proxy) handleClientCache(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(p.walkiedir.Directory)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error while encoding : %s", err)
	}
	return
}

func (p *Proxy) handleServerCache(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(p.serverdir)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error while encoding : %s", err)
	}
	return
}

func (p *Proxy) handleServeFile(w http.ResponseWriter, r *http.Request) {
	var err error

	logrus.Infof("%s", r.URL)
	path := chi.URLParam(r, "*")
	logrus.Infof("%s", path)

	// Don't show directories
	if strings.HasSuffix(path, "/") {
		w.WriteHeader(404)
		fmt.Fprint(w, "Not found")

	}

	// File not found on server
	_, ok := p.findFileSrv(path)
	if !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "Not found")
	}

	if r.Method != http.MethodGet {
		return
	}

	// File don't exist or is invalid
	if !p.checkFile(path) {
		err = p.getFile(path)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Not found")
		}
	}

	http.ServeFile(w, r, filepath.Join(p.path, path))

	return
}
