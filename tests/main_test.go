package tests

import (
	"github.com/code-scan/WpGo/module"
	"log"
	"testing"
)

func TestLogin(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	w := module.NewWpGo()
	//task := module.SiteTask{
	//	Host: "https://50.116.84.244",
	//	User: "admin",
	//	Pass: "password1",
	//}

	//task := module.SiteTask{
	//	Host: "https://saiyo-ac.jp",
	//	User: "master",
	//	Pass: "homelesspa",
	//}

	task := module.SiteTask{
		Host: "http://1sad12392.168.81.130/wordpress",
		User: "admin",
		Pass: "caonima123",
	}
	w.Login(task)

}
