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

	data := struct {
		LastPing  time.Duration `json:"name"`
		Running   bool          `json:"running"`
		ServerURL string        `json:"server_url"`
		LocalPath string        `json:"local_path"`
	}{
		time.Since(p.lastping),
		p.Running(),
		p.server_url,
		p.path,
	}

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	fmt.Fprint(w, json)

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

	path := chi.URLParam(r, "*")

	if r.Method != http.MethodGet {
		return
	}

	// Don't show directories
	if strings.HasSuffix(path, "/") {
		w.WriteHeader(404)
		fmt.Fprint(w, "Not found")
		return
	}

	// Special mode where we were unable to reach our server
	nulltime := time.Time{}
	if p.lastping == nulltime {
		if p.checkLocalFile(path) {
			w.Header().Add("X-ProxyWalkie-Cached", "true")
			// The file is in cache but we are not sure abobut it being correct
			// because we are unable to reach the server
			w.Header().Add("X-ProxyWalkie-Local", "true")
			http.ServeFile(w, r, filepath.Join(p.path, path))
			return
		}
	}

	// File not found on server
	_, ok := p.findFileSrv(path)
	if !ok {

		w.WriteHeader(404)
		fmt.Fprint(w, "Not found")
		return

	}

	// File don't exist or is invalid
	if !p.checkFile(path) {
		logrus.Infof("File %s not found or invalid", path)
		err = p.getFile(path)

		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Not found on server")
		}
	} else {
		logrus.Debugf("File %s hit cache", path)
		w.Header().Add("X-ProxyWalkie-Cached", "true")
	}

	http.ServeFile(w, r, filepath.Join(p.path, path))

	return
}
