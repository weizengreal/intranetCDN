package comhttp

import (
	"../download"
	"../base"
	"net/http"
	"fmt"
	"encoding/json"
	//"time"
	"net/textproto"
)

// 作为服务端搭建 http 服务器
func Server() {
	http.HandleFunc("/",handler)
	http.HandleFunc("/downloadStable",downloadStable)
	http.HandleFunc("/downloadChunk",downloadChunk)
	http.ListenAndServe("localhost:8000",nil)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	//w.Header().Set("Content-Type","application/zip")
	//bytes := base.ReadFile("/Users/weizeng/vagrant/go/intranetCDN/cache/3a6c910b45e644739a80522bfd92d4ea.zip")
	//w.Write(bytes)
	fmt.Fprintf(w,"%s","request will sleep 3 second! \r\n")
	//time.Sleep(time.Duration(3) * time.Second)
	fmt.Fprintf(w,"%s","i am weak!")
}

func downloadStable(w http.ResponseWriter, r *http.Request) {
	res := &ApiResult{
		Status:1,
		Message:"ok",
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w,err.Error(),500)
		return
	}
	uri := r.Form.Get("uri")
	if uri =="" {
		res.Status = -2
		res.Message = "lose param uri!"
		bytes,err := json.Marshal(res)
		if err != nil {
			fmt.Fprintf(w,"%s","json marshal error!")
		}
		fmt.Fprintf(w,"%s",string(bytes))
		return
	}
	// 获取该资源的上线文信息
	context := GetUriContext(uri)
	if context == nil {
		context = AddContext(uri)
	}

	// 检测当前的资源是否需要重复下载
	if context.Res == nil || context.Res.FileMd5 == ""{
		download.Download(uri,context)
	}
	fmt.Println(*context)
	w.Header().Set("Content-Type",textproto.MIMEHeader(context.Res.Header).Get("Content-Type"))
	w.Header().Set("Content-Length",textproto.MIMEHeader(context.Res.Header).Get("Content-Length"))
	bytes := base.ReadFile(context.Res.Path)
	w.Write(bytes)
}

func downloadChunk(w http.ResponseWriter, r *http.Request) {

}
