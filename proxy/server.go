package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Proxy struct {
	chi.Router

	walkiedir *walkie.Walkie
}

func NewProxy(path string) (proxy *Proxy, err error) {
	proxy = &Proxy{Router: chi.NewRouter()}

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

	return
}

func (p *Proxy) handleServeFile(w http.ResponseWriter, r *http.Request) {

	logrus.Errorf("%s", r.URL)

	if r.URL.Path == "/" {
		err := json.NewEncoder(w).Encode(p.walkiedir.Directory)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error while encoding : %s", err)
		}
		return
	}

	// http.ServeFile(w, r, )

	return
}
