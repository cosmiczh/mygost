package gost

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ginuerzh/gost/zbutil/loglv"
	glob "github.com/gobwas/glob"
)

// Matcher is a generic pattern matcher,
// it gives the match result of the given pattern for specific v.
type Matcher interface {
	Match(v string) bool
	String() string
}

// NewMatcher creates a Matcher for the given pattern.
// The acutal Matcher depends on the pattern:
// IP Matcher if pattern is a valid IP address.
// CIDR Matcher if pattern is a valid CIDR address.
// Domain Matcher if both of the above are not.
func NewMatcher(pattern string) Matcher {
	if pattern == "" {
		return nil
	}
	if ip := net.ParseIP(pattern); ip != nil {
		return IPMatcher(ip)
	}
	if _, inet, err := net.ParseCIDR(pattern); err == nil {
		return CIDRMatcher(inet)
	}
	return DomainMatcher(pattern)
}

type ipMatcher struct {
	ip net.IP
}

// IPMatcher creates a Matcher for a specific IP address.
func IPMatcher(ip net.IP) Matcher {
	return &ipMatcher{
		ip: ip,
	}
}

func (m *ipMatcher) Match(ip string) bool {
	if m == nil {
		return false
	}
	return m.ip.Equal(net.ParseIP(ip))
}

func (m *ipMatcher) String() string {
	return "ip " + m.ip.String()
}

type cidrMatcher struct {
	ipNet *net.IPNet
}

// CIDRMatcher creates a Matcher for a specific CIDR notation IP address.
func CIDRMatcher(inet *net.IPNet) Matcher {
	return &cidrMatcher{
		ipNet: inet,
	}
}

func (m *cidrMatcher) Match(ip string) bool {
	if m == nil || m.ipNet == nil {
		return false
	}
	return m.ipNet.Contains(net.ParseIP(ip))
}

func (m *cidrMatcher) String() string {
	return "cidr " + m.ipNet.String()
}

type domainMatcher struct {
	pattern string
	glob    glob.Glob
}

// DomainMatcher creates a Matcher for a specific domain pattern,
// the pattern can be a plain domain such as 'example.com',
// a wildcard such as '*.exmaple.com' or a special wildcard '.example.com'.
func DomainMatcher(pattern string) Matcher {
	p := pattern
	if strings.HasPrefix(pattern, ".") {
		p = pattern[1:] // trim the prefix '.'
		pattern = "*" + p
	}
	return &domainMatcher{
		pattern: p,
		glob:    glob.MustCompile(pattern),
	}
}

func (m *domainMatcher) Match(domain string) bool {
	if m == nil || m.glob == nil {
		return false
	}

	if domain == m.pattern {
		return true
	}
	return m.glob.Match(domain)
}

func (m *domainMatcher) String() string {
	return "domain " + m.pattern
}

// Bypass is a filter for address (IP or domain).
// It contains a list of matchers.
type Bypass struct {
	ischain            bool
	inwall, inwall_0   bool
	chkwall, chkwall_0 bool
	white, white_0     bool
	fakeip, fakeip_0   bool

	matchers []Matcher
	period   time.Duration // the period for live reloading
	stopped  chan struct{}
	mux      sync.RWMutex
}

// NewBypass creates and initializes a new Bypass using matchers as its match rules.
// The rules will be inwall if the inwall is true.
func NewBypass(inwall, chkwall, white, fakeip, ischain bool, matchers ...Matcher) *Bypass {
	return &Bypass{
		ischain:   ischain,
		inwall:    inwall,
		inwall_0:  inwall,
		chkwall:   chkwall,
		chkwall_0: chkwall,
		white:     white,
		white_0:   white,
		fakeip:    fakeip,
		fakeip_0:  fakeip,
		matchers:  matchers,
		stopped:   make(chan struct{}),
	}
}

// NewBypassPatterns creates and initializes a new Bypass using matcher patterns as its match rules.
// The rules will be reversed if the inwall is true.
func NewBypassPatterns(inwall, chkwall, white, fakeip, ischain bool, patterns ...string) *Bypass {
	var matchers []Matcher
	for _, pattern := range patterns {
		if m := NewMatcher(pattern); m != nil {
			matchers = append(matchers, m)
		}
	}
	bp := NewBypass(inwall, chkwall, white, fakeip, ischain)
	bp.AddMatchers(matchers...)
	return bp
}
func (bp *Bypass) chkInWall(addr string) int8 {
	var l_ipchn *ipchn
	if ipchn2, found := name_ip.Load(addr); found && l_ipchn != nil {
		l_ipchn = ipchn2.(*ipchn)
	} else if ip, _ := net.ResolveIPAddr("ip4", addr); ip == nil { //无法解析的域名
		l_ipchn = &ipchn{ip: ip.String()}
		name_ip.Store(addr, l_ipchn) //加上这行可以优化无法解析的域名的响应
	} else {
		l_ipchn = &ipchn{ip: ip.String()}
		name_ip.Store(addr, l_ipchn)
	}
	log.Printf("检测[IP:%s<=DN:%s]\n", l_ipchn.ip, addr)
	if l_ipchn.inwall == 0 {
		l_ipchn.inwall = chn_wall(l_ipchn.ip)
	}
	if l_ipchn.inwall > 0 {
		log.Printf("---addr:%s位于[墙内]---\n", addr)
	} else if l_ipchn.inwall < 0 {
		log.Printf("---addr:%s位于[墙外]---\n", addr)
	} else {
		log.Printf("---addr:%s位于[未知]---\n", addr)
	}
	return l_ipchn.inwall
}
func (bp *Bypass) matchInList(addr string) bool {
	bp.mux.RLock()
	defer bp.mux.RUnlock()
	for _, matcher := range bp.matchers {
		if matcher == nil {
			continue
		}
		if matcher.Match(addr) {
			return true
		}
	}
	return false
}

