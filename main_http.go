package main

import (
	"bytes"
	"log"
	"net/http"
)

func main() {
	router.Start()
	http.HandleFunc("/", r)

	var err error
	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func r(rw http.ResponseWriter, req *http.Request) {
	route, _, err := router.FindRoute(req.URL.Path)
	if route == nil || err != nil {
		return
	}
	routeMap := route.Dest.(map[string]func(map[string]interface{}, []byte) ([]byte, error, int))
	handler := routeMap[req.Method]
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	params := map[string]interface{}{}
	jsonData, _, _ := handler(params, buf.Bytes())
	rw.Header().Set("Content-Type", "text/html")
	rw.Write(jsonData)
}
