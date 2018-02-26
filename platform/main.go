package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"time"

	"github.com/o3labs/openpoint/platform/config"
	"github.com/o3labs/openpoint/platform/middleware"
	"github.com/o3labs/openpoint/platform/router"

	"github.com/gorilla/mux"
)

//Server mux
type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, DEL")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Channel-Application-Key, X-Channel-Client-ID, X-Channel-Bot-Token, X-Channel-Story-ID")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	record := &middleware.LogRecord{
		ResponseWriter: w,
	}

	t1 := time.Now()
	s.r.ServeHTTP(record, r)
	t2 := time.Now()
	if record.Status != 0 && record.Status != 301 {
		log.SetPrefix("\x1b[31;1m[Error] ")
		log.Printf("[%s] %q %v %v\n", r.Method, r.URL.String(), record.Status, t2.Sub(t1))
		log.Printf("%v\n", r.Header)
	} else {
		log.SetPrefix("\x1b[0m[Info]  ")
		log.Printf("[%s] %q %v %v\n", r.Method, r.URL.String(), http.StatusOK, t2.Sub(t1))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	mode := flag.String("mode", "", "Run in which mode. local | staging | production")
	flag.Parse()

	if *mode == "" {
		//default mode is local
		defaultEnv := "local"
		mode = &defaultEnv
	}
	pathPrefix := "/"
	r := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()

	http.Handle("/", &Server{r})

	config.Load(*mode)

	//config routing here
	router.RegisterRouting(r)

	log.Printf("Running %+v Environment Listening... %+v", *mode, config.Env.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", config.Env.Port), nil)

}
