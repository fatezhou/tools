package sys

import "strings"

type FilePath struct{
	Path string
}

func (f *FilePath)GetPathExt()string{
	pos := strings.LastIndex(f.Path, ".")
	if pos == -1{
		return ""
	}else{
		return f.Path[1 + pos:]
	}
}

func (f *FilePath)GetFileName()string{
	pos := strings.LastIndex(f.Path, "\\")
	pos2 := strings.LastIndex(f.Path, "/")
	if pos > pos2{
		return f.Path[pos + 1:]
	}else{
		return f.Path[pos2 + 1:]
	}
}

func (f *FilePath)GetDir()string{
	pos := strings.LastIndex(f.Path, "\\")
	pos2 := strings.LastIndex(f.Path, "/")
	if pos > pos2{
		return f.Path[:pos]
	}else{
		return f.Path[:pos2]
	}
}