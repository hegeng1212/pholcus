package proxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"fmt"

	"github.com/hegeng1212/pholcus/app/downloader/request"
	"github.com/hegeng1212/pholcus/app/downloader/surfer"
	"github.com/hegeng1212/pholcus/common/ping"
	"github.com/hegeng1212/pholcus/config"
	"github.com/hegeng1212/pholcus/logs"
)

type Proxy struct {
	ipRegexp           *regexp.Regexp
	proxyIPTypeRegexp  *regexp.Regexp
	proxyUrlTypeRegexp *regexp.Regexp
	allIps             map[string]string
	all                map[string]bool
	online             int32
	usable             map[string]*ProxyForHost
	ticker             *time.Ticker
	tickMinute         int64
	threadPool         chan bool
	surf               surfer.Surfer
	sync.Mutex
}

const (
	CONN_TIMEOUT = 4 //4s
	DAIL_TIMEOUT = 4 //4s
	TRY_TIMES    = 3
	// IP测速的最大并发量
	MAX_THREAD_NUM = 1000
)

func New() *Proxy {
	p := &Proxy{
		ipRegexp:           regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`),
		proxyIPTypeRegexp:  regexp.MustCompile(`https?://([\w]*:[\w]*@)?[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`),
		proxyUrlTypeRegexp: regexp.MustCompile(`((https?|ftp):\/\/)?(([^:\n\r]+):([^@\n\r]+)@)?((www\.)?([^/\n\r:]+)):?([0-9]{1,5})?\/?([^?\n\r]+)?\??([^#\n\r]*)?#?([^\n\r]*)`),
		allIps:             map[string]string{},
		all:                map[string]bool{},
		usable:             make(map[string]*ProxyForHost),
		threadPool:         make(chan bool, MAX_THREAD_NUM),
		surf:               surfer.New(),
	}
	go p.Update()
	return p
}

// 代理IP数量
func (self *Proxy) Count() int32 {
	return self.online
}

// 更新代理IP列表
func (self *Proxy) Update() *Proxy {

	b, err := self.getProxyByUrl()

	if err != nil {
		f, err := os.Open(config.PROXY)
		if err != nil {
			// logs.Log.Error("Error: %v\n", err)
			return self
		}
		b, _ = ioutil.ReadAll(f)
		f.Close()
	}

	proxysIPType := self.proxyIPTypeRegexp.FindAllString(string(b), -1)
	for _, proxy := range proxysIPType {
		self.allIps[proxy] = self.ipRegexp.FindString(proxy)
		self.all[proxy] = false
	}

	proxysUrlType := self.proxyUrlTypeRegexp.FindAllString(string(b), -1)
	for _, proxy := range proxysUrlType {
		gvalue := self.proxyUrlTypeRegexp.FindStringSubmatch(proxy)
		self.allIps[proxy] = gvalue[6]
		self.all[proxy] = false
	}

	log.Printf(" *     读取代理IP: %v 条\n", len(self.all))

	self.findOnline()

	return self
}

func (self *Proxy) getProxyByUrl() (proxys []byte, err error) {

	url := "http://39.103.185.241:44567/api/v1.0/get-all-proxy-ip" // 替换成您要请求的URL

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求出错:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应出错:", err)
		return
	}

	// 使用空格作为分隔符将字符串分割成数组
	arr := strings.Split(string(body), " ")

	newArr := make([]string, 0)
	// 打印数组中的每个元素
	for _, val := range arr {
		val = fmt.Sprintf("http://%s", strings.Trim(val, " "))
		newArr = append(newArr, val)
	}

	return []byte(strings.Join(newArr, "\n")), nil
}

// 筛选在线的代理IP
func (self *Proxy) findOnline() *Proxy {
	log.Printf(" *     正在筛选在线的代理IP……")
	self.online = 0
	for proxy := range self.all {
		self.threadPool <- true
		go func(proxy string) {
			alive, _, _ := ping.Ping(self.allIps[proxy], CONN_TIMEOUT)
			self.Lock()
			self.all[proxy] = alive
			self.Unlock()
			if alive {
				atomic.AddInt32(&self.online, 1)
			}
			<-self.threadPool
		}(proxy)
	}
	for len(self.threadPool) > 0 {
		time.Sleep(0.2e9)
	}
	self.online = atomic.LoadInt32(&self.online)
	log.Printf(" *     在线代理IP筛选完成，共计：%v 个\n", self.online)

	return self
}

