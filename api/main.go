package main

import "API/service"

func main() {
	s := service.NewService()
	s.StartWebService()
}
