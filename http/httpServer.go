package http

import (
	"strings"
	"sync"
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
	


	npos := strings.LastIndex(filePath, ".")
	if npos != -1{
		ext := filePath[npos + 1:]
		h.m.RLock()
		defer h.m.RUnlock()
		if _, ok := h.FileExtFilter[strings.ToLower(ext)]; ok{
			return true
		}
	}
	return false
}


