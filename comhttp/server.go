package comhttp

import (
	"../download"
	"../base"
	"net/http"
	"fmt"
	"encoding/json"
	"time"
)

// 作为服务端搭建 http 服务器
func Server() {
	http.HandleFunc("/",handler)
	http.HandleFunc("/downloadStable",downloadStable)
	http.HandleFunc("/downloadChunk",downloadChunk)
	http.ListenAndServe("localhost:8000",nil)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"%s","hello world,I am weizeng!")
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
	//jsonbytes,_ := json.Marshal(context)
	//fmt.Println(string(jsonbytes))

	w.Header().Set("Content-Type",context.Res.Header["Content-Type"])
	w.Header().Set("Content-Length",context.Res.Header["Content-Length"])

	bytes := base.ReadFile(context.Res.Path)
	w.Write(bytes)
}

func downloadChunk(w http.ResponseWriter, r *http.Request) {
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
	// 检测当前的资源是否需要重新下载
	if context.Res == nil || context.Res.FileMd5 == "" || download.IsChange(context){
		download.AssignInit(uri,context)
		go func(context *base.Context) {
			download.DownloadAllBlock(context)
			download.MergeBlock(context)
		}(context)
	}

	w.Header().Set("Content-Type",context.Res.Header["Content-Type"])
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", context.Res.Name))

	if context.Res.IsComplete {
		w.Header().Set("Content-Length",context.Res.Header["Content-Length"])
		w.Write(base.ReadFile(context.Res.Path))
	} else {
		for i := 0; i < len(context.TmpPath); i++ {
			block := context.FileMap[context.TmpPath[i]]
			for {
				if block.Status {
					w.Write(base.ReadFile(context.TmpPath[i]))
					break
				}
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}
	}
}