// 更新继时器
func (self *Proxy) UpdateTicker(tickMinute int64) {
	self.tickMinute = tickMinute
	self.ticker = time.NewTicker(time.Duration(self.tickMinute) * time.Minute)
	for _, proxyForHost := range self.usable {
		proxyForHost.curIndex++
		proxyForHost.isEcho = true
	}
}

// 获取本次循环中未使用的代理IP及其响应时长
func (self *Proxy) GetOne(u string) (curProxy string) {
	if self.online == 0 {
		return
	}
	u2, _ := url.Parse(u)
	if u2.Host == "" {
		logs.Log.Informational(" *     [%v]设置代理IP失败，目标url不正确\n", u)
		return
	}
	var key = u2.Host
	if strings.Count(key, ".") > 1 {
		key = key[strings.Index(key, ".")+1:]
	}

	self.Lock()
	defer self.Unlock()

	var ok = true
	var proxyForHost = self.usable[key]

	select {
	case <-self.ticker.C:
		proxyForHost.curIndex++
		if proxyForHost.curIndex >= proxyForHost.Len() {
			_, ok = self.testAndSort(key, u2.Scheme+"://"+u2.Host)
		}
		proxyForHost.isEcho = true

	default:
		if proxyForHost == nil {
			self.usable[key] = &ProxyForHost{
				proxys:    []string{},
				timedelay: []time.Duration{},
				isEcho:    true,
			}
			proxyForHost, ok = self.testAndSort(key, u2.Scheme+"://"+u2.Host)
		} else if l := proxyForHost.Len(); l == 0 {
			ok = false
		} else if proxyForHost.curIndex >= l {
			_, ok = self.testAndSort(key, u2.Scheme+"://"+u2.Host)
			proxyForHost.isEcho = true
		}
	}
	if !ok {
		logs.Log.Informational(" *     [%v]设置代理IP失败，没有可用的代理IP\n", key)
		return
	}
	curProxy = proxyForHost.proxys[proxyForHost.curIndex]
	if proxyForHost.isEcho {
		logs.Log.Informational(" *     设置代理IP为 [%v](%v)\n",
			curProxy,
			proxyForHost.timedelay[proxyForHost.curIndex],
		)
		proxyForHost.isEcho = false
	}
	return
}

// 测试并排序
func (self *Proxy) testAndSort(key string, testHost string) (*ProxyForHost, bool) {
	logs.Log.Informational(" *     [%v]正在测试与排序代理IP……", key)
	proxyForHost := self.usable[key]
	proxyForHost.proxys = []string{}
	proxyForHost.timedelay = []time.Duration{}
	proxyForHost.curIndex = 0
	for proxy, online := range self.all {
		if !online {
			continue
		}
		self.threadPool <- true
		go func(proxy string) {
			alive, timedelay := self.findUsable(proxy, testHost)
			if alive {
				proxyForHost.Mutex.Lock()
				proxyForHost.proxys = append(proxyForHost.proxys, proxy)
				proxyForHost.timedelay = append(proxyForHost.timedelay, timedelay)
				proxyForHost.Mutex.Unlock()
			}
			<-self.threadPool
		}(proxy)
	}
	for len(self.threadPool) > 0 {
		time.Sleep(0.2e9)
	}
	if proxyForHost.Len() > 0 {
		sort.Sort(proxyForHost)
		logs.Log.Informational(" *     [%v]测试与排序代理IP完成，可用：%v 个\n", key, proxyForHost.Len())
		return proxyForHost, true
	}
	logs.Log.Informational(" *     [%v]测试与排序代理IP完成，没有可用的代理IP\n", key)
	return proxyForHost, false
}

// 测试代理ip可用性
func (self *Proxy) findUsable(proxy string, testHost string) (alive bool, timedelay time.Duration) {
	t0 := time.Now()
	req := &request.Request{
		Url:         testHost,
		Method:      "HEAD",
		Header:      make(http.Header),
		DialTimeout: time.Second * time.Duration(DAIL_TIMEOUT),
		ConnTimeout: time.Second * time.Duration(CONN_TIMEOUT),
		TryTimes:    TRY_TIMES,
	}
	req.SetProxy(proxy)
	resp, err := self.surf.Download(req)

	if resp == nil || resp.StatusCode != http.StatusOK {
		return false, 0
	}

	return err == nil, time.Since(t0)
}
