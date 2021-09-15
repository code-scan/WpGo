package module

type Service interface {
	Login(SiteTask) bool    // 登录
	Check(SiteTask) bool    // 判断是否是符合条件
	CheckCache(string) bool //获取缓存
}
