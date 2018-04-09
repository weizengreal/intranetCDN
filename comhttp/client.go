package comhttp

// httputil 客户端请求方法

import (
	"../base"
	"net/http"
	"crypto/tls"
	"log"
	"io"
	"strconv"
)

/**
 发送 http GET 请求，分段下载
 */
func SendGet(url string, path string, start int64, end int64) (length int64,err error) {
	var req *http.Request
	req ,err = http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("NewRequese faild!",err)
		return 0,err
	}
	req.Header.Set("Range","bytes=" + strconv.FormatInt(start,10) + "-" + strconv.FormatInt(end,10))
	req.Header.Set("Connection","close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:true,
		},
	}
	client := http.Client{
		Transport : tr,
	}
	var resp *http.Response
	resp ,err = client.Do(req)
	defer resp.Body.Close()
	file ,err := base.CreateFile(path)
	defer file.Close()
	length, err = io.Copy(file,resp.Body)
	return length,err
}

/**
 发送 HEAD 请求，获取资源基本信息
 */
func SendHead(url string) (length int64,support bool,md5 string, err error) {
	var req *http.Request
	req ,err = http.NewRequest("HEAD",url,nil)
	if err != nil {
		return -1,false,"",err
	}
	// 要求服务器返回最新的数据而不是缓存
	req.Header.Set("Cache-Control","no-cache")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:true,
		},
	}
	client := http.Client{
		Transport : tr,
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return -1,false,"",err
	}
	defer resp.Body.Close()
	length,err = strconv.ParseInt(resp.Header.Get("Content-Length"),10,64)
	if err != nil {
		return -1,false,"",err
	}
	// Accept-Ranges 不存在，不支持断点续传
	if resp.Header.Get("Accept-Ranges") != "" {
		support = true
	}
	md5 = resp.Header.Get("Content-Md5")
	return length,support,md5,err
}