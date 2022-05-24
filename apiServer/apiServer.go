package main

import (
	"dsxy/apiServer/heartbeat"
	"dsxy/apiServer/locate"
	"dsxy/apiServer/objects"
	"dsxy/apiServer/versions"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	fmt.Println("Loading...")
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}
