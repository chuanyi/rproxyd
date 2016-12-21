package main

import (
	"bufio"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/kardianos/osext"
	"github.com/kardianos/service"
)

type RProxy struct {
	Pre   string
	Proxy *httputil.ReverseProxy
}

type RPServer struct {
	Addr string
	RPs  []RProxy
}

func (s *RPServer) ReadCfg(file string) {
	cfg, _ := os.Open(file)
	defer cfg.Close()
	scanner := bufio.NewScanner(cfg)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "listen") {
			s.Addr = strings.Split(line, " ")[1]
		} else if strings.HasPrefix(line, "proxy") {
			tmp := strings.Split(line, " ")
			remote, _ := url.Parse(tmp[2])
			proxy := httputil.NewSingleHostReverseProxy(remote)
			s.RPs = append(s.RPs, RProxy{Pre: tmp[1], Proxy: proxy})
		}
	}
}

func (s *RPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handled := false
	for _, rp := range s.RPs {
		if strings.HasPrefix(r.RequestURI, rp.Pre) {
			handled = true
			rp.Proxy.ServeHTTP(w, r)
			break
		}
	}
	if !handled {
		w.WriteHeader(404)
		w.Write([]byte("Bad Proxy\n"))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	s := &RPServer{}
	path, _ := osext.ExecutableFolder()
	s.ReadCfg(path + "\\rproxyd")
	http.ListenAndServe(s.Addr, s)
}
