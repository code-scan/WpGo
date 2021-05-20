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

	task := module.SiteTask{
		Host: "https://167.172.19.58/wp-login.php",
		User: "admin",
		Pass: "caonima123",
	}
	w.Login(task)

}
