package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

/**************************************************************
	Define Global Variables
 **************************************************************/

// If DEBUG is set to true, Debug information are printed out in /var/log/web-1.log
var DEBUG = true

// HTML Template for the index page
var indexTemplate = template.Must(template.New("index-template.html").
	Delims("[[", "]]").ParseFiles("templates/index-template.html"))

/**************************************************************
	Define Index (/) Page Handler
 **************************************************************/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	Info(">>>>> indexHandler")
	DebugInfo(r)
	PrintMemUsage()

	// Render index page
	if err := indexTemplate.Execute(w, template.FuncMap{
		"Version": 123,
	}); err != nil {
		Error("Error with indexTemplate: %v", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

/**************************************************************
	Handler to debug an http.Request to the web user
 **************************************************************/
func dumpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL:%v \n", r.URL)
	fmt.Fprintf(w, "Method:%v \n", r.Method)
	fmt.Fprintf(w, "Proto:%v \n", r.Proto)
	fmt.Fprintf(w, "Header:%v \n", r.Header)
	fmt.Fprintf(w, "ContentLength:%v \n", r.ContentLength)
	fmt.Fprintf(w, "Host:%v \n", r.Host)
	fmt.Fprintf(w, "Referer:%v \n", r.Referer())
	fmt.Fprintf(w, "Form:%v \n", r.Form)
	fmt.Fprintf(w, "PostForm:%v \n", r.PostForm)
	fmt.Fprintf(w, "MultipartForm:%v \n", r.MultipartForm)
	fmt.Fprintf(w, "RemoteAddr:%v \n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI:%v \n", r.RequestURI)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header %v = %v \n", k, v)
	}

	for _, v := range r.Cookies() {
		fmt.Fprintf(w, "Cookie %v = %v \n", v.Name, v.Value)
	}

	for _, v := range os.Environ() {
		fmt.Fprintf(w, "Env  %v \n", v)
	}
	request, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintf(w, "Error while dumping request: %v\n", err)
		return
	}
	fmt.Fprintf(w, "Request: %v\n", string(request))
}

/**************************************************************
	Main program
	Logs in /var/log/web-1.log and /var/log/web-1.error.log
 **************************************************************/
func main() {
	Info(">>>>> main")
	DebugOS()
	PrintMemUsage()

	// Define HTTP Router
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/dump", dumpHandler)
	r.HandleFunc("/event", eventHandler)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	r.PathPrefix("/").HandlerFunc(indexHandler) // Catch-all
	http.Handle("/", r)

	// Set port (default to 5000)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		Info("Defaulting to port %s\n", port)
	}

	// Serve application (plain HTTP protocol within Elastic Beanstalk network)
	Info("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		Error("Error with ListenAndServe: %v", err)
		log.Fatal(err)
	}

}
