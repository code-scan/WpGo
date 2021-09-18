package module

import (
	"log"
	"os"
)

type Service interface {
	Login(SiteTask) bool    // 登录
	Check(SiteTask) bool    // 判断是否是符合条件
	CheckCache(string) bool //获取缓存
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
