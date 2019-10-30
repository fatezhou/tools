package http

import (
	"fmt"
	"net/http"
	"strings"
)

type ServerHandler func(data map[string]string, writer http.ResponseWriter)

type HttpServer struct{
	router map[string]ServerHandler
}


func (httpServer *HttpServer)ServeHTTP(writer http.ResponseWriter, request *http.Request){
	strCmd := strings.ToLower(request.URL.Path)
	if handle, ok := httpServer.router[strCmd]; ok {
		mapData := make(map[string]string)
		value := request.URL.Query()
		for k, v := range value {
			if len(v) > 0 {
				mapData[k] = v[0]
			}
		}
		handle(mapData, writer)
		return
	}else{
		//404
		writer.Write([]byte("<h1>404</h1>"))
	}
}

func (httpServer *HttpServer)AddHandle(cmd string, handle ServerHandler){
	//     cmd:  /cmd/xxxxx
	cmd = strings.ToLower(cmd)
	if httpServer.router == nil{
		httpServer.router = make(map[string]ServerHandler)
	}
	httpServer.router[cmd] = handle
}

func (httpServer *HttpServer)Run(ip string, port int){
	strHttpIp := fmt.Sprintf("%s:%d", ip, port)
	http.ListenAndServe(strHttpIp, httpServer)
}

func (httpServer *HttpServer)RunAync(ip string, port int){
	go func(){
		httpServer.Run(ip, port)
	}()
}
