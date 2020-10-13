// main.go
package main

import (
	"net/http"

	greeting "github.com/aydenhex/US-Greetings"
)

func main() {
	service := greeting.NewInmemGreetingService()
	endpoints := greeting.MakeGreetingEndpoints(service)

	err := http.ListenAndServe(":8000", greeting.MakeTTPHandler(endpoints))
	if err != nil {
		panic(err)
	}
}
