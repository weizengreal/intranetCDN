package comhttp

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

type ApiResult struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

// 作为服务端搭建 http 服务器

func Server() {
	http.HandleFunc("/",handler)
	http.HandleFunc("/download",download)
	http.ListenAndServe("localhost:8000",nil)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println(*r)
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	res := &ApiResult{
		Status:1,
		Message:"ok",
	}
	if err := r.ParseForm(); err != nil {
		res.Status = -1
		res.Message = err.Error()
		fmt.Println(res)
		bytes,err := json.Marshal(res)
		if err != nil {
			fmt.Fprintf(w,"%s","json marshal error!")
		}
		fmt.Fprintf(w,"%s",string(bytes))
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
		fmt.Println(string(bytes))
		fmt.Fprintf(w,"%s",string(bytes))
		return
	}

	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	fmt.Println(*r)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.Form.Get("uri"))

}
