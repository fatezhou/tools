package http

import (
"fmt"
"net/http"
"strings"
)

type CmdHandler func(data map[string]string)(strResponseText string)

type ApiHttpServer struct{
	router map[string]CmdHandler
}

func (ApiHttpServer *ApiHttpServer)ServeHTTP(writer http.ResponseWriter, request *http.Request){
	strCmd := strings.ToLower(request.URL.Path)
	if handle, ok := ApiHttpServer.router[strCmd]; ok{
		mapData := make(map[string]string)
		value := request.URL.Query()
		for k, v := range value{
			if len(v) > 0{
				mapData[k] = v[0]
			}
		}
		response := handle(mapData)
		strJson := fmt.Sprintf(`
{
	"code":0, 
	"data":%s
}`		, response)
		writer.Write([]byte(strJson))
	}else{
		writer.Write([]byte(`{"code":-1, "data":{"text":"type http://ip:port/help to get cmds"}}`))
	}
}

func (ApiHttpServer *ApiHttpServer)AddHandle(cmd string, handle CmdHandler){
	//     cmd:  /cmd/xxxxx
	cmd = strings.ToLower(cmd)
	if ApiHttpServer.router == nil{
		ApiHttpServer.router = make(map[string]CmdHandler)
	}
	ApiHttpServer.router[cmd] = handle
}

func (ApiHttpServer *ApiHttpServer)Run(port int){
	strHttpIp := fmt.Sprintf("0.0.0.0:%d", port)
	http.ListenAndServe(strHttpIp, ApiHttpServer)
}

func (ApiHttpServer *ApiHttpServer)RunAync(port int){
	go func(){
		ApiHttpServer.Run(port)
	}()
}