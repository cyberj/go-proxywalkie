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

	// Interval between 2 server queries
	SyncInterval time.Duration

	// Clean files
	Clean bool

	server_url string
	lastping   time.Time

	serverfiles map[string]*walkie.File

	// router *chi.Router
	m sync.RWMutex

	done    chan bool
	running bool

	path string
}

func NewProxy(path string, server_url string) (proxy *Proxy, err error) {
	proxy = &Proxy{
		server_url:   server_url,
		path:         path,
		serverdir:    walkie.Directory{},
		SyncInterval: 1 * time.Minute,
		done:         make(chan bool),
	}

	// proxy.m.Lock()

	walkiedir, err := walkie.NewWalkie(path)
	if err != nil {
		return
	}

	proxy.walkiedir = walkiedir
	proxy.walkiedir.Explore()
	// time.Sleep(10 * time.Millisecond)

	proxy.Run()

	return
}

func (p *Proxy) Ready() bool {
	p.m.RLock()
	defer p.m.RUnlock()
	return !p.serverdir.DeepEquals(walkie.Directory{}) // eww
}

// Start server loop
func (p *Proxy) Run() (err error) {
	if p.running {
		return fmt.Errorf("Already running")
	}
	p.done = make(chan bool)

	// First tick free !
	err = p.getServerDirectory()
	if err != nil {
		logrus.Error(err)
		return
	}

	p.running = true

	go func() {
		ticker := time.NewTicker(p.SyncInterval)

		for {
			select {
			case <-ticker.C:
				p.getServerDirectory()
			case <-p.done:
				ticker.Stop()
				p.running = false
				return
			}

		}
	}()
	return
}

func (p *Proxy) Stop() {
	if p.running {
		close(p.done)

		for p.running {

		}
	}
}

func (p *Proxy) getServerDirectory() (err error) {

	res, err := http.Get(p.server_url)
	if err != nil {
		err = fmt.Errorf("Error on cache loop : %s", err)
		return
	}

	exporteddir := walkie.Directory{}
	err = json.NewDecoder(res.Body).Decode(&exporteddir)
	if err != nil {
		err = fmt.Errorf("Error on cache loop decondig json : %s", err)
		return
	}

	// Change only when needed or first run
	p.m.RLock()
	same := p.serverdir.DeepEquals(exporteddir)
	p.m.RUnlock()
	if p.running && same {
		return
	}

	p.m.Lock()
	p.serverdir = exporteddir
	p.serverfiles = p.serverdir.ListFiles()
	p.m.Unlock()
	p.syncDir()
	if p.Clean {
		p.cleanFiles()
	}
	return
}

// syncDir Adds and remove directories accordingly with server
func (p *Proxy) syncDir() {
	p.m.RLock()
	defer p.m.RUnlock()
	add, del, err := p.walkiedir.SyncDir(p.serverdir)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("syncdir : add=%v del=%v", add, del)

}

// cleanFiles swipe uneeded files
func (p *Proxy) cleanFiles() {
	p.m.RLock()
	deleted, err := p.walkiedir.CleanFiles(p.serverdir)
	if err != nil {
		logrus.Errorf("cleanFiles error %s", err)
		return
	}
	logrus.Infof("cleanFiles : deleted=%v", deleted)
	p.m.RUnlock()

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
