package comhttp

import "../base"

var contextMap map[string] *base.Context = make(map[string] *base.Context)

type ApiResult struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

// 判断某个 key-value 是否存在
func IsExist(url string) bool {
	if _,ok := contextMap[url];!ok {
		return false
	}
	return true
}

// 获得某个 url 的上下文信息
func GetUriContext(url string) *base.Context {
	if IsExist(url) {
		return contextMap[url]
	} else {
		return nil
	}
}

// 向 map 在中初始化一个上下文信息
func AddContext(url string) *base.Context {
	context := new(base.Context)
	contextMap[url] = context
	return context
}

// 删除 map 中的一个上下文
func DelContext(url string)  {
	delete(contextMap,url)
}