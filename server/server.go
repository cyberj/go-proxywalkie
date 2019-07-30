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
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ClientPingItemJson struct {
	Hostname string `json:"hostname"`
	LastPing int64  `json:"lastping_seconds"`
}

type ClientPingJson struct {
	List []ClientPingItemJson `json:"list"`
}

type Server struct {
	chi.Router

	dircache []byte

	// Interval between 2 server queries
	SyncInterval time.Duration

	walkiedir *walkie.Walkie
	path      string

	m sync.RWMutex

	// Watchdog for clients with his mutex
	clientpinglist map[string]time.Time
	clientMu       sync.RWMutex
}

// NewProxy is the default call to create a proxy
func NewServer(path string) (server *Server, err error) {
	return NewServerParams(path, 10*time.Minute)
}

func NewServerParams(path string, interval time.Duration) (server *Server, err error) {
	logrus.Debug("Creating a new server")

	server = &Server{Router: chi.NewRouter(), path: path, SyncInterval: interval, clientpinglist: map[string]time.Time{}}

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

	server.Get("/_files", server.handleFileList)
	server.Get("/ping", server.handlePing)
	server.Get("/_clients", server.handleClientList)
	server.HandleFunc("/$", server.handleServeFile)
	server.HandleFunc("/*", server.handleServeFile)

	// Cache first
	server.cache()

	go func(server *Server) {
		ticker := time.NewTicker(server.SyncInterval)

		for {
			select {
			case <-ticker.C:
				server.walkiedir.Explore()
				server.cache()
			}
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
		p.m.Lock()
		p.dircache = buf.Bytes()
		logrus.Debugf("Gzip cache 'dircache' : before=%vb after=%vb", len(buf2.Bytes()), len(p.dircache))
		p.m.Unlock()
	}

}

func (p *Server) handleServeFile(w http.ResponseWriter, r *http.Request) {

	path := chi.URLParam(r, "*")
	logrus.Debugf("URL='%s', Path='%s'", r.URL, path)

	if r.URL.Path == "/" {
		logrus.Info("Using cache")
		p.m.RLock()
		defer p.m.RUnlock()

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

func (p *Server) handleFileList(w http.ResponseWriter, r *http.Request) {

	for _, filepath := range p.walkiedir.ListFiles() {
		fmt.Fprintln(w, filepath)
	}

}

// handlePing registers someone in cache list
func (p *Server) handlePing(w http.ResponseWriter, r *http.Request) {

	hostname := r.Header.Get("X-ProxyWalkie-Hostname")
	if hostname == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing ProxyWalkie Header")
		return
	}

	p.clientMu.Lock()
	p.clientpinglist[hostname] = time.Now()
	p.clientMu.Unlock()

	return

}

// handleClientList is a Convenient list to check if someone didn't respond for a long time...
func (p *Server) handleClientList(w http.ResponseWriter, r *http.Request) {

	data := ClientPingJson{}

	p.clientMu.RLock()
	for key, val := range p.clientpinglist {
		data.List = append(data.List, ClientPingItemJson{Hostname: key, LastPing: int64(time.Since(val).Seconds())})
	}
	p.clientMu.RUnlock()

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "%s", json)

}
