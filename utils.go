package main

import (
	"bytes"	
	"fmt"
	"github.com/gorilla/mux"	
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"runtime"	
)

/**************************************************************
	Utility to print debug information for an http.Request
 **************************************************************/
func DebugInfo(r *http.Request) {
	if !DEBUG {
		return
	}
	Debug("URL:%v ", r.URL.String())
	Debug("Scheme:%v ", r.URL.Scheme)
	Debug("Opaque:%v ", r.URL.Opaque)
	Debug("Host:%v ", r.URL.Host)
	Debug("Path:%v ", r.URL.Path)
	Debug("RawQuery:%v ", r.URL.RawQuery)
	Debug("Fragment:%v ", r.URL.Fragment)
	Debug("Method:%v ", r.Method)
	Debug("Proto:%v ", r.Proto)
	Debug("Header:%v ", r.Header)
	Debug("ContentLength:%v ", r.ContentLength)
	Debug("Host:%v ", r.Host)
	Debug("Referer:%v ", r.Referer())
	Debug("Form:%v ", r.Form)
	Debug("PostForm:%v ", r.PostForm)
	Debug("MultipartForm:%v ", r.MultipartForm)
	Debug("RemoteAddr:%v ", r.RemoteAddr)
	Debug("RequestURI:%v ", r.RequestURI)
	for k, v := range r.Header {
		Debug("Header %v = %v ", k, v)
	}
	for _, v := range r.Cookies() {
		Debug("Cookie %v = %v", v.Name, v.Value)
	}
	for k, v := range mux.Vars(r) {
		Debug("mux Vars %v = %v ", k, v)
	}
	for _, v := range r.Cookies() {
		Debug("Cookie %v = %v", v.Name, v.Value)
	}
	request, err := httputil.DumpRequest(r, true)
	if err != nil {
		Debug("Error while dumping request: %v", err)
		return
	}
	Debug("Request: %v", string(request))
}

/**************************************************************
	Debug utility
	File /var/log/web-1.log
 **************************************************************/
func Debug(format string, a ...interface{}) {
	if !DEBUG {
		return
	}
	fmt.Printf(format+"\n", a...)
}

/**************************************************************
	Info utility
	File /var/log/web-1.log
 **************************************************************/
func Info(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

/**************************************************************
	Error utility
	Error in file /var/log/web-1.error.log
	Info in file /var/log/web-1.log
 **************************************************************/
func Error(format string, a ...interface{}) {
	Info("ERROR: "+format, a)
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a...)
}

/**************************************************************
	Utility to convert a []byte to a string
 **************************************************************/
func B2S(b []byte) (s string) {
	n := bytes.Index(b, []byte{0})
	if n > 0 {
		s = string(b[:n])
	} else {
		s = string(b)
	}
	return
}

/**************************************************************
	Utility to convert a Bytes to a MegaBytes
 **************************************************************/
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/**************************************************************
	Utility to pretty print memory usage
 **************************************************************/
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	Debug("Alloc = %v MiB \t TotalAlloc = %v MiB \t Sys = %v MiB", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys))
}

/**************************************************************
	Utility to pretty print OS information
 **************************************************************/
func DebugOS() {
	if !DEBUG {
		return
	}
	Debug("Environment variables:")
	for _, e := range os.Environ() {
		Debug("%v", e)
	}
	Debug("Process id: %v", os.Getpid())
	Debug("Parent Process id: %v", os.Getppid())
	if host, err := os.Hostname(); err == nil {
		Debug("Hostname: %v", host)
	}
}

/**************************************************************
	Utility to remove a directory and all its contents
 **************************************************************/
func RemoveDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		Error("Error removing temp folder: %v", err)
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		Error("Error removing temp folder: %v", err)
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			Error("Error removing temp folder: %v", err)
			return err
		}
	}
	return os.Remove(dir)
}
