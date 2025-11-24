package main

import (
	"commenttree/api"
	"commenttree/front"
	"commenttree/service"
	"fmt"
	"net/http"
)

func main() {
	serviceImpl := service.NewCommentService()
	server := api.NewServer(serviceImpl)

	server.Router.Get("/", front.WebInterfaceHandler)

	fmt.Println("Server is running on http://localhost:8082")
	err := http.ListenAndServe(":8082", server.Router)
	if err != nil {
		panic(err)
	}
}
