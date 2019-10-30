package main
import(
	"github.com/fatezhou/tools/http"
)

func main(){
	s := http.HttpServer{}
	s.DefaultInit()
	s.AddFilter("js")
	s.AddFilter("html")
	s.Run("0.0.0.0", 33445)
}
