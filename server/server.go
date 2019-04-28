package proxy

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	chi.Router

	dircache []byte

	walkiedir *walkie.Walkie
	path      string
}

func NewServer(path string) (server *Server, err error) {
	server = &Server{Router: chi.NewRouter(), path: path}

	walkiedir, err := walkie.NewWalkie(path)
	if err != nil {
		return
	}

	server.walkiedir = walkiedir
	server.walkiedir.Explore()

	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)
	// server.Use(middleware.DefaultCompress)

	server.HandleFunc("/", server.handleServeFile)
	server.HandleFunc("/*", server.handleServeFile)

	// Cache first
	server.cache()

	go func(server *Server) {
		for {
			time.Sleep(10 * time.Second)
			server.cache()
		}
	}(server)

	return
}

func (p *Server) cache() {

	buf := &bytes.Buffer{}
	// copy buffer
	buf2 := &bytes.Buffer{}

	gzip_encoder := gzip.NewWriter(buf)
	w := io.MultiWriter(buf2, gzip_encoder)

	json_encoder := json.NewEncoder(w)

	err := json_encoder.Encode(p.walkiedir.Directory)
	gzip_encoder.Flush()
	gzip_encoder.Close()
	if err == nil {
		p.dircache = buf.Bytes()
		logrus.Debugf("Cache 'dircache' : before=%v after=%v", len(buf2.Bytes()), len(p.dircache))
	}

}

func (p *Server) handleServeFile(w http.ResponseWriter, r *http.Request) {

	logrus.Debugf("%s", r.URL)
	path := chi.URLParam(r, "*")
	logrus.Debugf("%s", path)

	if r.URL.Path == "/" {
		logrus.Info("Using cache")

		// fallback
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gzip_reader, err := gzip.NewReader(bytes.NewReader(p.dircache))
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "Error when initializing gzip decoder : %s", err)
				return
			}
			io.Copy(w, gzip_reader)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")

		io.Copy(w, bytes.NewReader(p.dircache))

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

}
