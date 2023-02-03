package utils

import "net/http"

type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total int
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {

}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, -1, data, msg)
}
