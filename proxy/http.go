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

	logrus.Infof("%s", r.URL)
	path := chi.URLParam(r, "*")
	logrus.Infof("%s", path)

	if r.URL.Path == "/" {
		err := json.NewEncoder(w).Encode(p.walkiedir.Directory)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error while encoding : %s", err)
		}
		return
	}

	if strings.HasSuffix(path, "/") {
		dir, ok := p.walkiedir.GetDir(path)
		if !ok {
			w.WriteHeader(404)
			fmt.Fprint(w, "Not found")
		}
		err := json.NewEncoder(w).Encode(dir)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error while encoding : %s", err)
		}
		return

	} else {
		file, ok := p.walkiedir.GetFile(path)
		if !ok {
			w.WriteHeader(404)
			fmt.Fprint(w, "Not found")
		}
		w.Header().Add("X-ProxyWalkie-Hash", file.SHA256)
		w.Header().Add("X-ProxyWalkie-Size", fmt.Sprint(file.Size))
		w.Header().Add("X-ProxyWalkie-Mtime", file.Mtime.String())

		if r.Method == http.MethodGet {
			http.ServeFile(w, r, filepath.Join(p.path, path))
		}
		return

	}

	w.WriteHeader(404)
	fmt.Fprint(w, "Not found")

	// http.ServeFile(w, r, )

	return
}
