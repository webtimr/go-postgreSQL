package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/webtimr/go-postgreSQL.git/router"
)

func main() {
	r := router.Router()
	fmt.Println("start server on the port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
 