// Passable reports whether the bypass includes addr.
func (bp *Bypass) Passable(addr string) bool { //Skip Pass/Bypass
	if bp == nil || len(addr) == 0 {
		loglv.Inf.Stackf("[1]Passable(%s) ret:false\n", addr)
		return false
	}

	// try to strip the port
	if host, port, _ := net.SplitHostPort(addr); host != "" && port != "" {
		if p, _ := strconv.Atoi(port); p > 0 { // port is valid
			addr = host
		}
	}
	if !bp.ischain { //前接收端
		if bp.matchInList(addr) { //在黑白名单中
			return bp.white
		} else {
			return !bp.white
		}
	} else if !bp.white && bp.matchInList(addr) { //在转发端的黑名单中，直接拒绝了
		return false
	} else if bp.fakeip && bp.white && bp.matchInList(addr) { //伪装功能打开，白名单的网站直接强制转发
		return true
	} else {
		var l_inwall int8 = -2 //不检查墙，默认为墙外<0
		if bp.chkwall {
			l_inwall = bp.chkInWall(addr)
			if l_inwall == 0 { //出错，直接跳过
				return false
			} else if !bp.inwall { //在墙外，标志反转
				l_inwall = -l_inwall
			}
		}
		if l_inwall > 0 { //墙这一边的地址不让过
			return false
		} else if !bp.white { //墙另一边的且不在黑名单，能过
			return true
		} else if bp.fakeip { //墙另一边不在白名单，不能过
			return false
		} else { //墙另一边且是白名单，在白名单能过，不在白名单不让过
			return bp.matchInList(addr)
		}
	}
}

var name_ip sync.Map

type ipchn struct {
	ip     string
	inwall int8
}

// AddMatchers appends matchers to the bypass matcher list.
func (bp *Bypass) AddMatchers(matchers ...Matcher) {
	bp.mux.Lock()
	defer bp.mux.Unlock()

	bp.matchers = append(bp.matchers, matchers...)
}

// Matchers return the bypass matcher list.
func (bp *Bypass) Matchers() []Matcher {
	bp.mux.RLock()
	defer bp.mux.RUnlock()

	return bp.matchers
}

// Reload parses config from r, then live reloads the bypass.
func (bp *Bypass) Reload(r io.Reader) error {
	var matchers []Matcher
	var period time.Duration
	inwall, chkwall, white, fakeip :=
		true, true, true, true

	if r == nil || bp.Stopped() {
		return nil
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ss := splitLine(line)
		if len(ss) == 0 {
			continue
		}
		switch ss[0] {
		case "reload": // reload option
			if len(ss) > 1 {
				period, _ = time.ParseDuration(ss[1])
			}
		case "inwall": // in_wall option
			if len(ss) > 1 {
				inwall, _ = strconv.ParseBool(ss[1])
			}
		case "chkwall": // in_wall option
			if len(ss) > 1 {
				chkwall, _ = strconv.ParseBool(ss[1])
			}
		case "white":
			if len(ss) > 1 {
				white, _ = strconv.ParseBool(ss[1])
			}
		case "fakeip":
			if len(ss) > 1 {
				fakeip, _ = strconv.ParseBool(ss[1])
			}
		default:
			matchers = append(matchers, NewMatcher(ss[0]))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	bp.mux.Lock()
	defer bp.mux.Unlock()

	bp.matchers = matchers
	bp.period = period
	bp.inwall = bp.inwall_0 && inwall
	bp.chkwall = bp.chkwall_0 && chkwall
	bp.white = bp.white_0 && white
	bp.fakeip = bp.fakeip_0 && fakeip
	return nil
}

// Period returns the reload period.
func (bp *Bypass) Period() time.Duration {
	if bp.Stopped() {
		return -1
	}

	bp.mux.RLock()
	defer bp.mux.RUnlock()

	return bp.period
}

// Stop stops reloading.
func (bp *Bypass) Stop() {
	select {
	case <-bp.stopped:
	default:
		close(bp.stopped)
	}
}

// Stopped checks whether the reloader is stopped.
func (bp *Bypass) Stopped() bool {
	select {
	case <-bp.stopped:
		return true
	default:
		return false
	}
}

func (bp *Bypass) String() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "in_wall: %v\n", bp.inwall)
	fmt.Fprintf(b, "reload: %v\n", bp.Period())
	for _, m := range bp.Matchers() {
		b.WriteString(m.String())
		b.WriteByte('\n')
	}
	return b.String()
}
