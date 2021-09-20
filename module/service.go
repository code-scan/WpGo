package module

import (
	"encoding/xml"
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

type MethodResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text  string `xml:",chardata"`
				Array struct {
					Text string `xml:",chardata"`
					Data struct {
						Text  string `xml:",chardata"`
						Value []struct {
							Text   string `xml:",chardata"`
							Struct struct {
								Text   string `xml:",chardata"`
								Member []struct {
									Text  string `xml:",chardata"`
									Name  string `xml:"name"`
									Value struct {
										Text   string `xml:",chardata"`
										Int    string `xml:"int"`
										String string `xml:"string"`
									} `xml:"value"`
								} `xml:"member"`
							} `xml:"struct"`
						} `xml:"value"`
					} `xml:"data"`
				} `xml:"array"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}
