package main

import (
	"github.com/ant0ine/go-urlrouter"
)

var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/template/invoice",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": invoice,
			},
		},
		urlrouter.Route{
			PathExp: "/template/agreement",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": agreement,
			},
		},
	},
}
