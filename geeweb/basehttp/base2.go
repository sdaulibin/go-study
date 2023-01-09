package basehttp

import (
	"fmt"
	"net/http"
)

type Engine struct{}

func (engine Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "======%q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "======%q>>>>>>%q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 not found url %q\n", r.URL)
	}
}

// func main() {
// 	engine := new(Engine)
// 	log.Fatal(http.ListenAndServe("127.0.0.1:9999", engine))
// }
