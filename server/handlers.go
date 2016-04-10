package server

import (
	"fmt"
	"net/http"
)

func panicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, r *http.Request, i interface{}) {
		fmt.Println("wtf panic handler")
	}
}
