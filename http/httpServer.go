package http

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"github.com/fatezhou/tools/sys"
)

type HttpServer struct{
	FileExtFilter map[string]bool
	m sync.RWMutex
}

func (h *HttpServer)AddFilter(fileExt string){
	h.m.Lock()
	h.FileExtFilter[strings.ToLower(fileExt)] = true
	h.m.Unlock()
}

func (h *HttpServer)IsInFilter(filePath string)bool{
	file := sys.FilePath{filePath}
	ext := strings.ToLower(file.GetPathExt())
	h.m.RLock()
	defer h.m.RUnlock()
	if _, ok := h.FileExtFilter[ext]; ok{
		return true
	}
	return false
}

func (h *HttpServer)Run(ip string, port int)bool{
	strIp := fmt.Sprintf("%s:%d", ip, port)
	if nil == http.ListenAndServe(strIp, h){
		return true
	}
	return false
}

func (h *HttpServer)ServeHTTP(res http.ResponseWriter, req *http.Request){

}
