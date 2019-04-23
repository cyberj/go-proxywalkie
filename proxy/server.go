package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Proxy struct {
	chi.Router

	walkiedir *walkie.Walkie
	path      string
}

func NewProxy(path string) (proxy *Proxy, err error) {
	proxy = &Proxy{Router: chi.NewRouter(), path: path}

	walkiedir, err := walkie.NewWalkie(path)
	if err != nil {
		return
	}

	proxy.walkiedir = walkiedir
	proxy.walkiedir.Explore()

	proxy.Use(middleware.RequestID)
	proxy.Use(middleware.RealIP)
	proxy.Use(middleware.Logger)
	proxy.Use(middleware.Recoverer)

	proxy.HandleFunc("/", proxy.handleServeFile)
	proxy.HandleFunc("/*", proxy.handleServeFile)

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
