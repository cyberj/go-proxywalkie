package proxy

import (
	"encoding/json"
	"fmt"
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

	serverfiles map[string]*walkie.File
	myfiles     map[string]*walkie.File

	// router *chi.Router
	m sync.RWMutex

	path string
}

func NewProxy(path string, server_url string) (proxy *Proxy, err error) {
	proxy = &Proxy{
		server_url: server_url,
		path:       path,
	}

	// proxy.m.Lock()

	walkiedir, err := walkie.NewWalkie(path)
	if err != nil {
		return
	}
	proxy.getServerDirectory()

	proxy.walkiedir = walkiedir
	proxy.walkiedir.Explore()
	// time.Sleep(10 * time.Millisecond)

	// First sync
	proxy.syncDir()
	go func() {
		for {
			proxy.syncDir()
			time.Sleep(10 * time.Second)

		}
	}()

	return
}

func (p *Proxy) Ready() {
	p.m.RLock()
	p.m.RUnlock()
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

	p.m.Lock()
	p.serverdir = *exporteddir
	p.serverfiles = p.serverdir.ListFiles()
	p.m.Unlock()
}

func (p *Proxy) syncDir() {
	p.m.RLock()
	add, del, err := p.walkiedir.SyncDir(p.serverdir)
	if err != nil {
		p.m.RUnlock()
		logrus.Error(err)
	}
	p.m.RUnlock()
	logrus.Infof("syncdir : add=%v del=%v", add, del)

}

// Retrive file from server
func (p *Proxy) getFile(fileid string) (err error) {

	srvfile, ok := p.findFileSrv(fileid)
	if !ok {
		return fmt.Errorf("Unknow file")
	}

	// logrus.Error(p.server_url, fileid)
	res, err := http.Get(p.server_url + "/" + fileid)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("Server status %v", res.StatusCode)
		return
	}

	// localpath := filepath.Join(p.path, filepath.FromSlash(fileid))
	err = p.walkiedir.UpdateOrCreateFile(fileid, res.Body, *srvfile)
	defer res.Body.Close()
	if err != nil {
		return
	}
	return

}

// Check if file exists
func (p *Proxy) checkFile(filepath string) (ok bool) {
	srvfile, ok := p.findFileSrv(filepath)
	if !ok {
		return false
	}

	myfile, ok := p.walkiedir.GetFile(filepath)
	if !ok {
		return false
	}

	return myfile.Equals(*srvfile)

}

// Check if file exists on serveur list
func (p *Proxy) findFileSrv(filepath string) (srvfile *walkie.File, ok bool) {
	p.m.RLock()
	defer p.m.RUnlock()
	srvfile, ok = p.serverfiles[filepath]
	if !ok {
		return
	}
	return
}
