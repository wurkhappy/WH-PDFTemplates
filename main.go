// +build !DEBUG

package main

import (
	"encoding/json"
	"flag"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/mdp"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	// "time"
)

var production = flag.Bool("production", false, "Production settings")

func main() {
	flag.Parse()
	router.Start()
	if *production {
		config.Prod()
		log.SetOutput(os.Stdout)
	} else {
		config.Test()
	}

	gophers := 1

	// Create a channel to talk with the OS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, syscall.SIGTERM)
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	sigChan <- true
	// }()

	// Create a channel to shut down the program early
	shutChan := make(chan bool)
	var wg sync.WaitGroup

	for i := 0; i < gophers; i++ {
		log.Println(config.PDFTemplatesService)
		worker := mdp.NewWorker(config.MDPBroker, config.PDFTemplatesService, false)
		defer worker.Close()
		go route(worker, shutChan, wg)
	}

	select {
	case <-sigChan:
		log.Println("Main", "controller.Run", "******> Program Being Killed")

		// Signal the program to shutdown and wait for confirmation
		for i := 0; i < gophers; i++ {
			log.Println("shut chan")
			shutChan <- true
		}
	}
	wg.Wait()

	log.Println("Main", "controller.Run", "******> Shutting Down")
	return
}

type Resp struct {
	Body       []byte `json:"body"`
	StatusCode int    `json:"status_code"`
}

type ServiceReq struct {
	Method string
	Path   string
	Body   []byte
	UserID string
}

func route(worker mdp.Worker, shutChan chan bool, wg sync.WaitGroup) {
	wg.Add(1)
	for reply := [][]byte{}; ; {
		request := worker.Recv(reply, shutChan)
		if len(request) == 0 {
			log.Print("wg done")
			wg.Done()
			break
		}
		var req *ServiceReq
		json.Unmarshal(request[0], &req)

		log.Println(req.UserID, req.Path, req.Method, string(req.Body))

		//route to function based on the path and method
		route, pathParams, err := router.FindRoute(req.Path)
		if route == nil || err != nil {
			return
		}
		routeMap := route.Dest.(map[string]func(map[string]interface{}, []byte) ([]byte, error, int))
		handler := routeMap[req.Method]

		//add url params to params var
		params := make(map[string]interface{})
		for key, value := range pathParams {
			params[key] = value
		}
		//add url query params
		uri, _ := url.Parse(req.Path)
		values := uri.Query()
		for key, value := range values {
			params[key] = value
		}
		params["userID"] = req.UserID

		//run handler and do standard http stuff(write JSON, return err, set status code)
		jsonData, err, statusCode := handler(params, req.Body)
		if err != nil {
			log.Println(req.UserID, req.Path, req.Method, "ERROR", err.Error())
			resp := &Resp{[]byte(`{"description":"` + err.Error() + `"}`), statusCode}
			d, _ := json.Marshal(resp)
			reply = [][]byte{d}
			continue
		}
		resp := &Resp{jsonData, statusCode}
		d, _ := json.Marshal(resp)
		reply = [][]byte{d}

	}
}
