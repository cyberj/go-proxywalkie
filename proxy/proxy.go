package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	// Background sync
	Sync bool

	server_url string
	lastping   time.Time

	serverfiles map[string]*walkie.File

	// router *chi.Router
	m sync.RWMutex

	done chan bool

	// Chan to stop file checking
	stopCh chan bool
	// chan to transmit file to fetch
	fetchCh chan string
	// runningCh chan bool

	// Local directory
	path string
}

// NewProxy is the default call to create a proxy
func NewProxy(path string, server_url string) (proxy *Proxy, err error) {
	return NewProxyParams(path, server_url, 1*time.Minute, false, false)
}

// NewProxyParams is the customizable function
func NewProxyParams(path string, server_url string, interval time.Duration, clean bool, sync bool) (proxy *Proxy, err error) {
	proxy = &Proxy{
		server_url:   server_url,
		path:         path,
		serverdir:    walkie.Directory{},
		SyncInterval: interval,
		stopCh:       make(chan bool),
		fetchCh:      make(chan string),
		// runningCh:    make(chan bool),

		Clean: clean,
		Sync:  sync,
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

func (p *Proxy) Running() bool {

	if p.done == nil {
		return false
	}

	select {
	case <-p.done:
		return false
	default:
		return true
	}

}

// Start server loop
func (p *Proxy) Run() (err error) {

	if p.Running() {
		return fmt.Errorf("Already running")
	}

	// Reset done State

	// First tick free !
	err = p.getServerDirectory()
	if err != nil {
		logrus.Error(err)
	}
	doneCh := make(chan bool)

	go func(done chan bool) {
		ticker := time.NewTicker(p.SyncInterval)
		// Ping timer
		ticker2 := time.NewTicker(5 * time.Minute)

		for {
			select {
			case <-ticker.C:
				p.getServerDirectory()
			case <-ticker2.C:
				p.pingServer()
			case fileid := <-p.fetchCh:
				if fileid == "" {
					logrus.Fatal("Main loop recieved empty fileid = should not append")
				}
				p.getFile(fileid)
			case <-done:
				ticker.Stop()
				ticker2.Stop()

				return
			}

		}
	}(doneCh)
	p.done = doneCh

	return
}

func (p *Proxy) pingServer() {

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// res, err := http.Get(p.server_url + "/ping")
	req, err := http.NewRequest("GET", p.server_url+"/ping", nil)
	if err != nil {
		err = fmt.Errorf("Error on ping loop : %s", err)
		return
	}

	// Get Hostname
	hostname, err := os.Hostname()
	if err == nil {
		req.Header.Add("X-ProxyWalkie-Hostname", hostname)
	}

	_, err = client.Do(req)
	if err != nil {
		err = fmt.Errorf("Error on ping loop : %s", err)
		return
	}

	return
}

func (p *Proxy) Stop() {

	if p.Running() {
		close(p.done)

		for p.Running() {

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
	if same && p.Running() {
		logrus.Info("getServerDirectory : No changes detected")
		return
	}

	p.m.Lock()
	p.serverdir = exporteddir
	p.serverfiles = p.serverdir.ListFiles()
	p.m.Unlock()
	p.syncDir()

	// Clean files if needed
	if p.Clean {
		p.cleanFiles()
	}
	// Sync files in background
	if p.Sync {
		p.syncFiles(p.done)
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

// syncFiles Lists all files to fetch and add them in event loop
func (p *Proxy) syncFiles(done chan bool) {

	p.m.RLock()
	toadd, _ := p.walkiedir.Directory.DiffFiles(p.serverdir)
	p.m.RUnlock()

	close(p.stopCh)
	// <-p.stopCh // Sync closing

	stopCh := make(chan bool)
	p.stopCh = stopCh

	go func(add []string) {
		for _, v := range add {

			// already exists
			if p.checkFile(v) {
				continue
			}

			select {
			case <-done:
				// close(p.tofetch)
				return
			case <-stopCh:
				return
			default:
			}
			logrus.Debugf("syncFiles : Add new file to check : %s", v)
			p.fetchCh <- v
			logrus.Debugf("syncFiles : Done added new file to check : %s", v)

		}
	}(toadd)

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
		return fmt.Errorf("Unknown file")
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

// Check if local cache of that file is valid
func (p *Proxy) checkFile(filepath string) (ok bool) {
	srvfile, ok := p.findFileSrv(filepath)
	if !ok {
		return false
	}

	myfile, ok := p.walkiedir.GetFile(filepath)
	if !ok {
		return false
	}

	err := myfile.Compare(*srvfile)
	if err != nil {
		logrus.Infof("File '%s' raised a FileCompareError : %s", filepath, err.Error())
		return false
	}

	return true
}

// Check if local cache of that file exists
func (p *Proxy) checkLocalFile(filepath string) (ok bool) {

	_, ok = p.walkiedir.GetFile(filepath)
	if !ok {
		return false
	}

	return true
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
