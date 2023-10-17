package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Pallinder/sillyname-go"
	"github.com/gorilla/mux"
	// "path/filepath"
	// "sync"
	// "text/template"
)

// templ represents a single template
// type templateHandler struct {
// 	once     sync.Once
// 	filename string
// 	templ    *template.Template
// }

// ServeHTTP handles the HTTP request.
// func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	t.once.Do(func() {
// 		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
// 	})
// 	t.templ.Execute(w, r)
// }

type name struct{}

func (s *name) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "hx-request, hx-current-url, hx-target")

	name := sillyname.GenerateStupidName()
	htmlString := fmt.Sprintf(`<input id="nameInput" type="text" name="name" value="%s" class="focus:outline-none px-4 placeholder-gray-600 bg-[#f0f0f1] rounded-md py-2" required />`, name)
	htmlBytes := []byte(htmlString)

	w.Write(htmlBytes)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	router := mux.NewRouter()
	r := newRoom()

	// http.Handle("/", &templateHandler{filename: "chat.html"})
	router.Handle("/name", &name{})
	router.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
