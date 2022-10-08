package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/ginuerzh/gost"
	"github.com/ginuerzh/gost/zbutil"
	"github.com/ginuerzh/gost/zbutil/loglv"
)

var (
	routers []router
)

type baseConfig struct {
	route
	Routes []route
	Debug  bool
}

func parseBaseConfig(s string) (*baseConfig, error) {
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(baseCfg); err != nil {
		return nil, err
	}

	return baseCfg, nil
}

var (
	defaultCertFile = "cert.pem"
	defaultKeyFile  = "key.pem"
)

// Load the certificate from cert & key files and optional client CA file,
// will use the default certificate if the provided info are invalid.
func tlsConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	if certFile == "" || keyFile == "" {
		certFile, keyFile = defaultCertFile, defaultKeyFile
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	if pool, _ := loadCA(caFile); pool != nil {
		cfg.ClientCAs = pool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return cfg, nil
}

func loadCA(caFile string) (cp *x509.CertPool, err error) {
	if caFile == "" {
		return
	}
	cp = x509.NewCertPool()
	data, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	if !cp.AppendCertsFromPEM(data) {
		return nil, errors.New("AppendCertsFromPEM failed")
	}
	return
}

func parseKCPConfig(configFile string) (*gost.KCPConfig, error) {
	if configFile == "" {
		return nil, nil
	}
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &gost.KCPConfig{}
	if err = json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func parseUsers(authFile string) (users []*url.Userinfo, err error) {
	if authFile == "" {
		return
	}

	file, err := os.Open(authFile)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		s := strings.SplitN(line, " ", 2)
		if len(s) == 1 {
			users = append(users, url.User(strings.TrimSpace(s[0])))
		} else if len(s) == 2 {
			users = append(users, url.UserPassword(strings.TrimSpace(s[0]), strings.TrimSpace(s[1])))
		}
	}

	err = scanner.Err()
	return
}

func parseAuthenticator(s string) (gost.Authenticator, error) {
	if s == "" {
		return nil, nil
	}
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	au := gost.NewLocalAuthenticator(nil)
	au.Reload(f, false)

	go gost.PeriodReload(au, s)

	return au, nil
}

func parseIP(s string, port string) (ips []string) {
	if s == "" {
		return
	}
	if port == "" {
		port = "8080" // default port
	}

	file, err := os.Open(s)
	if err != nil {
		ss := strings.Split(s, ",")
		for _, s := range ss {
			s = strings.TrimSpace(s)
			if s != "" {
				// TODO: support IPv6
				if !strings.Contains(s, ":") {
					s = s + ":" + port
				}
				ips = append(ips, s)
			}

		}
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, ":") {
			line = line + ":" + port
		}
		ips = append(ips, line)
	}
	return
}

func parseBypass(s string, inwall, chkwall, white, fakeip, ischain bool) *gost.Bypass {
	if s == "" {
		if chkwall {
			return gost.NewBypass(inwall, chkwall, white, fakeip, ischain)
		}
		loglv.Inf.Stackf("[1]parseBypass(s:%s,chkwall:%v)\n", s, chkwall)
		return nil
	}
	var matchers []gost.Matcher
	if strings.HasPrefix(s, "~") { //白名单反成黑名单
		s = strings.TrimLeft(s, "~")
		white = false
	}

	f, err := os.Open(s)
	if err != nil {
		for _, s := range strings.Split(s, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			matchers = append(matchers, gost.NewMatcher(s))
		}
		return gost.NewBypass(inwall, chkwall, white, fakeip, ischain, matchers...)
	}
	defer f.Close()

	bp := gost.NewBypass(inwall, chkwall, white, fakeip, ischain)
	bp.Reload(f, false)
	go gost.PeriodReload(bp, s)

	return bp
}

func parseResolver(cfg string) gost.Resolver {
	if cfg == "" {
		return nil
	}
	var nss []gost.NameServer

	f, err := os.Open(cfg)
	if err != nil {
		for _, s := range strings.Split(cfg, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if strings.HasPrefix(s, "https") {
				p := "https"
				u, _ := url.Parse(s)
				if u == nil || u.Scheme == "" {
					continue
				}
				if u.Scheme == "https-chain" {
					p = u.Scheme
				}
				ns := gost.NameServer{
					Addr:     s,
					Protocol: p,
				}
				nss = append(nss, ns)
				continue
			}

			ss := strings.Split(s, "/")
			if len(ss) == 1 {
				ns := gost.NameServer{
					Addr: ss[0],
				}
				nss = append(nss, ns)
			}
			if len(ss) == 2 {
				ns := gost.NameServer{
					Addr:     ss[0],
					Protocol: ss[1],
				}
				nss = append(nss, ns)
			}
		}
		return gost.NewResolver(0, nss...)
	}
	defer f.Close()

	resolver := gost.NewResolver(0)
	resolver.Reload(f, false)

	go gost.PeriodReload(resolver, cfg)

	return resolver
}

func parseHosts(s string) *gost.Hosts {
	f, err := os.Open(s)
	if err != nil {
		return nil
	}
	defer f.Close()

	hosts := gost.NewHosts()
	hosts.Reload(f, false)

	go gost.PeriodReload(hosts, s)

	return hosts
}

func parseIPRoutes(s string) (routes []gost.IPRoute) {
	if s == "" {
		return
	}

	file, err := os.Open(s)
	if err != nil {
		ss := strings.Split(s, ",")
		for _, s := range ss {
			if _, inet, _ := net.ParseCIDR(strings.TrimSpace(s)); inet != nil {
				routes = append(routes, gost.IPRoute{Dest: inet})
			}
		}
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), "\t", " ", -1)
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var route gost.IPRoute
		var ss []string
		for _, s := range strings.Split(line, " ") {
			if s = strings.TrimSpace(s); s != "" {
				ss = append(ss, s)
			}
		}
		if len(ss) > 0 && ss[0] != "" {
			_, route.Dest, _ = net.ParseCIDR(strings.TrimSpace(ss[0]))
			if route.Dest == nil {
				continue
			}
		}
		if len(ss) > 1 && ss[1] != "" {
			route.Gateway = net.ParseIP(ss[1])
		}
		routes = append(routes, route)
	}
	return routes
}
func parseLF(s string) (L, F stringList) {
	if s == "" {
		return
	}

	file, err := os.Open(s)
	if err != nil {
		return
	}

	defer file.Close()
	var l_err error = nil

	for rdr, lineno := bufio.NewReader(file), 0; l_err != io.EOF; lineno++ {
		var l_byts []byte
		l_byts, l_err = rdr.ReadBytes('\n')
		if l_err != nil && l_err != io.EOF {
			fmt.Printf("文件末尾的格式错误.\n")
			return nil, nil
		} else if len(l_byts) < 1 {
			continue
		}
		line := strings.Replace(string(l_byts), "\t", " ", -1)
		line = strings.TrimSpace(line)
		fmt.Printf("line[%d]=%s\n", lineno, line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.LastIndex(line, "#"); idx > 0 && line[idx-1] == ' ' {
			line = line[:idx-1]
		}

		var ss []string
		for _, s := range zbutil.SplitFields(line, "\t ", false) {
			if s = strings.TrimSpace(s); s != "" {
				ss = append(ss, s)
			}
		}
		if len(ss) < 2 {
			continue
		}
		if ss[0] == "-L" {
			L.Set(ss[1])
		} else if ss[0] == "-F" {
			F.Set(ss[1])
		}
	}
	return
}
