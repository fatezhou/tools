package http

import (
	"fmt"
	"github.com/fatezhou/tools/cache"
	"github.com/fatezhou/tools/sys"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ServerHandler func(data map[string]string, writer http.ResponseWriter)

type HttpServerConfig struct{
	Ip string 	`json:"ip"`
	Port int	`json:"port"`
	CacheExpire int64	`json:"cache_expire"`
	CacheExpireAdd int64	`json:"cache_expire_add"`
	StaticDir string	`json:"static_dir"`
}

type HttpServer struct{
	FileExtFilter map[string]bool
	handler map[string]ServerHandler
	m sync.RWMutex
	cache cache.HashCache
	staticDir string
}

const (
	static_file_ok int = 0
	static_file_is_not int = 1
	static_file_404 int = 2
)


func (h *HttpServer)AddFilter(fileExt string){
	h.m.Lock()
	h.FileExtFilter[strings.ToLower(fileExt)] = true
	h.m.Unlock()
}

func (h *HttpServer)AddHandler(strPath string, handler ServerHandler){
	h.m.Lock()
	h.handler[strings.ToLower(strPath)] = handler
	h.m.Unlock()
}

func (h *HttpServer)GetHandler(strPath string)ServerHandler{
	h.m.Lock()
	defer h.m.Unlock()
	if handler, ok := h.handler[strPath]; ok{
		return handler
	}
	return nil
}

func (h *HttpServer)IsInFilter(filePath string)bool{
	file := sys.FilePath{Path:filePath}
	ext := strings.ToLower(file.GetPathExt())
	h.m.RLock()
	defer h.m.RUnlock()
	if _, ok := h.FileExtFilter[ext]; ok{
		return true
	}
	return false
}

func (h *HttpServer)DefaultInit(){
	h.cache.SetConfig(3600, 600)
	h.staticDir = "d:/static/"
	h.FileExtFilter = make(map[string]bool)
	h.handler = make(map[string]ServerHandler)
}

func (h *HttpServer)Init(conf HttpServerConfig){
	h.cache.SetConfig(conf.CacheExpire, conf.CacheExpireAdd)
	h.staticDir = conf.StaticDir
	h.FileExtFilter = make(map[string]bool)
	h.handler = make(map[string]ServerHandler)
}

func (h *HttpServer)Run(ip string, port int)bool{
	strIp := fmt.Sprintf("%s:%d", ip, port)
	if nil == http.ListenAndServe(strIp, h){
		return true
	}
	return false
}

func (h *HttpServer)ServeHTTP(res http.ResponseWriter, req *http.Request){
	staticState := h.DoStaticHttpRequest(res, req)
	if static_file_404 == staticState{
		h.Redirect404(res, req)
	}else if static_file_is_not == staticState{
		handler := h.GetHandler(req.URL.Path)
		if handler != nil{
			mapData := make(map[string]string)
			value := req.URL.Query()
			for k, v := range value{
				if len(v) > 0{
					mapData[k] = v[0]
				}
			}
			handler(mapData, res)
		}else{
			h.Redirect404(res, req)
		}
	}else{
		// finish do static request
	}
}

func (h *HttpServer)Redirect404(res http.ResponseWriter, req *http.Request){
	res.Header().Add("Location", "http://" + req.Host + "/404.html")
	res.WriteHeader(302)
	res.Write([]byte(""))
}

func (h *HttpServer)MakeDefault404(res http.ResponseWriter, req *http.Request){
	res.WriteHeader(404)
	res.Write([]byte("<h1>404</h1>"))
}

func (h *HttpServer)DoStaticHttpRequest(res http.ResponseWriter, req *http.Request)int{
	urlPath := strings.ToLower(req.URL.Path)
	filePath := sys.FilePath{Path:urlPath}
	if false == h.IsInFilter(urlPath){
		return static_file_is_not
	}

	fileData := h.cache.SaveGet(urlPath)
	if fileData != nil{
		res.Write(fileData.([]byte))
		return static_file_ok
	}else{
		file, err := os.Open(h.staticDir + filePath.GetFileName())
		defer file.Close()
		if err != nil{
			if strings.Contains(urlPath, "404.html"){
				h.MakeDefault404(res, req)
				return static_file_ok
			}
			return static_file_404
		}else{
			b, _ := ioutil.ReadAll(file)
			h.cache.SaveSet(urlPath, b, 3600)
			res.Write(b)
			return static_file_ok
		}
	}
}
