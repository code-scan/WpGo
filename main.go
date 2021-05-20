package main

import (
	"flag"
	"github.com/code-scan/WpGo/module"
	"log"
	"os"
)

var hostFile, userFile, passFile, outFile string
var autoUser bool
var autoCount, threadCount int

func main() {

	flag.StringVar(&hostFile, "w", "", "website list filepath")
	flag.StringVar(&passFile, "p", "", "password list filepath")
	flag.IntVar(&autoCount, "c", 5, "max auto get user count")
	flag.IntVar(&threadCount, "t", 20, "max thread")
	flag.StringVar(&outFile, "o", "result.txt", "out filepath")
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if hostFile == "" || passFile == "" {
		log.Println("website/password file is null")
		log.Println("usage:")
		log.Println(os.Args[0], " -h")

		return
	}
	var passlist []string
	var hostlist []string

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
	for _, site := range hostlist {
		module.SiteQueue <- site
	}
	for i := 0; i < threadCount; i++ {
		module.Wg.Add(2)
		go module.NewSend(passlist, autoCount)
		go module.NewWork()
	}
	module.Wg.Wait()
}
