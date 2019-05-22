go-proxywalkie
==============

Explore and share a filepath.

Was asked for a `proxy` cache for a proprietary app on windows. This cache must be able to serve files from another (linux) server, cache them, invalidate the cache when needed and be able to launch a background sync if needed. All in HTTP(s).

So the linux `server` watch a directory structure, and serve it after some SHA256 sums.

Much of groking, but fun project to introduce new developpers to Golang.


## Configure server

```
Serve a path for improved proxy (server)

Usage:
  proxywalkie serve [:8080] [flags]

Flags:
  -h, --help   help for serve

Global Flags:
  -C, --directory string   proxywalkie's workdir (default current directory) (default ".")
```

Example :

`proxywalkie serve -C /path/to/directory/to/serve :8080`

You can check if proxy is running at http://127.0.0.1:8080/ or http://YO.UR.SRV.IP:8080/

## Configure proxy

```
Usage:
  proxywalkie proxy SERVER [flags]

Flags:
  -b, --background          Background Sync
  -d, --delete              Delete files
  -h, --help                help for proxy
  -p, --port string         Local server URL (default "8081")
  -u, --sync-interval int   Sync interval (in minutes) (default 5)

Global Flags:
  -C, --directory string   proxywalkie's workdir (default current directory) (default ".")

```

Example :

`proxywalkie proxy -C /path/to/directory/to/sync --delete -b "http://IP.SE.RV.ER:8080/`

The proxy will be ready at http://127.0.0.1:8081/ or http://YO.UR.SRV.IP:8080/
