package proxy

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
)

type Proxy struct {
	// chi.Router

	walkiedir *walkie.Walkie
	serverdir walkie.Directory

	server_url string
	lastping   time.Time

	// router *chi.Router
	m sync.Mutex

	path string
}

func NewProxy(path string, server_url string) (proxy *Proxy, err error) {
	proxy = &Proxy{
		server_url: server_url,
		path:       path,
	}

	proxy.m.Lock()

	walkiedir, err := walkie.NewWalkie(path)
	if err != nil {
		return
	}

	proxy.walkiedir = walkiedir
	proxy.walkiedir.Explore()

	proxy.getServerDirectory()
	go proxy.syncDir()

	return
}

func (p *Proxy) Ready() {
	p.m.Lock()
	p.m.Unlock()
}

func (p *Proxy) getServerDirectory() {
	res, err := http.Get(p.server_url)
	if err != nil {
		logrus.Errorf("Error on cache loop : %s", err)
	}

	exporteddir := &walkie.Directory{}
	err = json.NewDecoder(res.Body).Decode(exporteddir)

	if err != nil {
		logrus.Errorf("Error on cache loop decondig json : %s", err)

	}

	p.serverdir = *exporteddir
	p.m.Unlock()
}

func (p *Proxy) syncDir() {
	for {
		p.Ready()
		add, del, err := p.walkiedir.SyncDir(p.serverdir)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Infof("syncdir : add=%v del=%v", add, del)
	}

}
