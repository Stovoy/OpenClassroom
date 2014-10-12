package main

import (
	"oc/db"
	"oc/server"
)

func main() {
	db.Start()
	server.Start()
}
