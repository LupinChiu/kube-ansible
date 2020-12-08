package content

import (
	"fmt"
	"net/http"
	"server/conf"
	database "server/db"
	"time"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

type Http struct {
	config conf.Conf
}

type echoHandler int

func (e echoHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	fmt.Println("receive from:", request.RemoteAddr)
	name := database.GetTestingSQLService()
	response.Write([]byte(name + ":" + time.Now().Format("2006-01-02 15:04:05")))
}

func (h *Http) Start() {
	var handler echoHandler
	//router.HandleFunc("/echo", echoHandler)
	http.Handle("/echo", handler)
	fmt.Println("http server start.")
	fmt.Println(h.config.Proto.Port)
	srv := &http.Server{
		Addr:         h.config.Proto.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
}
